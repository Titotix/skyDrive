package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

type ComparableNode struct {
	Id   string
	Ip   string
	Port string
}

type BasicNode struct {
	ComparableNode
	IdByte []byte
}

//I can delete pointer on succ et pred
type Node struct {
	BasicNode
	Successor   BasicNode
	Predecessor BasicNode
}

type DHTnode struct {
	Node
	storagePath    string
	storageContent string
	Fingers        []*Finger
	//	Predecessors   []*Node // larger index number = further away in ring
	//	Successors     []*Node // same as above
	joinViaIp   string
	joinViaPort string
}

/*
//DEPRECATED
func (self *DHTnode) updateIncorrectFingers() {

	start := self.Node
	newNode := self

	for start != *self.Successor {
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
*/

/*
DEPRECATED
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
*/

func createFirstNode(host string, port string) BasicNode {
	var firstNode BasicNode
	firstNode.Port = port
	firstNode.Ip = host

	id := sha1hash(host + port)
	idByte, _ := hex.DecodeString(id)

	firstNode.Id = id
	firstNode.IdByte = idByte
	return firstNode
}

func createNode(port string) DHTnode {

	//TODO get ip  addr
	NodeIp := "localhost"
	NodePort := port

	joinViaIp := "localhost"
	joinViaPort := "1111"

	node := makeDHTNode(NodeIp, NodePort, joinViaIp, joinViaPort)

	return node
}

