package main

import (
	"fmt"
	"bytes"
	"encoding/hex"	
)
// @param : node who is looking for the node responsible for key
// @param : key
func ( self *DHTnode ) lookup(key string) *DHTnode {

	keyByte, _ := hex.DecodeString(key)

	if (self.nodeId == key) {
		return self
	} else if (self.successor.nodeId == key ) {
		return self.successor
	} else if (between(self.nodeIdByte, self.successor.nodeIdByte, keyByte ) ) {
		return self
	} else if ( bytes.Compare(keyByte, self.successor.nodeIdByte) == 1 || bytes.Compare(keyByte, self.nodeIdByte) == -1 ) {
		
		return self.successor.lookup(key)
	}
	fmt.Printf("\n ***** Fail to lookup ***** \n\n");
	return nil
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


