package main

import (
	"fmt"
	"time"
)

// responds to status checks
func (n *Node) NodeStatus(arg *ArgStatus, statusReply *bool) error {
	*statusReply = true
	return nil
}

// checks status at neighbour nodes and reacts if other node is not ok
func checkStatus(n *Node, interval time.Duration) {
	//var statusReply bool
	//arg := &ArgStatus{}
	for {
		isAlive := isAlive(n.BasicNode)
		if isAlive == true {
		} else {

			//if n.nodeId = node.predecessor {	// predecessor unavailble
			//blockRemoteAccess("pred", "node")
			//reconnectRing(n.predecessor.predecessor)
			//moveData() // stores data from previous predecessor
			//replicateData("node", n.predecessor, "succ")  // replicates data to new predecessor
			//allowRemoteAccess("pred", "node")
			//	} else {		 		// successor unavailble
			//blockRemoteAccess("node", "succ")

			//TODO
			//n.getPredecessor().reconnectRing(n.Successor)
			//replicateData("succ", &n.Successor, "node") // restores lost data to new sucessor
			//replicateData("node", &n.Successor, "pred") // replicates own data to new succ
			//allowRemoteAccess("node", "succ")
			//	}
		}
		time.Sleep(interval * time.Millisecond)
	}
}