func makeDHTNode(NodeIp string, NodePort string, joinViaIp string, joinViaPort string) DHTnode {

	IdStr := sha1hash(NodeIp + NodePort)
	IdByte, _ := hex.DecodeString(IdStr)

	basicNode := BasicNode{ComparableNode{IdStr, NodeIp, NodePort}, IdByte}
	simpleNode := Node{basicNode, *new(BasicNode), *new(BasicNode)}
	node := DHTnode{simpleNode, "", "", nil, joinViaIp, joinViaPort}

	for i := 0; i < m; i++ {
		fingerNumber := i + 1
		newFingerKey, newFingerKeyByte := calcFinger(node.IdByte, fingerNumber, 160)
		fmt.Println("           " + newFingerKey)
		printIdByte(newFingerKeyByte)
		newFinger := &Finger{*new(Node), newFingerKey, newFingerKeyByte}
		node.Fingers = append(node.Fingers, newFinger)
	}

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
//		}
//	}
//	//self.updateAllFingerTables()
//}

/* AddToRing
Available for rpc
@arg.FirstNode is the node which nodeAdded
*/
//Work in progress
func (self *DHTnode) join(joinedNode BasicNode) {
	if isAlive(joinedNode) {
		self.initFingerTable(joinedNode)
		fmt.Println("***** FINI initFingerTable *********** \n \n ******")
		self.printFingers()
		self.updateOthers()
	} else {
		//First node on the ring
		for i := 0; i < m; i++ {
			self.Fingers[i].Node = self.Node
			self.Fingers[i].Predecessor = self.BasicNode
			self.Fingers[i].Successor = self.BasicNode
		}
		self.Successor = self.BasicNode
		self.Predecessor = self.BasicNode
	}
}

//func (self *DHTnode) AddToRing(arg *ArgAddToRing, reply *Node) error {

//Must add case of second node in the ring ... OR NOT
/*
 Initalize Finger[] table for current node
*/
func (self *DHTnode) initFingerTable(joinedNode BasicNode) {
	//self.Fingers[0].key, _ = calcFingerSha(self.IdByte, 0)

	//self.basicInit(joinedNode)
	successor := joinedNode.findSuccessor(self.Fingers[0].key)
	self.Fingers[0].Node = successor
	self.Successor = successor.BasicNode

	//findSuccessor give back Predecessor ?
	//fmt.Println("successor.Pred :" + successor.Predecessor.Id)
	self.Predecessor = successor.Predecessor
	self.Fingers[0].Predecessor = self.BasicNode

	self.Successor = successor.BasicNode
	fmt.Println("initFinger : self :")
	self.print()
	thisNode.print()
	for i := 0; i < m-1; i++ {

		//If finger i+1 key is between self and node pointed by fingers i
		if between(self.IdByte, self.Fingers[i].IdByte, self.Fingers[i+1].keyByte) {
			//Mean that finger i+1 must point to the same node as finger i
			self.Fingers[i+1].Node = self.Fingers[i].Node
		} else {
			//otherwise, ask to joinedNOde the successor of finger key i+1
			self.Fingers[i+1].Node = joinedNode.findSuccessor(self.Fingers[i+1].key)
			//But if the answer is between self and joinedNode AND finger 's key asked is between joinedNode and self
			//in this case, joinedNode doesn't know yet self so he is going to  answer wrong
			if (between(self.IdByte, joinedNode.IdByte, self.Fingers[i+1].IdByte) || self.Fingers[i+1].Id == joinedNode.Id) && between(joinedNode.IdByte, self.IdByte, (self.Fingers[i+1].keyByte)) {
				//if between(joinedNode.IdByte, self.IdByte, (self.Fingers[i+1].keyByte)) {
				self.Fingers[i+1].Node = self.Node
			}
		}
	}
	fmt.Println("fin initfingertablethisNODE:" + thisNode.Id)
}

//Iniatize successor, predeccessor and finger 1 of new node in ring
func (self *DHTnode) basicInit(joinedNode BasicNode) {

	//Let's look for responsible node for the first finger
	successor := joinedNode.findSuccessor(self.Fingers[0].key)
	self.Fingers[0].Node = successor
	self.Successor = successor.BasicNode

	//findSuccessor give back Predecessor ?
	//fmt.Println("successor.Pred :" + successor.Predecessor.Id)
	self.Predecessor = successor.Predecessor
	self.Fingers[0].Predecessor = self.BasicNode
	self.Successor = successor.BasicNode
}

func (self *DHTnode) updateOthers() {
	for i := 0; i < m; i++ {
		//find last node p whose i finger might be self
		lastFinger, _ := calcLastFinger(self.IdByte, i+1)
		p := self.findPredecessor(lastFinger)
		//	if self.Id == self.Successor.Id {
		//		//current node is the second node in the ring
		//		self.updateFingerTableFirstNode(self.Node, i)
		//		fmt.Println("p : \n")
		//		p.print()
		//	} else {

		//fmt.Println("p : \n")
		//p.print()
		//fmt.Println("SELF : *****")
		//self.print()
		p.updateFingerTable(thisNode.Node, i)
		//}
	}
}

//Useless reply
func (self *DHTnode) UpdateFingerTableFirstNode(arg *ArgUpdateFingerTable, reply *Node) error {
	self.Successor = arg.Node.BasicNode
	self.Predecessor = arg.Node.BasicNode

	if between(self.IdByte, self.Fingers[arg.I].IdByte, arg.Node.IdByte) {
		self.Fingers[arg.I].Node = arg.Node

		//get first node preceding n
		fmt.Println("self : " + self.Id)
		fmt.Println("UpdateFingerTable: self.ip before pred " + self.Ip + ":" + self.Port)

		p := self.Predecessor
		fmt.Println("UpdateFingerTable: p.pred.ip : " + p.Ip + ":" + p.Port)
		p.updateFingerTable(arg.Node, arg.I)
	}

	return nil
}

// if s is i finger of self, update self.Fingers with s
//Useless reply parameter, rpc doesn't work without
func (self *DHTnode) UpdateFingerTable(arg *ArgUpdateFingerTable, reply *Node) error {

	fmt.Println("\n***** Begin UpdateFingerTable " + strconv.Itoa(arg.I) + " Node.Id " + arg.Node.Id)
	fmt.Printf("arg.Node.IdBYte :%x %x %x", arg.Node.IdByte, self.IdByte, self.Fingers[arg.I].IdByte)
	if between(self.IdByte, self.Fingers[arg.I].IdByte, arg.Node.IdByte) {
		fmt.Printf("arg.Node.IdBYte :%x ", arg.Node.IdByte)
		self.Fingers[arg.I].Node = arg.Node

		fmt.Println("self : ")
		self.print()
		fmt.Println("UpdateFingerTable: self.ip before pred " + self.Ip + ":" + self.Port)
		//get first node preceding n
		// BUG TODO : self.Predecessor == self,  so infinite loop in the case of the join of the second node
		p := self.Predecessor
		fmt.Println("p : ")
		p.print()
		if self.ComparableNode == self.Predecessor.ComparableNode {
			fmt.Println("EGAL !!!!!!!!!!!!!!!!!!!!!!!!!")
			self.Predecessor = arg.Node.BasicNode
			self.Successor = arg.Node.BasicNode
			self.Predecessor.print()
			p = self.Predecessor
		}
		fmt.Println("UpdateFingerTable: p.pred.ip : " + p.Ip + ":" + p.Port)
		p.print()
		p.updateFingerTable(arg.Node, arg.I)
	}

	return nil
}

func (self *DHTnode) FindPredecessor(arg *ArgLookup, reply *Node) error {
	predecessor := *self

	if predecessor.Successor.Id == "" {

		fmt.Println("***********µFAIL FindPredecessor cant work properly self.Successor unset")
		predecessor.print()
	}
	//Hack to use Predecessor instead of Successor
	predecessor.print()
	//doesnt not respet doc algo TODO BUG
	//Fixed, to test...
	for !between(predecessor.IdByte, predecessor.Successor.IdByte, arg.KeyByte) && arg.Key != predecessor.Successor.Id {
		predecessor.Node = predecessor.closestPrecedingFinger(arg.Key)
	}
	*reply = predecessor.Node
	return nil
}

// FindPredecessor need SUccessor and Predecessor in FIngers struct
func (self *DHTnode) ClosestPrecedingFinger(arg *ArgLookup, reply *Node) error {
	idByte := arg.KeyByte
	fmt.Println("arg.Key :" + arg.Key)
	printIdByte(arg.KeyByte)
	for i := m - 1; i > -1; i-- {
		if inside(self.IdByte, idByte, self.Fingers[i].IdByte) {
			*reply = self.Fingers[i].Node
			return nil
		}
	}
	*reply = self.Node
	return nil
}

func (self *DHTnode) FindSuccessor(arg *ArgLookup, reply *Node) error {
	predecessor := self.findPredecessor(arg.Key)
	fmt.Println("\n**asked :" + arg.Key)
	new := predecessor.BasicNode
	new.print()
	reply.Predecessor = predecessor.BasicNode
	//fmt.Println("reply :" + reply.Id + "reply pred " + reply.Predecessor.Id)
	reply.BasicNode = predecessor.Successor
	fmt.Println("\nreply :")
	reply.print()
	return nil
}

func (self *DHTnode) isMyFinger(node Finger) bool {
	for i := 0; i < m; i++ {
		if self.Fingers[i].ComparableNode == node.ComparableNode {
			return true
		}
	}
	return false
}

//func (self *DHTnode) getPredecessor() BasicNode {
//	for i:= m-1; i > -1; i-- {
//		if self.Fingers[i].Id != self.Id {
//			if self.Fingers[i].Node.findSuccessor(self.Fingers[i].Id)
//
