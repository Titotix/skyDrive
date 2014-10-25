package main

import (
	"encoding/hex"
	"log"
)

// @param : node who is looking for the node responsible for key
// @param : key
//DEPRECATED
/*
func (t *DHTnode) Lookup(arg *ArgLookup, responsibleNode *DHTnode) error {

	//keyByte, _ := hex.DecodeString(arg.Key)
	fmt.Printf("node id : %s  :::: port : %s", thisNode.Id, thisNode.Port)

	if thisNode.Id == arg.Key {
		*responsibleNode = arg.Node
		return nil
		//} else if arg.Node.Successor.Id == arg.Key {
		//	responsibleNode = arg.Node.Successor
		//	return nil
		//} else if between(arg.Node.IdByte, arg.Node.Successor.IdByte, keyByte) {
		//	*responsibleNode = arg.Node
		//	return nil
		//} else if bytes.Compare(keyByte, arg.Node.Successor.IdByte) == 1 || bytes.Compare(keyByte, arg.Node.IdByte) == -1 {

		//arg.Node.Successor.Lookup(arg, responsibleNode)
		//	fmt.Printf("ECHEC")
	}
	return errors.New("Lookup failed")
}
*/

/* @param : node who is looking for the node responsible for key
@param : key
Use thisNode global variable
*/
/*
func (self *DHTnode) FingerLookup(arg *ArgLookup, responsibleNode *DHTnode) error {

	targetId := ""
	*responsibleNode = *self
	fmt.Println("debut FingerLookup on ", self.Port)

	if between([]byte(self.Predecessor.Id), []byte(self.Id), []byte(arg.Key)) || self.Id == arg.Key { // self is responsible for key
		fmt.Println("arg.key is between predeccessor.Id et self.id")
		*responsibleNode = *self
		fmt.Println("predecessor Id : " + self.Predecessor.Id)
		return nil

	} else {

		fmt.Println("else statement of FingerLookup")
		// deciding finger to use by iteration, replace with better algoritm???
		fingerFound := false
		i := 0
		for (!(fingerFound == true)) && (i < 159) {

			if between([]byte(self.Fingers[i].key), []byte(self.Fingers[i+1].key), []byte(arg.Key)) {

				targetId = self.Fingers[i].Id
				fingerFound = true

			} else {
				i = i + 1
			}
		}
		if !fingerFound {
			targetId = self.Fingers[159].Id
		}

		// traversing ring clockwise instead of send request directly via IP of node
		for !(self.Successor.Id == targetId) {
			self.Node = *self.Successor
		}

		// recursive request to closest node pointed to by finger
		*responsibleNode = *self.Successor.lookup(arg.Key)
		return nil
	}
	return nil
}
*/

func (self *DHTnode) ringLookup(hashedKey string) BasicNode {

	nodeFound := false
	key, err := hex.DecodeString(hashedKey)
	if err != nil {
		log.Fatal("decodeString :", err)
	}

	for nodeFound == false {

		id1 := (self.IdByte)
		id2 := (self.Successor.IdByte)

		if between(id1, id2, key) {
			nodeFound = true
		} else {
			self.BasicNode = self.Successor
		}
	}
	return self.Successor
}
