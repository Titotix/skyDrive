package main

import (
	"fmt"
	"time"
)

type Node struct {
	name string
}

type ArgStatus struct {}

var remoteCounter int

func (n *Node) NodeStatus(arg *ArgStatus, statusReply *bool) error {
    
    *statusReply = true

    return nil
}

func checkStatus(n *Node, interval time.Duration) {
	var statusReply bool
	arg := &ArgStatus{}
    for {
    	fmt.Printf("node." + n.nodeId + ": ")
		err := n.NodeStatus(arg, &statusReply) 
		if (statusReply == true && err == nil) {
			fmt.Printf("ok\n")
		} else {
			fmt.Printf("not ok\n")

			//if n.nodeId = node.predecessor {	// predecessor unavailble
				
				//blockRemoteAccess("pred", "node")
				//reconnectRing(n.predecessor.predecessor)
				//storeReplicatedData("storage/pred", "storage/node")
				//replicateOwnData(n.predecessor)
				//allowRemoteAccess("pred", "node")
			//} else {		 		// successor unavailble
				//blockRemoteAccess("node", "succ")
				//reconnectRing(n.successor.successor)
				//moveReplicatedData("storage/succ", n.successor)
				//replicateOwData(n.successor)
				//allowRemoteAccess("node", "succ")
			//}
			
		}
		time.Sleep(interval * time.Millisecond)
    }
}
