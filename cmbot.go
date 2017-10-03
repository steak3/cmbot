package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	BotID          string
	isGreetEnabled bool
	greetMessage   string
	session        discordgo.Session
)

type Data struct {
	IsGreetEnabled bool     `json:"isGreetEnabled"`
	GreetMessage   string   `json:"greetMessage"`
	Repeats        []string `json:"repeats"`
}

//TODO: repeats - docker - organize code
// ok 2 3
func main() {

	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var data Data
	json.Unmarshal(file, &data)
	isGreetEnabled = data.IsGreetEnabled
	greetMessage = data.GreetMessage

	// Create a new Discord session using the provided bot token from env vars. Test
	session, err := discordgo.New("Bot " + os.Getenv("CMBOT_TOKEN"))

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := session.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	// Store the account ID for later use.
	BotID = u.ID

	// Register callbacks for the events.
	session.AddHandler(messageCreate)
	session.AddHandler(guildMemberAdd)

	// Open the websocket and begin listening.
	err = session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

// When someone joins the discord server
func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if isGreetEnabled {
		format := strings.Replace(greetMessage, "%user%", "<@"+m.Member.User.ID+">", -1)
		_, _ = s.ChannelMessageSend(m.GuildID, format)
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	var confirmMessage string
	split := strings.SplitN(m.Content, " ", 2)
	//fmt.Println(split)

	switch split[0] {
	case ".repeat":
		repeat(split)
	case ".greet":
		confirmMessage = greet(split)
	}

	// Send nothing is confirmMessage is null
	_, _ = s.ChannelMessageSend(m.ChannelID, confirmMessage)

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func repeat(cmd []string) {
	if len(cmd) > 1 {
		fmt.Print("yeah repeat! -> ")
		fmt.Println(cmd)
	}
}

func greet(cmd []string) string {
	//.greet my_welcom_message
	if len(cmd) > 1 {
		greetMessage = cmd[1]
		return "New greet message set!"
		// .greet
	} else {
		isGreetEnabled = !isGreetEnabled
		if isGreetEnabled {
			return "Welcome message enabled"
		} else {
			return "Welcome message disabled"
		}
	}
}
