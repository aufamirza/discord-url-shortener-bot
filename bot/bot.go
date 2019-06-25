package bot

import (
	"discord-url-shortener-bot/persistence"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"regexp"
)

var urlRegexp *regexp.Regexp = regexp.MustCompile(`(https?://\S+\.\S+)`)
var linkStore persistence.LinkStore

func Start(stop chan os.Signal, token string) {
	//set error to empty to prevent locally scoping linkStore
	var err error = nil

	//configure persistence
	err, linkStore = persistence.New()
	if err != nil {
		log.Fatalf("error: %v", err)
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
		//persist URL
		id := linkStore.Add(url)
		_, err := session.ChannelMessageSend(message.ChannelID, id)
		if err != nil {
			log.Println(fmt.Sprintf("error: %v", err))
		}
	}
}
