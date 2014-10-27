package main

import (
	"fmt"
	"time"
)

type ArgStatus struct{}

// responds to status checks
func (n *Node) NodeStatus(arg *ArgStatus, statusReply *bool) error {
	*statusReply = true
	return nil
}

// checks status at neighbour nodes and reacts if other node is not ok
func checkStatus(n *Node, interval time.Duration) {
	var statusReply bool
	arg := &ArgStatus{}
	for {
		fmt.Printf("node." + n.Id + ": ")
		isAlive := n.nodeStatus()
		if isAlive == true {
			fmt.Printf("ok\n")
		} else {
			fmt.Printf("not ok\n")

			//if n.nodeId = node.predecessor {	// predecessor unavailble
			//blockRemoteAccess("pred", "node")
			//reconnectRing(n.predecessor.predecessor)
			//moveData() // stores data from previous predecessor
			//replicateData("node", n.predecessor, "succ")  // replicates data to new predecessor
			//allowRemoteAccess("pred", "node")
			//	} else {		 		// successor unavailble
			//blockRemoteAccess("node", "succ")
			reconnectRing(n.Successor.Successor)
			replicateData("succ", n.Successor, "node") // restores lost data to new sucessor
			replicateData("node", n.Successor, "pred") // replicates own data to new succ
			//allowRemoteAccess("node", "succ")
			//	}
		}
		time.Sleep(interval * time.Millisecond)
	}
}

// replicates or restores replicated data
func replicateData(oldStorageSpace string, newNode *Node, newStorageSpace string) {

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

	reader := bufio.NewReader(storageFile)
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
			contEOF = true
		} else {

			arg = &ArgStorage{key, data, newStorageSpace}
			dataStored := n.StoreData(arg, dataStored)
			if dataStored != true {
				fmt.Printf("Failed to store key: %s, data: %s, at node: %", key, data, n.nodeId)
			}
		}
	}
}

// moves replicated data to own storage when predeccesssor is lost
func moveData() {

	oldStorageFile, err := os.Open(predData.txt)
	if err != nil {
		log.Fatal(err)
	}
	defer oldStorageFile.Close()
	reader := bufio.NewReader(oldStorageFile)

	newStorageFile, err := os.Open(nodeData.txt)
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
		numbytes, err := newStorageFile.WriteAt([]byte(line), int64(newStorageFileSize))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d bytes written\n", numbytes)
	}
	oldStorageFile.Close()
	newStorageFile.Close()
}
