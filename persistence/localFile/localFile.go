package localFile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type LocalFile struct {
	filePath string
	links    map[string]string
}

func New() (error, LocalFile) {
	//filename is currently hardcoded
	//TODO parameterize
	const filePath = "data.json"
	localFile := LocalFile{
		filePath: filePath,
		links:    map[string]string{},
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

		err = json.Unmarshal(bytes, &localFile.links)
		if err != nil {
			return err, LocalFile{}
		}

		log.Println(fmt.Sprintf("loaded %v", filePath))
		return nil, localFile
	}

	return nil, localFile
}

func (localFile LocalFile) Get(id string) string {
	panic("unimplemented")
}

func (localFile LocalFile) Add(link string) string {
	//generate next ID
	id := generateID(localFile)

	//store long URL in map using generated ID as key
	localFile.links[id] = link
	return id
}

func generateID(localFile LocalFile) string {
	//mock ID out
	go updateFile(localFile)
	return strconv.Itoa(0)
}

//persist in memory map to object
func updateFile(localFile LocalFile) {
	bytes, err := json.Marshal(localFile.links)
	if err != nil {
		log.Println(fmt.Sprintf("error: %v", err))
	}
	err = ioutil.WriteFile(localFile.filePath, bytes, 0644)
	if err != nil {
		log.Println(fmt.Sprintf("error: %v", err))
	}
}
