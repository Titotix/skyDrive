package main

import (
	"fmt"
	"os"
	"io"
	"bufio"
	"bytes"
	"log"
	"strings"
)

// inits a folder and files for storing keys-data pair if they dont exist
func StorageInit() {

	CreateDir("storage")
	_ = os.Chdir("storage")

	CreateFile("succData.txt")
	CreateFile("nodeData.txt")
	CreateFile("predData.txt")

	_ = os.Chdir("..")	
}



// used by StorageInit() 
func CreateDir(dirname string) {
	err := os.Mkdir(dirname, 0666)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}
}

// always used by StorageInit(), can also be used if files are stored
func CreateFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}
 	defer file.Close()
}

// just for testing
func listStoredData(storageSpace string) {

	filename := ""
	if storageSpace == "node" {
		filename = "nodeData.txt"
	} else if storageSpace == "succ" {
		filename = "succData.txt"
	} else {
		filename = "predData.txt"
	}

	storageFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer storageFile.Close()

	reader := bufio.NewReader(storageFile)
	storageEOF := false
	fmt.Printf("\n\nFiles stored in %s space:\n", storageSpace)
	for (!storageEOF) {
		key_delim, err := reader.ReadBytes(',')
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		}
		key := bytes.TrimSuffix(key_delim, []byte(","))
		data, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		}
		if (len(data)) == 0 {
			storageEOF = true
		} else {
			fmt.Printf("key:%s\n", key)
			fmt.Printf("data:%s\n", data)
		}
	}
	storageFile.Close()
}