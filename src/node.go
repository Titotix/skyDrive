package main

import (
	"encoding/hex"
	"fmt"
	"log"
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

	m := 160
	for i := 0; i < m; i++ {
		fingerNumber := i + 1
		newFingerKey, newFingerKeyByte := calcFinger(node.IdByte, fingerNumber, 160)
		newFinger := &Finger{*new(Node), newFingerKey, newFingerKeyByte}
		node.Fingers = append(node.Fingers, newFinger)
	}

	return node
}



/* AddToRing
Available for rpc
@arg.FirstNode is the node which nodeAdded
*/
//Work in progress
func (self *DHTnode) join(joinedNode BasicNode) {
	if isAlive(joinedNode) {
		self.initFingerTable(joinedNode)
		self.updateOthers()
	} else {
		//First node on the ring
		m := 160
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
	self.Predecessor = successor.Predecessor
	self.Fingers[0].Predecessor = self.BasicNode
	m := 160
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
	self.initFingerSuccessor(joinedNode)
}

//TODO
func (self *DHTnode) initFingerSuccessor(joinedNode BasicNode) {
	m := 160
	for i := 0; i < m; i++ {
		//If finger i point to self node, assign self succcessor to finger successor
		if self.Fingers[i].Id == self.Id {
			self.Fingers[i].Successor = self.Successor
		} else {
			next, _ := add(self.Fingers[i].Id, 1)
			successor := joinedNode.findSuccessor(next)

			// if finger i key is after joinedNode or equal
			if between(joinedNode.IdByte, self.IdByte, self.Fingers[i].keyByte) {

				//If findSuccessor answer node after self, assign self as self.Fingers i+1 successor
				if between(self.IdByte, joinedNode.IdByte, successor.IdByte) || successor.Id == joinedNode.Id {
					self.Fingers[i].Successor = self.BasicNode
				} else {
					self.Fingers[i].Successor = successor.BasicNode
				}
				// if finger key is before joinedNode
			} else {
				//if finger i point to joinedNode
				if self.Fingers[i].Id == joinedNode.Id {
					// And findSucc is after self
					if between(self.IdByte, joinedNode.IdByte, successor.IdByte) || successor.Id == joinedNode.Id {
						//Take self as finger i successor
						self.Fingers[i].Successor = self.BasicNode
					} else {

						self.Fingers[i].Successor = successor.BasicNode
					}

				} else {
					self.Fingers[i].Successor = successor.BasicNode
				}

				//if successor.Id != joinedNode.Id {
				//	self.Fingers[i].Successor = successor.BasicNode
				//} else if self.Fingers[i].Id != self.Id {
				//	self.Fingers[i].Successor = self.BasicNode
			}
		}
	}
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
	var i int

	m := 160
	for i = 0; i < m; i++ {
		//find last node p whose i finger might be self
		lastFinger, _ := calcLastFinger(self.IdByte, i+1)
		p := self.findPredecessor(lastFinger)

		//Too lazy to go deep in UpdateFingerTable. If execute itself, false the last fingers.
		if p.Id != self.Id /* && p.Id != self.Id*/ {
			p.updateFingerTable(self.Node, i)
		}
	}
}

// if s is i finger of self, update self.Fingers with s
//Useless reply parameter, rpc doesn't work without
func (self *DHTnode) UpdateFingerTable(arg *ArgUpdateFingerTable, reply *Node) error {

	argIdByte, err := hex.DecodeString(arg.Node.Id)
	if err != nil {
		log.Fatal("err DecodeString in UpdateFingerTable :", err)
	}
	//Update finger with arg.Node  only if finger key is between self and arg.Node
	if between(self.IdByte, argIdByte, self.Fingers[arg.I].keyByte) {
		if between(self.IdByte, self.Fingers[arg.I].IdByte, arg.Node.IdByte) {
			self.Fingers[arg.I].Node = arg.Node

			//get first node preceding n
			p := self.Predecessor
			if self.ComparableNode == self.Predecessor.ComparableNode {
				self.Predecessor = arg.Node.BasicNode
				self.Successor = arg.Node.BasicNode
				p = self.Predecessor
			} //else {

			//Stop recursive updatefingerTale if self.Predecessor is the node who is updating Others
			if p.Id != arg.Node.Id /* && p.Id != self.Id*/ {
				p.updateFingerTable(arg.Node, arg.I)
			}
			return nil
			//}
		}
	}

	return nil
}

/**
* self has to be the predecessor of the deadNode
*
 */
func (self *DHTnode) UpdateFingerFromDeadOne(arg *ArgUpdateFingerFromDeadOne, reply *Node) error {
	//Could be improve in considering predecessor finger which point to the same node
	for i := 0; i < m; i++ {
		//is self fingers contains the dead node ?
		if self.Fingers[i].Id == arg.DeadNode.Id {
			fmt.Printf("\nupdate Dead :fger%d=deadNode", i+1)
			//easy update of fingers with his successor
			self.Fingers[i].BasicNode = self.Fingers[i].Successor
			next, _ := add(self.Fingers[i].Successor.Id, 1)
			self.Fingers[i].Successor = self.findSuccessor(next).BasicNode
			//TODO update also fingers predecessor
		}
	}
	//Update all node counter clockwise
	p := self.Predecessor
	fmt.Printf("\nupdate Dead :p=\"%s\"", p.Id)
	//we have to stop when p reach is arg.Node (when p is the deadNode)
	if p.Id != arg.DeadNode.Id {
		p.updateFingerFromDeadOne(arg.DeadNode)
	} else {
		//We finished the counter clockwise of nodes
		//So self.Predecessor is right now the dead node
		deadNodePred := self.findPredecessor(arg.DeadNode.Id)
		self.Predecessor = deadNodePred.BasicNode
	}
	return nil
}

//Useless reply
func (self *DHTnode) UpdateFingerTableFirstNode(arg *ArgUpdateFingerTable, reply *Node) error {
	self.Successor = arg.Node.BasicNode
	self.Predecessor = arg.Node.BasicNode

	if between(self.IdByte, self.Fingers[arg.I].IdByte, arg.Node.IdByte) {
		self.Fingers[arg.I].Node = arg.Node

		//get first node preceding n
		//fmt.Println("self : " + self.Id)
		//fmt.Println("UpdateFingerTable: self.ip before pred " + self.Ip + ":" + self.Port)

		p := self.Predecessor
		//fmt.Println("UpdateFingerTable: p.pred.ip : " + p.Ip + ":" + p.Port)

		//Stop recursive updatefingerTale if self.Predecessor is the node who is updating Others
		if p.Id != arg.Node.Id /*&& p.Id != self.Id*/ {
			p.updateFingerTable(arg.Node, arg.I)
		}
	}

	return nil
}

func (self *DHTnode) FindPredecessor(arg *ArgLookup, reply *Node) error {
	predecessor := *self

	if predecessor.Successor.Id == "" {
		log.Fatal("self.Successor unset in FindPredecessor")
	}
	for !between2(predecessor.IdByte, predecessor.Successor.IdByte, arg.KeyByte) {
		predecessor.Node = predecessor.closestPrecedingFinger(arg.Key)
	}
	*reply = predecessor.Node
	if reply.Id == "" {
		log.Fatal("FindPredecessor error : reply is empty")
	}
	return nil
}

// FindPredecessor need SUccessor and Predecessor in FIngers struct
func (self *DHTnode) ClosestPrecedingFinger(arg *ArgLookup, reply *Node) error {
	idByte := arg.KeyByte
	//fmt.Println("arg.Key :" + arg.Key)
	printIdByte(arg.KeyByte)
	for i := m - 1; i > -1; i-- {
		if inside(self.IdByte, idByte, self.Fingers[i].IdByte) {
			*reply = self.Fingers[i].Node
			if reply.Id == "" {
				log.Fatal("ClosestPrecedingFinger error : reply is empty")
			}
			return nil
		}
	}
	*reply = self.Node
	if reply.Id == "" {
		log.Fatal("ClosestPrecedingFinger error : reply is empty")
	}
	return nil
}

func (self *DHTnode) FindSuccessor(arg *ArgLookup, reply *Node) error {
	predecessor := self.findPredecessor(arg.Key)
	reply.Predecessor = predecessor.BasicNode
	reply.BasicNode = predecessor.Successor
	if reply.Id == "" {
		log.Fatal("FindSuccessor error : reply is empty")
	}
	return nil
}

func (self *DHTnode) isMyFinger(node Finger) bool {
	m := 160
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


func (thisNode *Node) retrieveData (unhashedKey) {

	hashedKey := sha1hash(unhashedKey)
	arg := &ArgLookup{hashedKey}
	reply := nil
	err := thisNode.findSuccessor(arg, &reply)
	if err != nil {
		log.Fatal(err)
	}

	storedAtNode := reply.BasicNode

	err := storedAtNode.getDataRemote(hashedKey)
	if err != nil {
		log.Fatal(err)
	}
}

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


type ArgGetting struct {
	Key string
}

func (n *BasicNode) GetData(arg *ArgGetting, dataFound *string) error {

	key := arg.Key

	_ = os.Chdir("..")
	_ = os.Chdir("..")
	_ = os.Chdir("storage")

	storageFile, err := os.Open("nodeData.txt")
	if err != nil {
		fmt.Printf("failed to open nodeData.txt")
		log.Fatal(err)
	}
	defer storageFile.Close()

	reader := bufio.NewReader(storageFile)
	searchDone := false
	fmt.Printf("\n\nFiles stored in %s space:\n", arg.storageSpace)
	for (!searchDone) {
		storedKeyDelim, err := reader.ReadBytes(',')
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		}
		storedKey := bytes.TrimSuffix(storedKeyDelim, []byte(","))
		data, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		}
		if (len(data)) == 0 {
			*dataFound = "no data was found"
			searchDone = true
		} else {
			if storedKey == key {
				*dataFound = data
				searchDone = true
			}
		}
	}
	storageFile.Close()

	_ = os.Chdir("..")
	_ = os.Chdir("new_git")	
	_ = os.Chdir("src")	

	return nil
}




type ArgStorage struct {
	Key string
	Data string
	StorageSpace string
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
/**
* Use this function when a node has been notice as dead
* deadNode has to be the successor of self
*
**/
func (self *DHTnode) reconnectRing(deadNode DHTnode) {
	self.Successor = self.Fingers[0].Successor
	self.updateFingerFromDeadOne(deadNode.BasicNode)
}

func (self *DHTnode) GetPredecessor(arg *ArgEmpty, reply *BasicNode) error {
	*reply = self.Predecessor
	return nil
}

func (self *DHTnode) GetSuccessor(arg *ArgEmpty, reply *BasicNode) error {
	*reply = self.Successor
	return nil
}
