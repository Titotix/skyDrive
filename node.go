package main

import (
	"encoding/hex"
)

//Could be nice to define minimal node struct
type Node struct {
	Id     string
	IdByte []byte
	Ip     string
	Port   string
}

type DHTnode struct {
	Node
	storagePath    string
	storageContent string
	Fingers        []*Finger
	Predecessor    *DHTnode
	Successor      *DHTnode
	Predecessors   []*NeighbourNode // larger index number = further away in ring
	Successors     []*NeighbourNode // same as above
	joinViaIp      string
	joinViaPort    string
}

type NeighbourNode struct {
	Id       string
	IdByte   []byte
	NodeIp   string
	NodePort string
	//storagePath string
	//storageContent string

}

func (self *DHTnode) updateIncorrectFingers() {

	start := self
	newNode := self

	for start != self.Successor {
		for i := 0; i < 160; i++ {

			if self.Fingers[i].key >= newNode.Id {
				PredecessorNode := self.ringLookup(self.Fingers[i].key)
				responsibleNode := PredecessorNode.Successor
				self.Fingers[i].Id = responsibleNode.Id[:len(responsibleNode.Id)]

			}
		}
		self = self.Successor
	}
}

func (self *DHTnode) updateAllFingerTables() { // updates all Fingers in fingerTables of all nodes, starts with self

	start := self

	for start != self.Successor {
		for i := 0; i < 160; i++ {
			responsibleNode := self.ringLookup(self.Fingers[i].key)
			self.Fingers[i].Id = responsibleNode.Id[:len(responsibleNode.Id)]
		}
		self = self.Successor
	}

	// filling finger table for last node before starting node
	for i := 0; i < 160; i++ {
		responsibleNode := self.ringLookup(self.Fingers[i].key)
		self.Fingers[i].Id = responsibleNode.Id[:len(responsibleNode.Id)]
	}

}

func createNode(port string) *DHTnode {

	NodeIp := "localhost"
	NodePort := port

	joinViaIp := "localhost"
	joinViaPort := "1111"

	node := makeDHTNode(NodeIp, NodePort, joinViaIp, joinViaPort)

	return node
}

func makeDHTNode(NodeIp string, NodePort string, joinViaIp string, joinViaPort string) *DHTnode {

	IdStr := sha1hash(NodeIp + NodePort)
	IdByte, _ := hex.DecodeString(IdStr)

	simpleNode := Node{IdStr, IdByte, NodeIp, NodePort}
	node := &DHTnode{simpleNode, "", "", nil, nil, nil, nil, nil, joinViaIp, joinViaPort}

	FingersWanted := 160
	for i := 0; i < FingersWanted; i++ {
		fingerNumber := i + 1
		newFingerKey := calcFinger(node.IdByte, fingerNumber, 160)
		newFinger := &Finger{newFingerKey, "", nil, "", ""}
		node.Fingers = append(node.Fingers, newFinger)
	}

	newPredecessor0 := &NeighbourNode{"", nil, "", ""}
	newPredecessor1 := &NeighbourNode{"", nil, "", ""}
	newSuccessor0 := &NeighbourNode{"", nil, "", ""}
	newSuccessor1 := &NeighbourNode{"", nil, "", ""}
	node.Predecessors = append(node.Predecessors, newPredecessor0)
	node.Predecessors = append(node.Predecessors, newPredecessor1)
	node.Successors = append(node.Successors, newSuccessor0)
	node.Successors = append(node.Successors, newSuccessor1)

	return node
}

