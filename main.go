package main

import (
	"discord-url-shortener-bot/bot"
	"discord-url-shortener-bot/persistence"
	"discord-url-shortener-bot/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//set the ENV var to read for the Discord bot token
	const tokenEnvVar = "DISCORD_BOT_TOKEN"
	var stop = make(chan os.Signal)
	//get the ENV var
	token := os.Getenv(tokenEnvVar)

	//if ENV var wasn't set then throw error
	if token == "" {
		log.Fatal(fmt.Sprintf("error: could not find env var $%v", tokenEnvVar))
	}

	err, URLStore := persistence.New()
	if err != nil {
		log.Fatal(fmt.Sprintf("error: could not find env var $%v", tokenEnvVar))
	}

	//start the server to serve redirect URL's
	go server.Start(stop, URLStore)
	go bot.Start(stop, token, URLStore)

	//make channel to listen to OS signals
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	//block until OS exit signal received and then send it to goroutines
	stop <- <-sc
}
