package main

import (
	"encoding/hex"
	"log"
	"time"
)

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

func createLocalNode(port string) DHTnode {

	NodeIp := "localhost"
	NodePort := port

	node := makeDHTNode(NodeIp, NodePort)

	return node
}

func makeDHTNode(NodeIp string, NodePort string) DHTnode {

	IdStr := sha1hash(NodeIp + NodePort)
	IdByte, _ := hex.DecodeString(IdStr)

	basicNode := BasicNode{IdStr, NodeIp, NodePort, IdByte}
	simpleNode := Node{basicNode, *new(BasicNode), *new(BasicNode)}
	node := DHTnode{simpleNode, nil}

	m := 160
	for i := 0; i < m; i++ {
		//We define the i th finger
		fingerNumber := i + 1
		newFingerKey, newFingerKeyByte := calcFinger(node.IdByte, fingerNumber, 160)

		newFinger := &Finger{*new(Node), i, newFingerKey, newFingerKeyByte}
		node.Fingers = append(node.Fingers, newFinger)
	}

	return node
}

/* AddToRing
Available for rpc
@arg.FirstNode is the node which nodeAdded
*/
func (self *DHTnode) join(joinedNode BasicNode) {
	if isAlive(joinedNode) {
		self.initFingerTable(joinedNode)
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
//TODO IMPROVE
func (self *DHTnode) initFingerTable(joinedNode BasicNode) {

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
			//TODO involve joinedNode of the new node before to ask him anything...
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

//TODO HORRIBLE STUFF
func (self *DHTnode) initFingerSuccessor(joinedNode BasicNode) {
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
//USELESS
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
			if self.Id == self.Predecessor.Id {
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

//Update finger table of the current node with the node provided in argument
// Will basically try to put arg.Node in its finger table where it fits
func (self *DHTnode) UpdateFinger(arg *ArgUpdateFinger, reply *Node) error {

	argIdByte, err := hex.DecodeString(arg.Node.Id)
	if err != nil {
		log.Fatal("err DecodeString in UpdateFingerTable :", err)
	}

	var i int
	i = 0
	var done bool
	done = false
	for i < m || done {
		//Update finger with arg.Node  only if finger key is between self and arg.Node
		if inside(self.IdByte, argIdByte, self.Fingers[i].keyByte) {
			if between(self.IdByte, self.Fingers[i].IdByte, arg.Node.IdByte) {
				self.Fingers[i].Node = arg.Node
				done = true
			}
		}
		i++
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
			//easy update of fingers with his successor
			self.Fingers[i].BasicNode = self.Fingers[i].Successor
			next := self.Fingers[i].getSuccessor()
			self.Fingers[i].Successor = next
			self.Fingers[i].Predecessor = self.Fingers[i].getPredecessor()
		}
	}
	//Update all node counter clockwise
	p := self.Predecessor
	//we have to stop when p reach is arg.Node (when p is the deadNode)
	if p.Id != arg.DeadNode.Id && p.Id != thisNode.Id && p.getPredecessor().Id != p.getSuccessor().Id {
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
		if self.Fingers[i].Id == node.Id {
			return true
		}
	}
	return false
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

//Check if i finger of self is alive or not
func (self *DHTnode) isFingerAlive(i int) bool {
	return isAlive(self.Fingers[i].BasicNode)
}

/*
Check all fingers of self with interval between each finger
 param : time in Millisecond  between fingers check
 IMPORTANT : has to be launch has***** thread *****
*/
func (self *DHTnode) checkFingers(interval time.Duration) {
	for {
		for i := 0; i < m; i++ {
			if !self.isFingerAlive(i) {
				handleDeadNode(self.Fingers[i].BasicNode)
			}

			time.Sleep(interval * time.Millisecond)
		}
	}
}
