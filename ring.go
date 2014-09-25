package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	var nodeList []*DHTnode 
	var firstNode *DHTnode

	wantedNodes := 10

	for i := 0; i < wantedNodes; i++ {
		port := (i*1) + 1111
		newNode := createNode(port)
		nodeList = append(nodeList, newNode)
		nodesCreated := len(nodeList)
		if nodesCreated == 1 {
			firstNode = newNode
		}

		if nodesCreated > 1 {

			firstNode.addToRing(newNode)
			
			if nodesCreated == 2 {	
				newNode.updateAllFingerTables()
			} else {
				newNode.updateIncorrectFingers()
			}
			
		}
		
	}

	firstNode.printRing()
	
	fmt.Printf("\nSearch for key: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		{break}
	}
	testKey := scanner.Text()
	testHash := sha1hash(testKey)
	fmt.Printf("Key hashed to: %s\n\n", testHash)
	fmt.Printf("ringLookup, nodeId: %s\n", firstNode.fingerLookup(testHash).nodeId)

}