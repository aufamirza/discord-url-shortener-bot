package bot

import (
	"discord-url-shortener-bot/persistence"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"regexp"
	"strings"
)

//url match regex
var urlRegexp *regexp.Regexp = regexp.MustCompile(`(https?://\S+\.\S+)`)

//ignore list regex
var urlIgnoreRegex *regexp.Regexp = regexp.MustCompile(`(open.spotify.com)`)
var URLStore persistence.URLStore
var hostname string
var protocol string
var port string

func Start(stop chan os.Signal, token string, newProtocol string, newHost string, newPort string, newURLStore persistence.URLStore) {
	//configure hostname
	hostname = newHost
	//configure persistence
	URLStore = newURLStore
	//configure protocol
	protocol = newProtocol
	//configure port
	if newPort != "" && newPort != "80" && newPort != "443" {
		port = fmt.Sprintf(":%v", newPort)
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
	//ignore messages from self and bots
	if message.Author.ID == session.State.User.ID || message.Author.Bot {
		return
	}

	//get all matched URL's in the message
	urls := urlRegexp.FindAllString(message.Content, -1)
	if 0 < len(urls) {
		newMessage := message.Content
		for _, url := range urls {
			//url ignore list
			if ignore(url) {
				return
			}
			//persist URL
			id := URLStore.Add(url)
			newMessage = strings.Replace(newMessage, url, fmt.Sprintf("%v://%v%v/%v", protocol, hostname, port, id), -1)
		}
		_, err := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%v\n%v", message.Author.Mention(), newMessage))
		if err != nil {
			log.Println(fmt.Sprintf("error: %v", err))
		}
		if err != nil {
			log.Println(fmt.Sprintf("error: %v", err))
		}
		err = session.ChannelMessageDelete(message.ChannelID, message.ID)

	}
}

//match ignore list
func ignore(URL string) bool {
	return urlIgnoreRegex.MatchString(URL)
}
