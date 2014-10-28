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

/*
Abstract RPC for Lookup method
@arg : ArgLookup{nodeTarget.Successor, keyTarget}
With nodeTarget.Successor is the node which are going to respond to this rpc
keyTarget is the key which we are looking for
*/
func (self *DHTnode) callLookup(clientSocket *rpc.Client, arg *ArgLookup) *DHTnode {
	var reply DHTnode
	err := clientSocket.Call("DHTnode.FingerLookup", arg, &reply)
	if err != nil {
		log.Fatal("remote lookup error on :", self.Ip, ":", self.Port, " ", err)
	}
	fmt.Printf("reply : %s", reply.Id)
	return &reply
}

//Abstract callLookup method
// nodeTarget is the node where rpc is computed
func (nodeTarget *DHTnode) lookup(keyTarget string, keyByte []byte) *DHTnode {

	clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
	arg := ArgLookup{keyTarget, keyByte}
	reply := nodeTarget.callLookup(clientSocket, &arg)
	clientSocket.Close()
	return reply
}

func callUpdateFingerTable(clientSocket *rpc.Client, arg *ArgUpdateFingerTable) {
	var reply Node
	err := clientSocket.Call("DHTnode.UpdateFingerTable", arg, &reply)
	if err != nil {
		log.Fatal("remote updateFingerTable error:", err)
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
		callUpdateFingerTable(clientSocket, arg)
		clientSocket.Close()
	}
}

func callUpdateFingerTableFirstNode(clientSocket *rpc.Client, arg *ArgUpdateFingerTable) {
	var reply int
	err := clientSocket.Call("DHTnode.UpdateFingerTableFirstNode", arg, &reply)
	if err != nil {
		log.Fatal("remote updateFingerTableFirstNode error:", err)
	}
}

func (nodeTarget *BasicNode) updateFingerTableFirstNode(s Node, i int) {

	arg := new(ArgUpdateFingerTable)
	arg.Node = s
	arg.I = i
	if nodeTarget.Id == thisNode.Id {
		// execute in local
		fmt.Println("exec in local")
		reply := new(Node)
		_ = thisNode.UpdateFingerTableFirstNode(arg, reply)
	} else {
		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
		callUpdateFingerTable(clientSocket, arg)
		clientSocket.Close()
	}
}

func callFindSuccessor(clientSocket *rpc.Client, arg *ArgLookup) *Node {
	var reply Node
	err := clientSocket.Call("DHTnode.FindSuccessor", arg, &reply)
	if err != nil {
		log.Fatal("remote FindSuccessor error:", err)
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
		reply := callFindSuccessor(clientSocket, arg)
		clientSocket.Close()
		return *reply
	}
}

func callFindPredecessor(clientSocket *rpc.Client, arg *ArgLookup) *Node {
	var reply Node
	err := clientSocket.Call("DHTnode.FindPredecessor", arg, &reply)
	if err != nil {
		log.Fatal("remote FindPredecessor error:", err)
	}
	return &reply
}

func (nodeTarget *DHTnode) findPredecessor(key string) Node {

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
		_ = nodeTarget.FindPredecessor(arg, reply)
		return *reply
	} else {
		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
		reply := callFindPredecessor(clientSocket, arg)
		clientSocket.Close()
		return *reply
	}
}

func callClosestPrecedingFinger(clientSocket *rpc.Client, arg *ArgLookup) *Node {
	var reply Node
	err := clientSocket.Call("DHTnode.ClosestPrecedingFinger", arg, &reply)
	if err != nil {
		log.Fatal("remote closestPrecedingFinger error:", err)
	}
	return &reply
}

func (nodeTarget *DHTnode) closestPrecedingFinger(key string) Node {

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
		_ = nodeTarget.ClosestPrecedingFinger(arg, reply)
		return *reply
	} else {
		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
		reply := callClosestPrecedingFinger(clientSocket, arg)
		clientSocket.Close()
		return *reply
	}
}

func callNodeStatus(clientSocket *rpc.Client, arg *ArgStatus) *bool {
	var reply bool
	err := clientSocket.Call("DHTnode.NodeStatus", arg, &reply)
	if err != nil {
		log.Fatal("remote NodeStatus error:", err)
	}
	return &reply
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
		return *reply
	}
}
