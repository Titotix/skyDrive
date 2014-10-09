package main

import (
	//	"encoding/hex"
	"errors"
	"fmt"
)

// @param : node who is looking for the node responsible for key
// @param : key
func (thisNode *DHTnode) Lookup(arg *ArgLookup, nodeResponsible *DHTnode) error {

	//keyByte, _ := hex.DecodeString(arg.Key)
	fmt.Printf("node id : %s  :::: port : %s", thisNode.NodeId, thisNode.NodePort)

	if arg.Node.NodeId == arg.Key {
		*nodeResponsible = arg.Node
		return nil
		//} else if arg.Node.Successor.NodeId == arg.Key {
		//	nodeResponsible = arg.Node.Successor
		//	return nil
		//} else if between(arg.Node.NodeIdByte, arg.Node.Successor.NodeIdByte, keyByte) {
		//	*nodeResponsible = arg.Node
		//	return nil
		//} else if bytes.Compare(keyByte, arg.Node.Successor.NodeIdByte) == 1 || bytes.Compare(keyByte, arg.Node.NodeIdByte) == -1 {

		//arg.Node.Successor.Lookup(arg, nodeResponsible)
		//	fmt.Printf("ECHEC")
	}
	return errors.New("Lookup failed")
}

// I am working on it, not yet ready
//func (f *DHTnode) FingerLookup(arg *ArgLookup, nodeResponsible *DHTnode) error {
//
//	fmt.Printf("checking node: %s\n", self.NodeId)
//
//	targetNodeId := ""
//	responsibleNode := self
//
//	if between([]byte(self.Predecessor.NodeId), []byte(self.NodeId), []byte(key)) || self.NodeId == key { // self is responsible for key
//		return self
//
//	} else {
//
//		// deciding finger to use by iteration, replace with better algoritm???
//		fingerFound := false
//		i := 0
//		for (!(fingerFound == true)) && (i < 159) {
//
//			if between([]byte(self.Fingers[i].key), []byte(self.Fingers[i+1].key), []byte(key)) {
//
//				targetNodeId = self.Fingers[i].NodeId
//				fingerFound = true
//
//			} else {
//				i = i + 1
//			}
//		}
//		if !fingerFound {
//			targetNodeId = self.Fingers[159].NodeId
//		}
//
//		// traversing ring clockwise instead of send request directly via IP of node
//		for !(self.Successor.NodeId == targetNodeId) {
//			self = self.Successor
//		}
//
//		// recursive request to closest node pointed to by finger
//		responsibleNode = self.Successor.fingerLookup(key)
//		return responsibleNode
//	}
//
//}

func (self *DHTnode) ringLookup(hashedKey string) *DHTnode {

	nodeFound := false
	key := []byte(hashedKey)

	for nodeFound == false {

		id1 := []byte(self.NodeId)
		id2 := []byte(self.Successors[0].NodeId)

		if between(id1, id2, key) {
			nodeFound = true
		} else {
			self = self.Successor
		}
	}
	return self.Successor
}
