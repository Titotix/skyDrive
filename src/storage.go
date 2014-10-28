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

type ArgDeletion struct {
	storagesSpace string
	key string
}


type ArgStorage struct {
	key string
	data string
	StorageSpace string
}


// inits a folder and files for storing keys-data pair if they dont exist
func storageInit() {

	CreateDir("storage")
	_ = os.Chdir("storage")

	CreateFile("succData.txt")
	CreateFile("nodeData.txt")
	CreateFile("predData.txt")

	_ = os.Chdir("..")	
}

// deletes key-data pair, can be called from another node
// arg[0] is for node/succ/pred, arg[1] is for key 
func (n *Node) DeleteData (arg *ArgDeletion, dataDeleted *bool) error {

	storageSpace := ArgDeletion[0] 
	key := ArgDeletion[1]

	filename := ""
	if storageSpace == "node" {
		currentFilename = "nodeData.txt"
		oldFileName = "oldNodeData.txt"
	} else if storageSpace == "succ" {
		currentFilename = "succData.txt"
		oldFileName = "oldSuccData.txt"
	} else {
		currentFilename = "predData.txt"
		oldFileName = "oldPredData.txt"
	}

	err = os.Rename(currentFilename,oldFileName)
	if err != nil {
		log.Fatal(err)
	}


	oldStorageFile, err := os.Open(oldFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldStorageFile.Close()
	reader := bufio.NewReader(oldStorageFile)

	newStorageFile, err := os.Create(currentFileName)
	if err != nil {
		log.Fatal(err)
	}
 	defer newStorageFile.Close()


 	oldStorageEOF := false
 	for !oldStorageEOF {
		line, err := reader.ReadBytes('\n')
		if err != nil {
 			if err == io.EOF {
 				oldStorageEOF = true
 			} else {
				log.Fatal(err)
			}
		}

		stringLine := string(line[:])

		if !strings.Contains(stringLine, key) {
			newStorageFileInfo, err := newStorageFile.Stat()
 			if err != nil {
				log.Fatal(err)
			}
	 		newStorageFileSize := newStorageFileInfo.Size()
			numbytes, err := newStorageFile.WriteAt([]byte(line), int64(newStorageFileSize))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%d bytes written\n", numbytes)
		}
	}
	oldStorageFile.Close()
	newStorageFile.Close()

	err = os.Remove(oldFileName)
	if err != nil {
		log.Fatal(err)
	}

	*dataDeleted = true
	return nil
}

// stores data at current node, can be called from another node
// arg[0] is key, arg[1] is for data 
func (n *Node) StoreData(arg *ArgStorage, dataStored *bool) error {

	key := ArgStorage[0]
	data := ArgStorage[1]
	storageSpace := ArgStorage[2]
	appendDataToStorage(key, data, storageSpace)
	*dataStored = true
	return nil
}

// used by storeData()
func appendDataToStorage(key string, data string, storageSpace string) {

	filename := ""
	if storageSpace == "node" {
		filename = "nodeData.txt"
	} else if storageSpace == "succ" {
		filename = "succData.txt"
	} else {
		filename = "predData.txt"
	}

	storageFile, err := os.OpenFile(filename, os.O_APPEND, 0666) 
	if err != nil {
		log.Fatal(err)
	}
	defer storageFile.Close()

	storageFileInfo, _ := storageFile.Stat()
	lastchar := storageFileInfo.Size()

	line := key + "," + data + "\r\n"
	numbytes, _ := storageFile.WriteAt([]byte(line), int64(lastchar))
	storageFile.Close()
	fmt.Printf("%d bytes written to contents file\n", numbytes)	
}

// used by storageInit() 
func CreateDir(dirname string) {
	err := os.Mkdir(dirname, 0666)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}
}

// always used by storageInit(), can also be used if files are stored
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
			contEOF = true
		} else {
			fmt.Printf("key:%s\n", key)
			fmt.Printf("data:%s\n", data)
		}
	}
	storageFile.Close()
}