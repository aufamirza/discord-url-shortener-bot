package persistence

import "discord-url-shortener-bot/persistence/localFile"

//interface for a link store to support multiple backends in future
type URLStore interface {
	Add(URL string) string
	Get(id string) string
}

func New() (error, URLStore) {
	//only support local file backend at the moment
	return localFile.New()
}
