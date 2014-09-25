package main
    
import (
    "encoding/hex"
    "strconv"
    "fmt"
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

    nodeIp := "localhost"
    portString := strconv.Itoa(port)
    nodePort := portString

    joinViaIp := "localhost"
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

        if self.successors[1].nodeId == self.nodeId {       // new node connects to a ring of two nodes

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
            

        } else {    // new node connects to a ring of at least three nodes
            
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

