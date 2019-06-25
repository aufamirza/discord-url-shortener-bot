package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"regexp"
)

var urlRegexp = regexp.MustCompile(`(https?://\S+\.\S+)`)

func Start(stop chan os.Signal, token string) {
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

	//block until stop signal received
	<-stop
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	//ignore messages from self
	if message.Author.ID == session.State.User.ID {
		return
	}

	//get all matched URL's in the message
	urls := urlRegexp.FindAllString(message.Content, -1)

	for _, url := range urls {
		//echo url back
		_, err := session.ChannelMessageSend(message.ChannelID, url)
		if err != nil {
			log.Println(fmt.Sprintf("error %v", err))
		}
	}
}
