package main

import ()

func (thisNode *Node) removeData (storageSpace string, unhashedKey ) {

	hashedKey := sha1hash(unhashedKey)
	arg := &ArgLookup{hashedKey}
	reply := nil
	err := thisNode.findSuccessor(arg, &reply)
	if err != nil {
		log.Fatal(err)
	}

	nodeToRemoveAt := reply.BasicNode

	err := nodeToRemoveAt.deleteDataRemote(storageSpace, hashedKey)
	if err != nil {
		log.Fatal(err)
	}


}

func (thisNode *Node) uploadData (unhashedKey string, data string) {

	hashedKey := sha1hash(unhashedKey)
	arg := &ArgLookup{hashedKey}
	reply := nil
	err := thisNode.findSuccessor(arg, &reply)
	if err != nil {
		log.Fatal(err)
	}

	nodeToStoreAt := reply.BasicNode

	err := nodeToStoreAt.storeDataRemote(hashedKey, data, "node")
	if err != nil {
		log.Fatal(err)
	}
}




// stores data at current node, can be called from another node
func (n *BasicNode) StoreData(arg *ArgStorage, dataStored *bool) error {

	key := arg.Key
	data := arg.Data
	storageSpace := arg.StorageSpace
	appendDataToStorage(key, data, storageSpace)
	if storageSpace == "node" {		
		replicateData("node", n.predeccessor, "node")
		replicateData("node", n.successor, "node")
	}

	*dataStored = true
	return nil
}

// used by StoreData()
func appendDataToStorage(key string, data string, storageSpace string) {

	_ = os.Chdir("..")
	_ = os.Chdir("..")
	_ = os.Chdir("storage")

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

	_ = os.Chdir("..")
	_ = os.Chdir("new_git")	
	_ = os.Chdir("src")	
}



// deletes key-data pair, can be called from another node
func (n *DHTnode) DeleteData (arg *ArgDeletion, dataDeleted bool) error {

	_ = os.Chdir("..")
	_ = os.Chdir("..")
	_ = os.Chdir("storage")

	storageSpace := arg.StorageSpace
	key := arg.Key

	currentFileName := ""
	oldFileName := ""
	if storageSpace == "node" {
		currentFileName = "nodeData.txt"
		oldFileName = "oldNodeData.txt"
	} else if storageSpace == "succ" {
		currentFileName = "succData.txt"
		oldFileName = "oldSuccData.txt"
	} else {
		currentFileName = "predData.txt"
		oldFileName = "oldPredData.txt"
	}

	err := os.Rename(currentFileName,oldFileName)
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

	_ = os.Chdir("..")
	_ = os.Chdir("new_git")	
	_ = os.Chdir("src")	

	dataDeleted = true
	return nil
}

// prints all key/data-pairs in one of the storage spaces of the node
func (n *DHTnode) ListStoredData (arg *ArgListing, dataListed *bool) error {
//func (n *DHTnode) ListStoredData(storageSpace string) {

	_ = os.Chdir("..")
	_ = os.Chdir("..")
	_ = os.Chdir("storage")

	filename := ""
	if arg.storageSpace == "node" {
		filename = "nodeData.txt"
	} else if arg.storageSpace == "succ" {
		filename = "succData.txt"
	} else {
		filename = "predData.txt"
	}

	storageFile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("failed to open nodeData.txt")
		log.Fatal(err)
	}
	defer storageFile.Close()

	reader := bufio.NewReader(storageFile)
	storageEOF := false
	fmt.Printf("\n\nFiles stored in %s space:\n", arg.storageSpace)
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

	_ = os.Chdir("..")
	_ = os.Chdir("new_git")	
	_ = os.Chdir("src")	

	*dataListed = true;
	return nil
}

// inits a folder (ip for unique name when on same computer) and files for storing keys-data pair if they dont exist
func (n *DHTnode) StorageInit() {

	_ = os.Chdir("..")
	_ = os.Chdir("..")
	
	folderName := "storage" + n.Ip
	CreateDir(folderName)
	_ = os.Chdir(folderName)

	CreateFile("succData.txt")
	CreateFile("nodeData.txt")
	CreateFile("predData.txt")

	_ = os.Chdir("..")
	_ = os.Chdir("new_git")	
	_ = os.Chdir("src")	
}

// replicates or restores replicated data
func replicateData(oldStorageSpace string, newNode *BasicNode, newStorageSpace string) {

	filename := ""
	if oldStorageSpace == "node" {
		filename = "nodeData.txt"
	} else if oldStorageSpace == "pred" {
		filename = "predData.txt"
	} else {
		filename = "succData.txt"
	}

	oldStorageFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer oldStorageFile.Close()

	reader := bufio.NewReader(oldStorageFile)
	storageEOF := false
	for !storageEOF {
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
			arg := &ArgStorage{string(key[:]), string(data[:]), newStorageSpace}
			var dataStored bool
			err = newNode.StoreData(arg, &dataStored)
			if err != nil {
				fmt.Printf("Failed to store key: %s, data: %s, at node: %", key, data, newNode.Id)
			}
		}
	}
}

// moves replicated data to own storage when predeccesssor is lost
func moveData() {

	oldStorageFile, err := os.Open("predData.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer oldStorageFile.Close()
	reader := bufio.NewReader(oldStorageFile)

	newStorageFile, err := os.Open("nodeData.txt")
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

		newStorageFileInfo, err := newStorageFile.Stat()
		if err != nil {
			log.Fatal(err)
		}
		newStorageFileSize := newStorageFileInfo.Size()
		numbytes, err := newStorageFile.WriteAt([]byte(stringLine), int64(newStorageFileSize))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d bytes written\n", numbytes)
	}
	oldStorageFile.Close()
	newStorageFile.Close()
}
