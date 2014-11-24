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

type ArgStorage struct {
	Key          string
	Data         string
	StorageSpace string
}

type ArgDeletion struct {
	StorageSpace string
	Key          string
}

type ArgListing struct {
	storageSpace string
}

type ArgEmpty struct{}

type ArgUpdateFingerFromDeadOne struct {
	DeadNode BasicNode
}

type ArgStoreData struct {
	Key          string
	Data         string
	StorageSpace string
}

type ArgDeleteData struct {
	StorageSpace string
	Key          string
}

type ArgGetData struct {
	Key string
}

// Current node is going to connect to remote http server (@host, @port)
func connect(host string, port string) *rpc.Client {
	client, err := rpc.DialHTTP("tcp", host+":"+port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	return client
}

//Just an abstraction of method connect
func connectToNode(nodeTarget BasicNode) *rpc.Client {
	return connect(nodeTarget.Ip, nodeTarget.Port)
}

func isAlive(node BasicNode) bool {
	client, err := rpc.DialHTTP("tcp", node.Ip+":"+node.Port)
	if err != nil {
		fmt.Println(err)
		return false
	}
	client.Close()
	return true
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
	err := clientSocket.Call("DHTnode.UpdateFingerTable", arg, &reply)
	if err != nil {
		if false == handleDeadNode(nodeTarget) {
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

func callGetPredecessor(clientSocket *rpc.Client, nodeTarget BasicNode, arg *ArgEmpty) BasicNode {
	var reply BasicNode
	err := clientSocket.Call("DHTnode.GetPredecessor", arg, &reply)
	if err != nil {
		if false == handleDeadNode(nodeTarget) {
			log.Fatal("remote getPredecessor error:", err)
		} else {
			//The node is dead, happy predecessor deadNode !
			deadPred := thisNode.findPredecessor(nodeTarget.Id)
			if deadPred.Id == nodeTarget.Id {
				log.Fatal("Shit happens")
			} else {

				reply = deadPred.BasicNode
			}
		}
	}
	return reply
}

func (nodeTarget *BasicNode) getPredecessor() BasicNode {

	arg := new(ArgEmpty)
	if nodeTarget.Id == thisNode.Id {
		// execute in local
		//fmt.Println("exec in local")
		reply := new(BasicNode)
		_ = thisNode.GetPredecessor(arg, reply)
		return *reply
	} else {
		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
		reply := callGetPredecessor(clientSocket, *nodeTarget, arg)
		clientSocket.Close()
		return reply
	}
}

func callGetSuccessor(clientSocket *rpc.Client, nodeTarget BasicNode, arg *ArgEmpty) BasicNode {
	var reply BasicNode
	err := clientSocket.Call("DHTnode.GetSuccessor", arg, &reply)
	if err != nil {
		if false == handleDeadNode(nodeTarget) {
			log.Fatal("remote getSuccessor error:", err)
		} else {
			//The node is dead, happy predecessor deadNode !
			deadSucc := thisNode.findSuccessor(nodeTarget.Id)
			if deadSucc.Id == nodeTarget.Id {
				log.Fatal("Shit happens")
			} else {

				reply = deadSucc.BasicNode
			}
		}
	}
	return reply
}

func (nodeTarget *BasicNode) getSuccessor() BasicNode {

	arg := new(ArgEmpty)
	if nodeTarget.Id == thisNode.Id {
		// execute in local
		//fmt.Println("exec in local")
		reply := new(BasicNode)
		_ = thisNode.GetSuccessor(arg, reply)
		return *reply
	} else {
		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
		reply := callGetSuccessor(clientSocket, *nodeTarget, arg)
		clientSocket.Close()
		return reply
	}
}

func callUpdateFingerFromDeadOne(clientSocket *rpc.Client, nodeTarget BasicNode, arg *ArgUpdateFingerFromDeadOne) {
	var reply Node
	err := clientSocket.Call("DHTnode.UpdateFingerFromDeadOne", arg, &reply)
	if err != nil {
		if false == handleDeadNode(nodeTarget) {
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
		reply := new(Node)
		_ = thisNode.UpdateFingerFromDeadOne(arg, reply)
	} else {
		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
		callUpdateFingerFromDeadOne(clientSocket, *nodeTarget, arg)
		clientSocket.Close()
	}
}

func callUpdateFingerTableFirstNode(clientSocket *rpc.Client, nodeTarget BasicNode, arg *ArgUpdateFingerTable) {
	var reply int
	err := clientSocket.Call("DHTnode.UpdateFingerTableFirstNode", arg, &reply)
	if err != nil {
		if false == handleDeadNode(nodeTarget) {
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
	err := clientSocket.Call("DHTnode.FindSuccessor", arg, &reply)
	if err != nil {
		if false == handleDeadNode(nodeTarget) {
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
	err := clientSocket.Call("DHTnode.FindPredecessor", arg, &reply)
	if err != nil {
		if false == handleDeadNode(nodeTarget) {
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

//func callNodeStatus(clientSocket *rpc.Client, arg *ArgStatus) bool {
//	var reply bool
//	err := clientSocket.Call("DHTnode.NodeStatus", arg, &reply)
//	if err != nil {
//		return false
//	}
//	return reply
//}

//func (nodeTarget *Node) nodeStatus() bool {
//
//	arg := new(ArgStatus)
//	if nodeTarget.Id == thisNode.Id {
//		// execute in local
//		reply := new(bool)
//		_ = nodeTarget.NodeStatus(arg, reply)
//		return *reply
//	} else {
//		clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
//		reply := callNodeStatus(clientSocket, arg)
//		clientSocket.Close()
//		return reply
//	}
//}

/*
Purpose : handle error in rpc call to dead node.
Check with isAlive() if deadode is really dead
If so, updateFingerFromDeadOne
return true if deadNode was really dead
return false if  not
*/
func handleDeadNode(deadNode BasicNode) bool {
	//Check a second time if deadNode is really dead
	reply := isAlive(deadNode)
	if reply == false {
		//BUG TODO findPredecessor is going to fall on deadNode -> By the way, pb handle by callFindPredecessor (normally
		deadNodePred := thisNode.findPredecessor(deadNode.Id)
		deadNodePred.updateFingerFromDeadOne(deadNode)
		// Add replication problematic
		return true
	} else {
		//DeadNode is not dead at all !
		return false
	}
}

//Code used in data management. Code which exist but dont compiled.
///*
//
//
//Abstract RPC for GetData method
//@arg : ArgDeleteData{ storageSpace, key}
//*/
//func (self *DHTnode) callGetData(clientSocket *rpc.Client, arg *ArgGetData) string {
//	var reply string
//	err := clientSocket.Call("DHTnode.GetData", arg, &reply)
//	if err != nil {
//		log.Fatal("remote GetData error on :", self.Ip, ":", self.Port, " ", err)
//	}
//	return reply
//}
//
////Abstract callGetData method
//// nodeTarget is the node where rpc is computed
//func (nodeTarget *DHTnode) getDataRemote(key string) string {
//
//	clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
//	arg := ArgGetData{key}
//	reply := nodeTarget.callGetData(clientSocket, &arg)
//	clientSocket.Close()
//	return reply
//}
//
///*
//Abstract RPC for DeleteData method
//@arg : ArgDeleteData{ storageSpace, key}
//*/
//func (self *Node) callDeleteData(clientSocket *rpc.Client, arg *ArgDeleteData) bool {
//	var reply bool
//	err := clientSocket.Call("DHTnode.DeleteData", arg, &reply)
//	if err != nil {
//		log.Fatal("remote DeleteData error on :", self.Ip, ":", self.Port, " ", err)
//	}
//	return reply
//}
//
////Abstract callDeleteData method
//// nodeTarget is the node where rpc is computed
//func (nodeTarget *Node) deleteDataRemote(storageSpace string, key string) bool {
//
//	clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
//	arg := ArgDeleteData{storageSpace, key}
//	reply := nodeTarget.callDeleteData(clientSocket, &arg)
//	clientSocket.Close()
//	return reply
//}
//
///*
//Abstract RPC for StoreData method
//@arg : ArgStoreData{ key, data, storageSpace}
//*/
//func (self *DHTnode) callStoreData(clientSocket *rpc.Client, arg *ArgStoreData) bool {
//	var reply bool
//	err := clientSocket.Call("DHTnode.StoreData", arg, &reply)
//	if err != nil {
//		log.Fatal("remote StoreData error on :", self.Ip, ":", self.Port, " ", err)
//	}
//	return reply
//}
//
////Abstract callStoreData method
//// nodeTarget is the node where rpc is computed
//func (nodeTarget *DHTnode) storeDataRemote(key string, data string, storageSpace string) bool {
//
//	clientSocket := connect(nodeTarget.Ip, nodeTarget.Port)
//	arg := ArgStoreData{key, data, storageSpace}
//	reply := nodeTarget.callStoreData(clientSocket, &arg)
//	clientSocket.Close()
//	return reply
//}
