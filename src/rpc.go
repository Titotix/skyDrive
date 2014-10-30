package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/rpc"
)

type ArgLookup struct {
	Key     string
	KeyByte []byte
}

type ArgAddToRing struct {
	FirstNode DHTnode
}

type ArgUpdateFingerTable struct {
	Node Node
	I    int
}

type ArgFirstUpdate struct {
	secondNode Node
}

type ArgStatus struct{}

type ArgUpdateFingerFromDeadOne struct {
	DeadNode BasicNode
}

/*
Abstract RPC for Lookup method
@arg : ArgLookup{nodeTarget.Successor, keyTarget}
With nodeTarget.Successor is the node which are going to respond to this rpc
keyTarget is the key which we are looking for
*/
//func (self *DHTnode) callLookup(clientSocket *rpc.Client, arg *ArgLookup) *DHTnode {
//	var reply DHTnode
//	err := clientSocket.Call("DHTnode.FingerLookup", arg, &reply)
//	if err != nil {
//		log.Fatal("remote lookup error on :", self.Ip, ":", self.Port, " ", err)
//	}
//	fmt.Printf("reply : %s", reply.Id)
//	return &reply
//}
//
////Abstract callLookup method
//// nodeTarget is the node where rpc is computed
//func (nodeTarget *DHTnode) lookup(keyTarget string, keyByte []byte) *DHTnode {
//
//	clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
//	arg := ArgLookup{keyTarget, keyByte}
//	reply := nodeTarget.callLookup(clientSocket, &arg)
//	clientSocket.Close()
//	return reply
//}

func callUpdateFingerTable(clientSocket *rpc.Client, nodeTarget BasicNode, arg *ArgUpdateFingerTable) {
	var reply Node
	fmt.Printf("\nupdateFingerTable RPC:\"%s\"\n", nodeTarget.Id)
	err := clientSocket.Call("DHTnode.UpdateFingerTable", arg, &reply)
	if err != nil {
		if false == handleDeadNode(clientSocket, nodeTarget, err) {
			log.Fatal("remote updateFingerTable error:", err)
		} else {
			//The node is dead, happy predecessor deadNode !
			deadPred := thisNode.findPredecessor(nodeTarget.Id)
			if deadPred.Id == nodeTarget.Id {
				log.Fatal("Shit happens")
			} else {
				deadPred.updateFingerTable(arg.Node, arg.I)
			}
		}
	}
}

func (nodeTarget *BasicNode) updateFingerTable(s Node, i int) {

	arg := new(ArgUpdateFingerTable)
	arg.Node = s
	arg.I = i
	if nodeTarget.Id == thisNode.Id {
		// execute in local
		//fmt.Println("exec in local")
		//reply := new(Node)
		//_ = thisNode.UpdateFingerTable(arg, reply)
		log.Fatal("updateFingerTable himself")
	} else {
		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
		callUpdateFingerTable(clientSocket, *nodeTarget, arg)
		clientSocket.Close()
	}
}

func callUpdateFingerFromDeadOne(clientSocket *rpc.Client, nodeTarget BasicNode, arg *ArgUpdateFingerFromDeadOne) {
	var reply Node
	fmt.Printf("\nupdateFingerFromDeadOne RPC:\"%s\"\n", nodeTarget.Id)
	err := clientSocket.Call("DHTnode.UpdateFingerFromDeadOne", arg, &reply)
	if err != nil {
		if false == handleDeadNode(clientSocket, nodeTarget, err) {
			log.Fatal("remote updateFingerTable error:", err)
		} else {
			deadPred := thisNode.findPredecessor(nodeTarget.Id)
			if deadPred.Id == nodeTarget.Id {
				log.Fatal("Shit happens")
			} else {
				deadPred.updateFingerFromDeadOne(arg.DeadNode)
			}
		}
	}
}

func (nodeTarget *BasicNode) updateFingerFromDeadOne(dead BasicNode) {

	arg := new(ArgUpdateFingerFromDeadOne)
	arg.DeadNode = dead
	if nodeTarget.Id == thisNode.Id {
		// execute in local
		//fmt.Println("exec in local")
		//reply := new(Node)
		//_ = thisNode.UpdateFingerFromDeadOne(arg, reply)
		log.Fatal("updateFingerFromDeadOne himself")
	} else {
		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
		callUpdateFingerFromDeadOne(clientSocket, *nodeTarget, arg)
		clientSocket.Close()
	}
}

func callUpdateFingerTableFirstNode(clientSocket *rpc.Client, nodeTarget BasicNode, arg *ArgUpdateFingerTable) {
	var reply int
	fmt.Printf("\nupdateFingerFirstNode RPC:\"%s\"\n", nodeTarget.Id)
	err := clientSocket.Call("DHTnode.UpdateFingerTableFirstNode", arg, &reply)
	if err != nil {
		if false == handleDeadNode(clientSocket, nodeTarget, err) {
			log.Fatal("remote updateFingerTable error:", err)
		} else {
			deadPred := thisNode.findPredecessor(nodeTarget.Id)
			if deadPred.Id == nodeTarget.Id {
				log.Fatal("Shit happens")
			} else {
				deadPred.updateFingerTableFirstNode(arg.Node, arg.I)
			}
		}
	}
}

func (nodeTarget *BasicNode) updateFingerTableFirstNode(s Node, i int) {

	arg := new(ArgUpdateFingerTable)
	arg.Node = s
	arg.I = i
	if nodeTarget.Id == thisNode.Id {
		// execute in local
		reply := new(Node)
		_ = thisNode.UpdateFingerTableFirstNode(arg, reply)
	} else {
		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
		callUpdateFingerTable(clientSocket, *nodeTarget, arg)
		clientSocket.Close()
	}
}

