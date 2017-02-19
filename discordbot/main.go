package main
//Made with heavy reference to discordgo examples

import (
	"github.com/bwmarrin/discordgo"
	"flag"
	"fmt"
	"os"
)

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
	
	dg.AddHandler(messageCreate) 

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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Content)
}
