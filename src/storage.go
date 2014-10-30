package main

import (
	"os"
	"log"
)



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

