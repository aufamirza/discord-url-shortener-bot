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

//type LocalFile struct {
//	filePath string
//	uRLs     map[string]string
//}
type addRequest struct {
	//the long URL
	URL string
	//the channel to return the ID to
	iDChan chan string
}

type LocalFile struct {
	filePath       string
	uRLs           map[string]string
	addRequestChan chan addRequest
}

func New() (error, LocalFile) {
	//filename is currently hardcoded
	//TODO parameterize
	const filePath = "data.json"
	localFile := LocalFile{
		filePath:       filePath,
		uRLs:           map[string]string{},
		addRequestChan: make(chan addRequest),
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

		err = json.Unmarshal(bytes, &localFile.uRLs)
		if err != nil {
			return err, LocalFile{}
		}

		log.Println(fmt.Sprintf("loaded %v", filePath))
	}

	//start the goroutine to handle map writes
	go loop(localFile)
	return nil, localFile
}

//match the interface
func (localFile LocalFile) Get(id string) string {
	return localFile.uRLs[id]
}

//match the interface
func (localFile LocalFile) Add(URL string) string {
	iDChan := make(chan string)
	localFile.addRequestChan <- addRequest{
		URL:    URL,
		iDChan: iDChan,
	}
	return <-iDChan
}

//this single goroutine handles writes on the map to prevent race conditions
func loop(localFile LocalFile) {
	for addRequest := range localFile.addRequestChan {
		id := generateID(localFile)
		localFile.uRLs[id] = addRequest.URL
		addRequest.iDChan <- id
		//update the database file async because this could be slow and it is not mission critical the file is accurate
		go updateFile(localFile)
	}
}

func generateID(localFile LocalFile) string {
	const randomLength = 2
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	randomBuilder := strings.Builder{}
	for i := 0; i < randomLength; i++ {
		randomBuilder.WriteRune(chars[rand.Intn(len(chars))])
	}
	index := len(localFile.uRLs)
	idString := fmt.Sprintf("%v%v", index, randomBuilder.String())
	log.Println(idString)
	id := base64.StdEncoding.EncodeToString([]byte(idString))
	return id
}

//persist in memory map to object
func updateFile(localFile LocalFile) {
	bytes, err := json.Marshal(localFile.uRLs)
	if err != nil {
		log.Println(fmt.Sprintf("error: %v", err))
	}
	err = ioutil.WriteFile(localFile.filePath, bytes, 0644)
	if err != nil {
		log.Println(fmt.Sprintf("error: %v", err))
	}
}
