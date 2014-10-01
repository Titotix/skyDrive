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

func (self *DHTnode) fingerLookup(key string) *DHTnode {

	fmt.Printf("current node: %s\n", self.nodeId)

	targetNodeId := ""
	responsibleNode := self

	if ( between([]byte(self.predecessor.nodeId), []byte(self.nodeId), []byte(key)) || self.nodeId == key ) {  // self is responsible for key
		return self
	
	} else { 

		// deciding finger to use by iteration, replace with better algoritm???
		fingerFound := false
		i := 0
		for ( (!(fingerFound == true)) && (i < 159) ) {

			if between( []byte(self.fingers[i].key), []byte(self.fingers[i+1].key), []byte(key)) {
				
				targetNodeId = self.fingers[i].nodeId
				fingerFound = true
				
			} else {
				i = i + 1
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
		responsibleNode = self.successor.fingerLookup(key)
		return responsibleNode
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
    return self.successor
}


