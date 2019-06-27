package localFileBackend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
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
	IDChan chan string
}

type localFile struct {
	filePath       string
	uRLs           map[string]string
	addRequestChan chan addRequest
}

func New() (error, localFile) {
	//filename is currently hardcoded
	//TODO parameterize
	const filePath = "data.json"
	newLocalFile := localFile{
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
			return err, localFile{}
		}

		err = json.Unmarshal(bytes, &newLocalFile.uRLs)
		if err != nil {
			return err, localFile{}
		}

		log.Println(fmt.Sprintf("loaded %v", filePath))
	}

	//start the goroutine to handle map writes
	go loop(newLocalFile)
	return nil, newLocalFile
}

//match the interface
func (localFile localFile) Get(id string) string {
	return localFile.uRLs[id]
}

//match the interface
func (localFile localFile) Add(uRL string) string {
	iDChan := make(chan string)
	localFile.addRequestChan <- addRequest{
		URL:    uRL,
		IDChan: iDChan,
	}
	return <-iDChan
}

//this single goroutine handles writes on the map to prevent race conditions
func loop(localFile localFile) {
	for addRequest := range localFile.addRequestChan {
		id := generateID(localFile)
		localFile.uRLs[id] = addRequest.URL
		addRequest.IDChan <- id
		//update the database file async because this could be slow and it is not mission critical the file is accurate
		go updateFile(localFile)
	}
}

//ID must be unique
//although the redirect URL is not treated as confidential, a slight randomness is applied to prevent guessing being totally trivial
//TODO
//convert the index to a string
//split the index
//convert the chars to ints using atoi
//append the matching runes
func generateID(localFile localFile) string {
	const randomLength = 2
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	randomBuilder := strings.Builder{}
	for i := 0; i < randomLength; i++ {
		randomBuilder.WriteRune(chars[rand.Intn(len(chars))])
	}
	randomBuilder.WriteString(strconv.Itoa(len(localFile.uRLs)))
	for i := 0; i < randomLength; i++ {
		randomBuilder.WriteRune(chars[rand.Intn(len(chars))])
	}
	idString := fmt.Sprintf("%v", randomBuilder.String())
	return idString
}

//persist in memory map to object
func updateFile(localFile localFile) {
	bytes, err := json.Marshal(localFile.uRLs)
	if err != nil {
		log.Println(fmt.Sprintf("error: %v", err))
	}
	err = ioutil.WriteFile(localFile.filePath, bytes, 0644)
	if err != nil {
		log.Println(fmt.Sprintf("error: %v", err))
	}
}