func callFindSuccessor(clientSocket *rpc.Client, nodeTarget BasicNode, arg *ArgLookup) *Node {
	var reply Node
	fmt.Printf("\nFindSuccessor RPC:\"%s\"", nodeTarget.Id)
	err := clientSocket.Call("DHTnode.FindSuccessor", arg, &reply)
	if err != nil {
		if false == handleDeadNode(clientSocket, nodeTarget, err) {
			log.Fatal("remote updateFingerTable error:", err)
		} else {
			deadPred := thisNode.findPredecessor(nodeTarget.Id)
			if deadPred.Id == nodeTarget.Id {
				log.Fatal("Shit happens")
			} else {
				deadPred.findSuccessor(arg.Key)
			}
		}
	}
	return &reply
}

func (nodeTarget *BasicNode) findSuccessor(key string) Node {

	arg := new(ArgLookup)
	arg.Key = key
	keyByte, err := hex.DecodeString(key)
	if err != nil {
		log.Fatal("findSuccessor conversion error:", err)
	}
	arg.KeyByte = keyByte
	if nodeTarget.Id == thisNode.Id {
		// execute in local
		reply := new(Node)
		_ = thisNode.FindSuccessor(arg, reply)
		return *reply
	} else {
		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
		reply := callFindSuccessor(clientSocket, *nodeTarget, arg)
		clientSocket.Close()
		return *reply
	}
}

func callFindPredecessor(clientSocket *rpc.Client, nodeTarget BasicNode, arg *ArgLookup) *Node {
	var reply Node
	fmt.Printf("\nFindPred RPC:\"%s\"", nodeTarget.Id)
	err := clientSocket.Call("DHTnode.FindPredecessor", arg, &reply)
	if err != nil {
		if false == handleDeadNode(clientSocket, nodeTarget, err) {
			log.Fatal("remote updateFingerTable error:", err)
		} else {
			deadPred := thisNode.findPredecessor(nodeTarget.Id)
			if deadPred.Id == nodeTarget.Id {
				log.Fatal("Shit happens")
			} else {
				deadPred.findPredecessor(arg.Key)
			}
		}
	}
	return &reply
}

func (nodeTarget *BasicNode) findPredecessor(key string) Node {

	arg := new(ArgLookup)
	arg.Key = key
	keyByte, err := hex.DecodeString(key)
	if err != nil {
		log.Fatal("findPredecessor ERROR DecodeString", err)
	}
	arg.KeyByte = keyByte
	if nodeTarget.Id == thisNode.Id {
		// execute in local
		reply := new(Node)
		_ = thisNode.FindPredecessor(arg, reply)
		return *reply
	} else {
		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
		reply := callFindPredecessor(clientSocket, *nodeTarget, arg)
		clientSocket.Close()
		return *reply
	}
}

func callClosestPrecedingFinger(clientSocket *rpc.Client, nodeTarget BasicNode, arg *ArgLookup) *Node {
	var reply Node
	fmt.Printf("\nclosestPreceding RPC:\"%s\"", nodeTarget.Id)
	err := clientSocket.Call("DHTnode.ClosestPrecedingFinger", arg, &reply)
	if err != nil {
		log.Fatal("remote closestPrecedingFinger error:", err)
	}
	return &reply
}

func (nodeTarget *BasicNode) closestPrecedingFinger(key string) Node {

	arg := new(ArgLookup)
	arg.Key = key
	keyByte, err := hex.DecodeString(key)
	if err != nil {
		log.Fatal("closestPrecedingFinger ERROR DecodeString", err)
	}
	arg.KeyByte = keyByte
	if nodeTarget.Id == thisNode.Id {
		// execute in local
		reply := new(Node)
		_ = thisNode.ClosestPrecedingFinger(arg, reply)
		return *reply
	} else {
		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
		reply := callClosestPrecedingFinger(clientSocket, *nodeTarget, arg)
		clientSocket.Close()
		return *reply
	}
}

func callNodeStatus(clientSocket *rpc.Client, arg *ArgStatus) bool {
	var reply bool
	fmt.Printf("\nNodeStatus RPC")
	err := clientSocket.Call("DHTnode.NodeStatus", arg, &reply)
	if err != nil {
		return false
	}
	return reply
}

func (nodeTarget *Node) nodeStatus() bool {

	arg := new(ArgStatus)
	if nodeTarget.Id == thisNode.Id {
		// execute in local
		reply := new(bool)
		_ = nodeTarget.NodeStatus(arg, reply)
		return *reply
	} else {
		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
		reply := callNodeStatus(clientSocket, arg)
		clientSocket.Close()
		return reply
	}
}

/*
Purpose : handle error in rpc call to dead node.
Check with nodeStatus if deadode is really dead
If so, updateFingerFromDeadOne
return true if deadNode was really dead
return false if  not
*/
func handleDeadNode(clientSocket *rpc.Client, deadNode BasicNode, err error) bool {
	arg := new(ArgStatus)
	//Check a second time if deadNode is really dead
	reply := callNodeStatus(clientSocket, arg)
	if reply == false {
		deadNodePred := thisNode.findPredecessor(deadNode.Id)
		deadNodePred.updateFingerFromDeadOne(deadNode)
		return true
	} else {
		//DeadNode is not dead at all !
		return false
	}
}
