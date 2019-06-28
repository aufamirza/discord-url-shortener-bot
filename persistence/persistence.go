package persistence

import (
	"discord-url-shortener-bot/persistence/localFileBackend"
	"discord-url-shortener-bot/persistence/sqlBackend"
	"errors"
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

func New(backend BackendType, location string) (error, URLStore) {
	switch backend {
	case BackendTypeLocalFile:
		{
			err, localFile := localFileBackend.New(location)
			if err != nil {
				log.Fatal(fmt.Sprintf("error: %v", err))
			}
			return nil, localFile
		}
	case BackendTypeSQL:
		{
			err, sql := sqlBackend.New(location)
			if err != nil {
				log.Fatal(fmt.Sprintf("error: %v", err))
			}
			return nil, sql
		}
	default:
		{
			return errors.New("persistence backend not recognised"), nil
		}
	}
}
