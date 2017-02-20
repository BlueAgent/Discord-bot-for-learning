package main
//Made with heavy reference to discordgo examples

import (
	"github.com/bwmarrin/discordgo"
	"flag"
	"fmt"
	"os"
	"sync"
)

type LastData struct {
	Sync sync.Mutex
	Message string
	Counter int
	Reply string
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

	dg.AddHandler(b.messageCreate)

	//Open websocket and start listening
	err = dg.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open websocket: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Bot has started. Exit with CTRL-C.")
	//I have no idea how it works
	//Discordgo example says it keeps the program running
	<-make(chan struct{})
	return
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
				Message: "",
				Counter: 1,
				Reply: "",
			}
		}
		b.LastSync.Unlock()
	}

	if data.Message == m.Content {
		go s.ChannelMessageDelete(m.ChannelID, m.ID)
	}

	//Only one thread per user
	data.Sync.Lock()
	//data.Sync.Unlock() //Should get run after the function returns

	fmt.Println(m.Content)

	if data.Message == m.Content {
		fmt.Println("Same")
		go s.ChannelMessageDelete(m.ChannelID, m.ID)

		data.Counter++
		var msg *discordgo.Message
		var err error
		if data.Reply != "" {
			msg, err = s.ChannelMessageEdit(
				m.ChannelID,
				data.Reply,
				fmt.Sprintf("(x%d)[%s] %s", data.Counter, m.Author.Username, data.Message),
			)
			if err != nil {
				s.ChannelMessageDelete(m.ChannelID, data.Reply)
				data.Reply = ""
			}
		}
		if data.Reply == "" || err != nil {
			msg, err = s.ChannelMessageSend(
				m.ChannelID, 
				fmt.Sprintf("(x%d)[%s] %s", data.Counter, m.Author.Username, data.Message),
			)
			if err == nil {
				data.Reply = (*msg).ID
			} else {
				data.Reply = ""
				fmt.Fprintf(os.Stderr, "Failed to send new message: %s\n", err)
			}
		}
	} else {
		fmt.Println("Different")
		data.Counter = 1
		data.Reply = "" 
	}
	data.Message = m.Content
	data.Sync.Unlock()
	b.Last[aID] = data
	fmt.Println("Unlocked")
}
