package main

import (
	"bufio"
	"fmt"
	"os"
	"encoding/hex"
	"strconv"
)


type DHTnode struct {
    nodeId string
    nodeIdByte []byte
    nodeIp string
    nodePort string
    storagePath string
    storageContent string
    fingers []*Finger
    predecessor *DHTnode
    successor *DHTnode
    predecessors []*NeighbourNode  // larger index number = further away in ring
    successors []*NeighbourNode    // same as above
    joinViaIp string
    joinViaPort string
}

type NeighbourNode struct {
	nodeId string
    nodeIdByte []byte
    nodeIp string
    nodePort string
    //storagePath string
    //storageContent string
}

type Finger struct {
	key string
	nodeId string
	nodeIdByte []byte
	nodeIp string
	nodePort string	
}




func main() {

	var nodeList []*DHTnode 
	var firstNode *DHTnode

	wantedNodes := 3

	for i := 0; i < wantedNodes; i++ {
	//for {
	
		/*
		fmt.Println("Press Enter to create a node")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			{break}
		}
		*/
		
		/*
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		*/


		//port := (i*10) + 1111
		port := (i*1) + 1111
		newNode := createNode(port)
		nodeList = append(nodeList, newNode)
		nodesCreated := len(nodeList)
		if nodesCreated == 1 {
			firstNode = newNode
		}

		//fmt.Printf("\nNodes created: %d\n", nodesCreated)
		//for i := 0; i < (nodesCreated); i++ {
		//	fmt.Printf("Node %d: %s\n", i+1, nodeList[i].nodeId)
		//}
		//fmt.Println()
		if nodesCreated > 1 {

			firstNode.addToRing(newNode)
			
			if nodesCreated == 2 {	
				newNode.updateAllFingerTables()
			} else {
				newNode.updateIncorrectFingers()
			}
			
		}
		
	}


	//firstNode.updateFingerTables()

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



func (self *DHTnode) updateIncorrectFingers() {

	start := self
	newNode := self

	for start != self.successor { 
		for i:= 0; i < 160; i++ {

			if self.fingers[i].key >= newNode.nodeId {

				responsibleNode := self.ringLookup(self.fingers[i].key)		
				self.fingers[i].nodeId = responsibleNode.nodeId[:len(responsibleNode.nodeId)]		
			
			}
		}
		self = self.successor
	}
}


func (self *DHTnode) updateAllFingerTables() {  // updates all fingers in fingerTables of all nodes, starts with self

	start := self

	for start != self.successor { 
		for i:= 0; i < 160; i++ {
			responsibleNode := self.ringLookup(self.fingers[i].key)		
			self.fingers[i].nodeId = responsibleNode.nodeId[:len(responsibleNode.nodeId)]		
		}
		self = self.successor
	}

	// filling finger table for last node before starting node
	for i:= 0; i < 160; i++ {
		responsibleNode := self.ringLookup(self.fingers[i].key)		
		self.fingers[i].nodeId = responsibleNode.nodeId[:len(responsibleNode.nodeId)]		
	}

}

func (self *DHTnode) fingerLookup(hashedKey string) *DHTnode {

	targetNodeId := ""
	responsibleNode := self
	key := []byte(hashedKey)
	ownId := []byte(self.nodeId)
	nextNodeId := []byte(self.successor.nodeId)
	if (between(ownId, nextNodeId, key)) {  // starting node is responsible for key
		
		return self

	} else { 

		// deciding finger to use by iteration, replace with better algoritm???
		fingerFound := false
		i := 0
		for (!fingerFound && i < 160) {
			fingerA := []byte(self.fingers[i].nodeId)
			fingerB := []byte(self.fingers[i+1].nodeId)
			if (between(fingerA, fingerB, key)) {
				targetNodeId = self.fingers[i].nodeId
				fingerFound = true
			} else {
				i++
			}
		}
		if (!fingerFound) {
			targetNodeId = self.fingers[159].nodeId
		}

		// traversing ring clockwise instead of send request directly via IP of node
		for (!(self.successor.nodeId == targetNodeId)) {
			self = self.successor
		}
		
		// recursive request to closest node pointed to by finger
		responsibleNode = self.fingerLookup(hashedKey)
		return responsibleNode
	}
	
}



func (self *DHTnode) ringLookup(hashedKey string) *DHTnode{
	
	nodeFound := false
    key := []byte(hashedKey)

	for nodeFound == false { 
		
		id1 := []byte(self.nodeId)
    	id2 := []byte(self.successors[0].nodeId)
		
    	if (between(id1, id2, key)) {
      		nodeFound = true
    	} else {
	      	self = self.successor
    	}
    }
    return self
}


func createNode(port int) *DHTnode{

	//scanner := bufio.NewScanner(os.Stdin)

	/*
	fmt.Printf("IP for new node (localhost): ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		{break}
	}
	nodeIp := scanner.Text()
	*/
	nodeIp := "localhost"

	/*
	fmt.Printf("\nPort for new node: ")
	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		{break}
	}
	nodePort := scanner.Text()
	*/

	portString := strconv.Itoa(port)
	nodePort := portString

	/*
	fmt.Printf("Join via IP: ")
	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		{break}
	}
	joinViaIp := scanner.Text()
	*/
	joinViaIp := "localhost"

	/*
	fmt.Printf("Join via port: ")
	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		{break}
	}
	joinViaPort := scanner.Text()
	*/
	joinViaPort := "1111"

	node := makeDHTNode(nodeIp, nodePort, joinViaIp, joinViaPort)

	return node
}


func makeDHTNode(nodeIp string, nodePort string, joinViaIp string, joinViaPort string) *DHTnode {

	nodeIdStr := sha1hash(nodeIp+nodePort)
    nodeIdByte, _ := hex.DecodeString(nodeIdStr)
   

    node := &DHTnode { nodeIdStr, nodeIdByte, nodeIp, nodePort, "", "", nil, nil, nil, nil, nil, joinViaIp, joinViaPort}
    
    fingersWanted := 160
    for i := 0; i < fingersWanted; i++ {
    	fingerNumber := i+1
		newFingerKey := calcFinger(node.nodeIdByte, fingerNumber, 160)
    	newFinger := &Finger{newFingerKey, "", nil, "", ""}
    	node.fingers = append(node.fingers, newFinger)
    }

	newPredecessor0 := &NeighbourNode{"", nil, "", ""}
	newPredecessor1 := &NeighbourNode{"", nil, "", ""}
    newSuccessor0 := &NeighbourNode{"", nil, "", ""}
    newSuccessor1 := &NeighbourNode{"", nil, "", ""}
    node.predecessors = append(node.predecessors, newPredecessor0)
    node.predecessors = append(node.predecessors, newPredecessor1)
    node.successors = append(node.successors, newSuccessor0)
    node.successors = append(node.successors, newSuccessor1)

    return node
}

 

func (self *DHTnode) addToRing(node *DHTnode){

	/*
	// instead of traverings all nodes from self until finding point of insertion, 
	//fingers of existing nodes should be used
	*/

	if self.successor == nil {  // new node connects to a single node, forming a ring of two nodes

	    self.successor = node
	    node.predecessor = self
		self.successors[0].nodeId = node.nodeId[:len(node.nodeId)]
	    self.successors[0].nodeIp = node.nodeIp[:len(node.nodeIp)] 
	    self.successors[0].nodePort = node.nodePort[:len(node.nodePort)] 
		self.successors[1].nodeId = self.nodeId[:len(self.nodeId)]
	    self.successors[1].nodeIp = self.nodeIp[:len(self.nodeIp)]
	    self.successors[1].nodePort = self.nodePort[:len(self.nodePort)]
		node.predecessors[0].nodeId = self.nodeId[:len(self.nodeId)]
	    node.predecessors[0].nodeIp = self.nodeIp[:len(self.nodeIp)] 
	    node.predecessors[0].nodePort = self.nodePort[:len(self.nodePort)] 
		node.predecessors[1].nodeId = node.nodeId[:len(node.nodeId)]
	    node.predecessors[1].nodeIp = node.nodeIp[:len(node.nodeIp)]
	    node.predecessors[1].nodePort = node.nodePort[:len(node.nodePort)]

	    node.successor = self
	    self.predecessor = node
	    node.successors[0].nodeId = self.nodeId[:len(self.nodeId)]
	    node.successors[0].nodeIp = self.nodeIp[:len(self.nodeIp)] 
	    node.successors[0].nodePort = self.nodePort[:len(self.nodePort)] 
	    node.successors[1].nodeId = node.nodeId[:len(node.nodeId)]
	    node.successors[1].nodeIp = node.nodeIp[:len(node.nodeIp)]
	    node.successors[1].nodePort = node.nodePort[:len(node.nodePort)]
		self.predecessors[0].nodeId = node.nodeId[:len(node.nodeId)]
	    self.predecessors[0].nodeIp = node.nodeIp[:len(node.nodeIp)] 
	    self.predecessors[0].nodePort = node.nodePort[:len(node.nodePort)] 
		self.predecessors[1].nodeId = self.nodeId[:len(self.nodeId)]
	    self.predecessors[1].nodeIp = self.nodeIp[:len(self.nodeIp)] 
	    self.predecessors[1].nodePort = self.nodePort[:len(self.nodePort)]  

	} else {

    	for(!between([]byte(self.nodeId), []byte(self.successors[0].nodeId), []byte(node.nodeId))) {
    
      		self = self.successor
		
		}

	    if self.successors[1].nodeId == self.nodeId {		// new node connects to a ring of two nodes

	    	node.successor = self.successor
	    	node.successor.predecessor = node
		   	node.successors[0].nodeId = self.successors[0].nodeId[:len(self.successors[0].nodeId)]
		    node.successors[0].nodeIp = self.successors[0].nodeIp[:len(self.successors[0].nodeIp)] 
		    node.successors[0].nodePort = self.successors[0].nodePort[:len(self.successors[0].nodePort)] 
		    node.successors[1].nodeId = self.successors[1].nodeId[:len(self.successors[1].nodeId)]
		    node.successors[1].nodeIp = self.successors[1].nodeIp[:len(self.successors[1].nodeIp)] 
		    node.successors[1].nodePort = self.successors[1].nodePort[:len(self.successors[1].nodePort)] 
		    node.successor.predecessors[0].nodeId = node.nodeId[:len(self.nodeId)]
		    node.successor.predecessors[0].nodeIp = node.nodeIp[:len(self.nodeIp)] 
		    node.successor.predecessors[0].nodePort = node.nodePort[:len(self.nodePort)] 
		    node.successor.predecessors[1].nodeId = self.nodeId[:len(self.nodeId)]
		    node.successor.predecessors[1].nodeIp = self.nodeIp[:len(self.nodeIp)]
		    node.successor.predecessors[1].nodePort = self.nodePort[:len(self.nodePort)] 
		    node.successor.successors[1].nodeId = node.nodeId[:len(node.nodeId)]
		    node.successor.successors[1].nodeIp = node.nodeIp[:len(node.nodeIp)]
		    node.successor.successors[1].nodePort = node.nodePort[:len(node.nodePort)]

		    self.successor = node
		    node.predecessor = self
			self.successors[0].nodeId = node.nodeId[:len(node.nodeId)]
		    self.successors[0].nodeIp = node.nodeIp[:len(node.nodeIp)] 
		    self.successors[0].nodePort = node.nodePort[:len(node.nodePort)] 
		    self.successors[1].nodeId = node.successors[0].nodeId[:len(node.successors[0].nodeId)]
		    self.successors[1].nodeIp = node.successors[0].nodeIp[:len(node.successors[0].nodeIp)] 
		    self.successors[1].nodePort = node.successors[0].nodePort[:len(node.successors[0].nodePort)] 
			node.predecessors[0].nodeId = self.nodeId[:len(self.nodeId)]
		    node.predecessors[0].nodeIp = self.nodeIp[:len(self.nodeIp)] 
		    node.predecessors[0].nodePort = self.nodePort[:len(self.nodePort)] 
		    node.predecessors[1].nodeId = self.predecessors[0].nodeId[:len(node.predecessors[0].nodeId)]
		    node.predecessors[1].nodeIp = self.predecessors[0].nodeIp[:len(node.predecessors[0].nodeIp)]
		    node.predecessors[1].nodePort = self.predecessors[0].nodePort[:len(node.predecessors[0].nodePort)]
		    self.predecessors[1].nodeId = node.nodeId[:len(node.nodeId)]
		    self.predecessors[1].nodeIp = node.nodeIp[:len(node.nodeIp)] 
		    self.predecessors[1].nodePort = node.nodePort[:len(node.nodePort)]
		    

	    } else {	// new node connects to a ring of at least three nodes
		   	
		   	node.successor = self.successor
		   	node.successor.predecessor = node
		   	node.successors[0].nodeId = self.successors[0].nodeId[:len(self.successors[0].nodeId)]
		    node.successors[0].nodeIp = self.successors[0].nodeIp[:len(self.successors[0].nodeIp)] 
		    node.successors[0].nodePort = self.successors[0].nodePort[:len(self.successors[0].nodePort)] 
		    node.successors[1].nodeId = self.successors[1].nodeId[:len(self.successors[1].nodeId)]
		    node.successors[1].nodeIp = self.successors[1].nodeIp[:len(self.successors[1].nodeIp)] 
		    node.successors[1].nodePort = self.successors[1].nodePort[:len(self.successors[1].nodePort)] 
		    node.successor.predecessors[0].nodeId = node.nodeId[:len(self.nodeId)]
		    node.successor.predecessors[0].nodeIp = node.nodeIp[:len(self.nodeIp)] 
		    node.successor.predecessors[0].nodePort = node.nodePort[:len(self.nodePort)] 
		    node.successor.predecessors[1].nodeId = self.nodeId[:len(self.nodeId)]
		    node.successor.predecessors[1].nodeIp = self.nodeIp[:len(self.nodeIp)]
		    node.successor.predecessors[1].nodePort = self.nodePort[:len(self.nodePort)] 
		    node.successor.successor.predecessors[1].nodeId = node.nodeId[:len(self.nodeId)]
		    node.successor.successor.predecessors[1].nodeIp = node.nodeIp[:len(self.nodeIp)] 
		    node.successor.successor.predecessors[1].nodePort = node.nodePort[:len(self.nodePort)] 

		    self.successor = node
		    node.predecessor = self
		    self.predecessor.successors[1].nodeId = node.nodeId[:len(node.nodeId)]
		    self.predecessor.successors[1].nodeIp = node.nodeIp[:len(node.nodeIp)] 
		    self.predecessor.successors[1].nodePort = node.nodePort[:len(node.nodePort)] 
			self.successors[0].nodeId = node.nodeId[:len(node.nodeId)]
		    self.successors[0].nodeIp = node.nodeIp[:len(node.nodeIp)] 
		    self.successors[0].nodePort = node.nodePort[:len(node.nodePort)] 
		    self.successors[1].nodeId = node.successors[0].nodeId[:len(node.successors[0].nodeId)]
		    self.successors[1].nodeIp = node.successors[0].nodeIp[:len(node.successors[0].nodeIp)] 
		    self.successors[1].nodePort = node.successors[0].nodePort[:len(node.successors[0].nodePort)] 
			node.predecessors[0].nodeId = self.nodeId[:len(self.nodeId)]
		    node.predecessors[0].nodeIp = self.nodeIp[:len(self.nodeIp)] 
		    node.predecessors[0].nodePort = self.nodePort[:len(self.nodePort)] 
		    node.predecessors[1].nodeId = self.predecessors[0].nodeId[:len(node.predecessors[0].nodeId)]
		    node.predecessors[1].nodeIp = self.predecessors[0].nodeIp[:len(node.predecessors[0].nodeIp)]
		    node.predecessors[1].nodePort = self.predecessors[0].nodePort[:len(node.predecessors[0].nodePort)]

		}
  	}

  	//self.updateAllFingerTables()
}


