package main

import (
	"discord-url-shortener-bot/bot"
	"discord-url-shortener-bot/persistence"
	"discord-url-shortener-bot/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	//set the ENV var to read
	const tokenEnvVar = "DISCORD_BOT_TOKEN"
	//get the ENV var
	token := os.Getenv(tokenEnvVar)
	//if ENV var wasn't set then throw error
	if token == "" {
		log.Fatal(fmt.Sprintf("error: could not find env var $%v", tokenEnvVar))
	}

	//set the ENV var to read
	const hostEnvVar = "DISCORD_BOT_HOST"
	//get the ENV var
	host := os.Getenv(hostEnvVar)
	//if ENV var wasn't set then throw error
	if host == "" {
		log.Fatal(fmt.Sprintf("error: could not find env var $%v", hostEnvVar))
	}

	//set the ENV var to read
	const postEnvVar = "DISCORD_BOT_PORT"
	//get the ENV var
	port := os.Getenv(postEnvVar)

	//set the ENV var to read
	const protocolEnvVar = "DISCORD_BOT_PROTOCOL"
	//get the ENV var
	protocol := strings.ToLower(os.Getenv(protocolEnvVar))
	//if ENV var wasn't set then throw error
	if protocol == "" {
		log.Fatal(fmt.Sprintf("error: could not find env var $%v", protocolEnvVar))
	} else if !(protocol == "https" || protocol == "http") {
		log.Fatal("error: protocol must be http or https")
	}

	err, URLStore := persistence.New()
	if err != nil {
		log.Fatal(fmt.Sprintf("error: could not find env var $%v", tokenEnvVar))
	}

	var stop = make(chan os.Signal)

	//start the server to serve redirect URL's
	go server.Start(stop, port, URLStore)
	go bot.Start(stop, token, protocol, host, port, URLStore)

	//make channel to listen to OS signals
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	//block until OS exit signal received and then send it to goroutines
	stop <- <-sc
}
