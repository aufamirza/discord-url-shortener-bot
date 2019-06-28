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
	const persistenceBackendTypeEnvVar = "DISCORD_BOT_PERSISTENCE_BACKEND_TYPE"
	//get the ENV var
	persistenceBackend := os.Getenv(persistenceBackendTypeEnvVar)
	//if ENV var wasn't set then throw error
	if persistenceBackend == "" {
		log.Fatal(fmt.Sprintf("error: could not find env var $%v", persistenceBackendTypeEnvVar))
	}

	const persistenceBackendSqlConnectionStringEnvVar = "DISCORD_BOT_PERSISTENCE_BACKEND_SQL_CONNECTION_STRING"
	persistenceBackendSqlHost := os.Getenv(persistenceBackendSqlConnectionStringEnvVar)
	if persistenceBackend == "SQL" {
		//set the ENV var to read
		if persistenceBackendSqlHost == "" {
			log.Fatal(fmt.Sprintf("error: persistence type SQL selected but could not find env var $%v", persistenceBackendSqlConnectionStringEnvVar))
		}
	}

	err, URLStore := persistence.New(persistence.BackendTypeLocalFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("error: could not find env var $%v", tokenEnvVar))
	}

	var stop = make(chan os.Signal)

	//start the server to serve redirect URL's
	go server.Start(stop, URLStore)
	go bot.Start(stop, token, host, URLStore)

	//make channel to listen to OS signals
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	//block until OS exit signal received and then send it to goroutines
	stop <- <-sc
}
