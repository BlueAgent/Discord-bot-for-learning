package main
//Made with heavy reference to discordgo examples

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"sync"
	"time"
)

type LastData struct {
	Sync *sync.Mutex
	Message string
	Counter int
	Reply chan int
}

type Bot struct {
	LastSync sync.Mutex
	Last map[string]LastData
	BotID string
}

func NewBot(botId string) *Bot {
	b := Bot{
		Last: make(map[string]LastData),
		BotID: botId,
	}
	return &b
}

func main() {
	var (
		botToken = flag.String("token", "", "Bot Token")
	)
	flag.Parse()
	if botToken == nil {
		fmt.Fprintf(os.Stderr, "Missing bot token.\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "BotToken: %s\n", *botToken)
	
	//Connect to Discord using bot token
	dg, err := discordgo.New("Bot " + *botToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Discord session: %s\n", err)
		os.Exit(1)
	}

	//Get account information
	u, err := dg.User("@me")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to retrive account information: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "BotID: %s\n", u.ID)

	b := NewBot(u.ID)

	dg.AddHandler(b.MessageCreate)

	//Open websocket and start listening
	err = dg.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open websocket: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Bot has started. Exit with CTRL-C.")

	//This works by reading from a channel that never gets anything written to it
	<-make(chan struct{})
	return
}

func MessageIdentifier(m *discordgo.Message) string {
	var id = m.Content	
	if len(m.Attachments) > 0 {
		for _,a := range m.Attachments {
			id += "|" + a.URL
		}
	}
	if len(m.Embeds) > 0 {
		for _,e := range m.Embeds {
			id += "|" + e.URL
		}
	}
	return id
}

func (b *Bot) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == b.BotID {
		return
	}

	var aID = m.Author.ID

	//Retrive Last
	data, ok := b.Last[aID]
	
	//If not found, get a lock to create one
	if !ok {
		b.LastSync.Lock()
		//Check Again
		if data, ok = b.Last[aID]; !ok {
			//Safe to create now
			b.Last[aID] = LastData{
				Sync: new(sync.Mutex),
				Message: "",
				Counter: 1,
				Reply: nil,
			}
		}
		b.LastSync.Unlock()
		data, ok = b.Last[aID]
	}

	data.Sync.Lock()
	var msgID = MessageIdentifier(m.Message)

	if data.Message == msgID {
		data.Counter++
		go s.ChannelMessageDelete(m.ChannelID, m.ID)
		if data.Reply == nil {
			data.Reply = make(chan int)
			go ReplyCreate(s, m.ChannelID, m.Author.Username, m.Content, data.Reply)
			fmt.Fprintf(os.Stdout, "Limiting [%s]: %s\n", m.Author.Username, m.Content)
		}
		go func(count int, cc chan int) {
			cc <- count
		}(data.Counter, data.Reply)
	} else {
		data.Counter = 1
		if data.Reply != nil {
			go func(cc chan int) {
				cc <- -1
			}(data.Reply)
			data.Reply = nil
		}
	}
	data.Message = msgID
	b.Last[aID] = data

	data.Sync.Unlock()
}

func ReplyCreate(s *discordgo.Session, channelID string, author string, content string, cc chan int) {
	var messageID string = ""
	var largest = 1
	var changed = false
	defer close(cc)

	for true {
		var cooldown = time.Second * 4;
		if !changed {
			cooldown = time.Second * 256
		}
		select {
		case <-time.After(cooldown):
			if !changed {
				continue
			}

			var msg *discordgo.Message
			var err error

			if messageID == "" {
				msg, err = s.ChannelMessageSend(
					channelID, 
					fmt.Sprintf("(x%d)[%s] %s", largest, author, content),
				)
			} else {
				msg, err = s.ChannelMessageEdit(
					channelID,
					messageID,
					fmt.Sprintf("(x%d)[%s] %s", largest, author, content),
				)
			}
			if err == nil {
				messageID = (*msg).ID
				fmt.Fprintf(os.Stdout, "Updating [%s] x%d\n", author, largest)
				changed = false
			} else {
				s.ChannelMessageDelete(channelID, messageID)
				messageID = ""
				fmt.Fprintf(os.Stderr, "Failed to send/update message: %s\n", err)
			}
		case counter, ok := <-cc:
			if !ok || counter < 0 {
				break
			}
			if largest >= counter {
				continue
			}
			largest = counter
			changed = true

			if messageID == "" {
				msg, err := s.ChannelMessageSend(
					channelID, 
					fmt.Sprintf("(x%d)[%s] %s", counter, author, content),
				)
				if err == nil {
					messageID = (*msg).ID
				} else {
					s.ChannelMessageDelete(channelID, messageID)
					messageID = ""
					fmt.Fprintf(os.Stderr, "Failed to send message: %s\n", err)
				}
			}
		}
	}
}

