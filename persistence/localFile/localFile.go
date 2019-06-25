package localFile

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
)

type LocalFile struct {
	filePath string
	URLs     map[string]string
}

func New() (error, LocalFile) {
	//filename is currently hardcoded
	//TODO parameterize
	const filePath = "data.json"
	localFile := LocalFile{
		filePath: filePath,
		URLs:     map[string]string{},
	}

	//check if file exists
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		//file exists
		//read local database file
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err, LocalFile{}
		}

		err = json.Unmarshal(bytes, &localFile.URLs)
		if err != nil {
			return err, LocalFile{}
		}

		log.Println(fmt.Sprintf("loaded %v", filePath))
		return nil, localFile
	}

	return nil, localFile
}

func (localFile LocalFile) Get(id string) string {
	return localFile.URLs[id]
}

func (localFile LocalFile) Add(link string) string {
	//generate next ID
	id := generateID(localFile)

	//store long URL in map using generated ID as key
	localFile.URLs[id] = link
	go updateFile(localFile)
	return id
}

func generateID(localFile LocalFile) string {
	const RANDOM_LENGTH = 2
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	randomBuilder := strings.Builder{}
	for i := 0; i < RANDOM_LENGTH; i++ {
		randomBuilder.WriteRune(chars[rand.Intn(len(chars))])
	}
	index := len(localFile.URLs)
	idString := fmt.Sprintf("%v%v", index, randomBuilder.String())
	log.Println(idString)
	id := base64.StdEncoding.EncodeToString([]byte(idString))
	return id
}

//persist in memory map to object
func updateFile(localFile LocalFile) {
	bytes, err := json.Marshal(localFile.URLs)
	if err != nil {
		log.Println(fmt.Sprintf("error: %v", err))
	}
	err = ioutil.WriteFile(localFile.filePath, bytes, 0644)
	if err != nil {
		log.Println(fmt.Sprintf("error: %v", err))
	}
}