//func (self *DHTnode) addToRing(node *DHTnode) {
//
//	/*
//	   // instead of traverings all nodes from self until finding point of insertion,
//	   //Fingers of existing nodes should be used
//	*/
//
//	if self.Successor == nil { // new node connects to a single node, forming a ring of two nodes
//
//		self.Successor = node
//		node.Predecessor = self
//		self.Successors[0].Id = node.Id[:len(node.Id)]
//		self.Successors[0].NodeIp = node.NodeIp[:len(node.NodeIp)]
//		self.Successors[0].NodePort = node.NodePort[:len(node.NodePort)]
//		self.Successors[1].Id = self.Id[:len(self.Id)]
//		self.Successors[1].NodeIp = self.NodeIp[:len(self.NodeIp)]
//		self.Successors[1].NodePort = self.NodePort[:len(self.NodePort)]
//		node.Predecessors[0].Id = self.Id[:len(self.Id)]
//		node.Predecessors[0].NodeIp = self.NodeIp[:len(self.NodeIp)]
//		node.Predecessors[0].NodePort = self.NodePort[:len(self.NodePort)]
//		node.Predecessors[1].Id = node.Id[:len(node.Id)]
//		node.Predecessors[1].NodeIp = node.NodeIp[:len(node.NodeIp)]
//		node.Predecessors[1].NodePort = node.NodePort[:len(node.NodePort)]
//
//		node.Successor = self
//		self.Predecessor = node
//		node.Successors[0].Id = self.Id[:len(self.Id)]
//		node.Successors[0].NodeIp = self.NodeIp[:len(self.NodeIp)]
//		node.Successors[0].NodePort = self.NodePort[:len(self.NodePort)]
//		node.Successors[1].Id = node.Id[:len(node.Id)]
//		node.Successors[1].NodeIp = node.NodeIp[:len(node.NodeIp)]
//		node.Successors[1].NodePort = node.NodePort[:len(node.NodePort)]
//		self.Predecessors[0].Id = node.Id[:len(node.Id)]
//		self.Predecessors[0].NodeIp = node.NodeIp[:len(node.NodeIp)]
//		self.Predecessors[0].NodePort = node.NodePort[:len(node.NodePort)]
//		self.Predecessors[1].Id = self.Id[:len(self.Id)]
//		self.Predecessors[1].NodeIp = self.NodeIp[:len(self.NodeIp)]
//		self.Predecessors[1].NodePort = self.NodePort[:len(self.NodePort)]
//
//	} else {
//
//		for !between([]byte(self.Id), []byte(self.Successors[0].Id), []byte(node.Id)) {
//
//			self = self.Successor
//
//		}
//
//		if self.Successors[1].Id == self.Id { // new node connects to a ring of two nodes
//
//			node.Successor = self.Successor
//			node.Successor.Predecessor = node
//			node.Successors[0].Id = self.Successors[0].Id[:len(self.Successors[0].Id)]
//			node.Successors[0].NodeIp = self.Successors[0].NodeIp[:len(self.Successors[0].NodeIp)]
//			node.Successors[0].NodePort = self.Successors[0].NodePort[:len(self.Successors[0].NodePort)]
//			node.Successors[1].Id = self.Successors[1].Id[:len(self.Successors[1].Id)]
//			node.Successors[1].NodeIp = self.Successors[1].NodeIp[:len(self.Successors[1].NodeIp)]
//			node.Successors[1].NodePort = self.Successors[1].NodePort[:len(self.Successors[1].NodePort)]
//			node.Successor.Predecessors[0].Id = node.Id[:len(self.Id)]
//			node.Successor.Predecessors[0].NodeIp = node.NodeIp[:len(self.NodeIp)]
//			node.Successor.Predecessors[0].NodePort = node.NodePort[:len(self.NodePort)]
//			node.Successor.Predecessors[1].Id = self.Id[:len(self.Id)]
//			node.Successor.Predecessors[1].NodeIp = self.NodeIp[:len(self.NodeIp)]
//			node.Successor.Predecessors[1].NodePort = self.NodePort[:len(self.NodePort)]
//			node.Successor.Successors[1].Id = node.Id[:len(node.Id)]
//			node.Successor.Successors[1].NodeIp = node.NodeIp[:len(node.NodeIp)]
//			node.Successor.Successors[1].NodePort = node.NodePort[:len(node.NodePort)]
//
//			self.Successor = node
//			node.Predecessor = self
//			self.Successors[0].Id = node.Id[:len(node.Id)]
//			self.Successors[0].NodeIp = node.NodeIp[:len(node.NodeIp)]
//			self.Successors[0].NodePort = node.NodePort[:len(node.NodePort)]
//			self.Successors[1].Id = node.Successors[0].Id[:len(node.Successors[0].Id)]
//			self.Successors[1].NodeIp = node.Successors[0].NodeIp[:len(node.Successors[0].NodeIp)]
//			self.Successors[1].NodePort = node.Successors[0].NodePort[:len(node.Successors[0].NodePort)]
//			node.Predecessors[0].Id = self.Id[:len(self.Id)]
//			node.Predecessors[0].NodeIp = self.NodeIp[:len(self.NodeIp)]
//			node.Predecessors[0].NodePort = self.NodePort[:len(self.NodePort)]
//			node.Predecessors[1].Id = self.Predecessors[0].Id[:len(node.Predecessors[0].Id)]
//			node.Predecessors[1].NodeIp = self.Predecessors[0].NodeIp[:len(node.Predecessors[0].NodeIp)]
//			node.Predecessors[1].NodePort = self.Predecessors[0].NodePort[:len(node.Predecessors[0].NodePort)]
//			self.Predecessors[1].Id = node.Id[:len(node.Id)]
//			self.Predecessors[1].NodeIp = node.NodeIp[:len(node.NodeIp)]
//			self.Predecessors[1].NodePort = node.NodePort[:len(node.NodePort)]
//
//		} else { // new node connects to a ring of at least three nodes
//
//			node.Successor = self.Successor
//			node.Successor.Predecessor = node
//			node.Successors[0].Id = self.Successors[0].Id[:len(self.Successors[0].Id)]
//			node.Successors[0].NodeIp = self.Successors[0].NodeIp[:len(self.Successors[0].NodeIp)]
//			node.Successors[0].NodePort = self.Successors[0].NodePort[:len(self.Successors[0].NodePort)]
//			node.Successors[1].Id = self.Successors[1].Id[:len(self.Successors[1].Id)]
//			node.Successors[1].NodeIp = self.Successors[1].NodeIp[:len(self.Successors[1].NodeIp)]
//			node.Successors[1].NodePort = self.Successors[1].NodePort[:len(self.Successors[1].NodePort)]
//			node.Successor.Predecessors[0].Id = node.Id[:len(self.Id)]
//			node.Successor.Predecessors[0].NodeIp = node.NodeIp[:len(self.NodeIp)]
//			node.Successor.Predecessors[0].NodePort = node.NodePort[:len(self.NodePort)]
//			node.Successor.Predecessors[1].Id = self.Id[:len(self.Id)]
//			node.Successor.Predecessors[1].NodeIp = self.NodeIp[:len(self.NodeIp)]
//			node.Successor.Predecessors[1].NodePort = self.NodePort[:len(self.NodePort)]
//			node.Successor.Successor.Predecessors[1].Id = node.Id[:len(self.Id)]
//			node.Successor.Successor.Predecessors[1].NodeIp = node.NodeIp[:len(self.NodeIp)]
//			node.Successor.Successor.Predecessors[1].NodePort = node.NodePort[:len(self.NodePort)]
//
//			self.Successor = node
//			node.Predecessor = self
//			self.Predecessor.Successors[1].Id = node.Id[:len(node.Id)]
//			self.Predecessor.Successors[1].NodeIp = node.NodeIp[:len(node.NodeIp)]
//			self.Predecessor.Successors[1].NodePort = node.NodePort[:len(node.NodePort)]
//			self.Successors[0].Id = node.Id[:len(node.Id)]
//			self.Successors[0].NodeIp = node.NodeIp[:len(node.NodeIp)]
//			self.Successors[0].NodePort = node.NodePort[:len(node.NodePort)]
//			self.Successors[1].Id = node.Successors[0].Id[:len(node.Successors[0].Id)]
//			self.Successors[1].NodeIp = node.Successors[0].NodeIp[:len(node.Successors[0].NodeIp)]
//			self.Successors[1].NodePort = node.Successors[0].NodePort[:len(node.Successors[0].NodePort)]
//			node.Predecessors[0].Id = self.Id[:len(self.Id)]
//			node.Predecessors[0].NodeIp = self.NodeIp[:len(self.NodeIp)]
//			node.Predecessors[0].NodePort = self.NodePort[:len(self.NodePort)]
//			node.Predecessors[1].Id = self.Predecessors[0].Id[:len(node.Predecessors[0].Id)]
//			node.Predecessors[1].NodeIp = self.Predecessors[0].NodeIp[:len(node.Predecessors[0].NodeIp)]
//			node.Predecessors[1].NodePort = self.Predecessors[0].NodePort[:len(node.Predecessors[0].NodePort)]
//
//		}
//	}
//
//	//self.updateAllFingerTables()
//}

/* AddToRing
Available for rpc
@arg.FirstNode is the node which nodeAdded
*/ /*Work in progress
/*func (t *DHTnode) AddToRing(arg *ArgAddToRing, nodeAdded *DHTnode) error {

}
*/

//implem algo from chord doc p6
/*
func (t *DHTnode) init_finger_table (nodeJoined *DHTnode) {
	thisNode.Fingers[0] = nodeJoined.lookup(calcFingerSha(0)
	thisNode.Predecessor = thisNode.Successor.Predecessor
	thisNode.Successor.Predecessor = thisNode

	for i :=0; i< m-1; i++ {
		fingerStart := calcFingerSha(thisNode.Finger[i+1]
		if (between(thisNode.Id, thisNode.Finger[i], ([]byte)fingerStart)) {
			thisNode.Finger[i+1].node = thisNode.Finger[i]
		}
		else {
			//Can't work like that. lookup return DHTnode and I must receive Finger
			//need to create node struct
			thisNode.Finger[i+1] = nodeJoined.lookup(fingerStart)
		}
*/
