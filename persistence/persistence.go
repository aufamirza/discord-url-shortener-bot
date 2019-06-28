package persistence

import (
	"discord-url-shortener-bot/persistence/localFileBackend"
	"fmt"
	"log"
)

type BackendType int

const BackendTypeLocalFile = 0
const BackendTypeSQL = 1

//interface for a link store to support multiple backends in future
type URLStore interface {
	Add(URL string) string
	Get(id string) string
}

func New() (error, URLStore) {
	err, localFile := localFileBackend.New()
	if err != nil {
		log.Fatal(fmt.Sprintf("error: %v", err))
	}
	return nil, localFile
}
