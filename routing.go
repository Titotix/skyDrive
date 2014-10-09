package main

import (
	"errors"
	"fmt"
)

// @param : node who is looking for the node responsible for key
// @param : key
func (t *DHTnode) Lookup(arg *ArgLookup, nodeResponsible *DHTnode) error {

	//keyByte, _ := hex.DecodeString(arg.Key)
	fmt.Printf("node id : %s  :::: port : %s", thisNode.Id, thisNode.Port)

	if thisNode.Id == arg.Key {
		*nodeResponsible = arg.Node
		return nil
		//} else if arg.Node.Successor.Id == arg.Key {
		//	nodeResponsible = arg.Node.Successor
		//	return nil
		//} else if between(arg.Node.IdByte, arg.Node.Successor.IdByte, keyByte) {
		//	*nodeResponsible = arg.Node
		//	return nil
		//} else if bytes.Compare(keyByte, arg.Node.Successor.IdByte) == 1 || bytes.Compare(keyByte, arg.Node.IdByte) == -1 {

		//arg.Node.Successor.Lookup(arg, nodeResponsible)
		//	fmt.Printf("ECHEC")
	}
	return errors.New("Lookup failed")
}

/* @param : node who is looking for the node responsible for key
@param : key
Use thisNode global variable
*/
func (f *DHTnode) FingerLookup(arg *ArgLookup, nodeResponsible *DHTnode) error {

	fmt.Printf("checking node: %s\n", thisNode.Id)

	targetId := ""
	nodeResponsible = thisNode

	if between([]byte(thisNode.Predecessor.Id), []byte(thisNode.Id), []byte(arg.Key)) || thisNode.Id == arg.Key { // self is responsible for key
		*nodeResponsible = *thisNode
		return nil

	} else {

		// deciding finger to use by iteration, replace with better algoritm???
		fingerFound := false
		i := 0
		for (!(fingerFound == true)) && (i < 159) {

			if between([]byte(thisNode.Fingers[i].key), []byte(thisNode.Fingers[i+1].key), []byte(arg.Key)) {

				targetId = thisNode.Fingers[i].Id
				fingerFound = true

			} else {
				i = i + 1
			}
		}
		if !fingerFound {
			targetId = thisNode.Fingers[159].Id
		}

		// traversing ring clockwise instead of send request directly via IP of node
		for !(thisNode.Successor.Id == targetId) {
			thisNode = thisNode.Successor
		}

		// recursive request to closest node pointed to by finger
		//		*responsibleNode = thisNode.Successor.fingerLookup(arg.Key)
		return nil
	}

}

func (self *DHTnode) ringLookup(hashedKey string) *DHTnode {

	nodeFound := false
	key := []byte(hashedKey)

	for nodeFound == false {

		id1 := []byte(self.Id)
		id2 := []byte(self.Successor.Id)

		if between(id1, id2, key) {
			nodeFound = true
		} else {
			self = self.Successor
		}
	}
	return self.Successor
}
