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


