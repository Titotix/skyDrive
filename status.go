package main

import (
	"fmt"
	"time"
)

type ArgStatus struct {}

// returns true whenever node is online
func (n *Node) NodeStatus(arg *ArgStatus, statusReply *bool) error {
    *statusReply = true
    return nil
}


// checks and prints status of node given in argument (i.e. node.successor or node.predecessor)
// should call local and/or remote methods to reconnect ring and handle data when remote node is offline
func checkStatus(n *Node) {
	var statusReply bool
	arg := &ArgStatus{}
    for {
    	fmt.Printf("node." + n.nodeId + ": ")
		err := n.NodeStatus(arg, &statusReply) 
		if (statusReply == true && err == nil) {
			fmt.Printf("ok\n")
		} else {
			fmt.Printf("not ok\n")
			lockStorage()
			if n == node.predecessor {	// predecessor of calling node unavailble
				//blockRemoteAccess("pred", "node")
				//reconnectRing(n.predecessor.predecessor)
				//storeReplicatedData("storage/pred", "storage/node")
				//replicateOwnData(n.predecessor)
				//allowRemoteAccess("pred", "node")
			} else {		 			// successor of calling node unavailble
				//blockRemoteAccess("node", "succ")
				//reconnectRing(n.successor.successor)
				//moveReplicatedData("storage/succ", n.successor)
				//replicateOwData(n.successor)
				//allowRemoteAccess("node", "succ")
			}
			
		}
		// Delay of 1 second before next status check
		time.Sleep(1000 * time.Millisecond)
    }
}
