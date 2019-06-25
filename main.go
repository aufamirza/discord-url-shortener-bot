package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

var urlRegexp = regexp.MustCompile(`(https?://\S+\.\S+)`)

func main() {
	//set the ENV var to read for the Discord bot token
	const tokenEnvVar = "DISCORD_BOT_TOKEN"

	//get the ENV var
	token := os.Getenv(tokenEnvVar)

	//if ENV var wasn't set then throw error
	if token == "" {
		log.Fatal(fmt.Sprintf("error: could not find env var $%v", tokenEnvVar))
	}

	//create bot
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("error: could not create bot: %v", err)
	}

	//add message created handler
	//events can be found here https://discordapp.com/developers/docs/topics/gateway#event-names
	bot.AddHandler(messageCreate)

	//attempt to open the bot websocket connection
	err = bot.Open()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Println("started bot")

	//make channel to listen to OS signals
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	//block until OS exit signal received
	<-sc

	//close bot and don't handle closing errors
	_ = bot.Close()
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	//ignore messages from self
	if message.Author.ID == session.State.User.ID {
		return
	}

	//get all matched url's in the message
	urls := urlRegexp.FindAllString(message.Content, -1)

	for _, url := range urls {
		//echo url back
		_, err := session.ChannelMessageSend(message.ChannelID, url)
		if err != nil {
			log.Println(fmt.Sprintf("error %v", err))
		}
	}
}
