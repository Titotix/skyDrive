package main

import (
	"bufio"
	"fmt"
	"os"
)

//GLobal variable
var thisNode *DHTnode
var m int

/*
/ When finished this file should be able to replace ring.go
/ A single node is created, the node stars to lisen for connection attempts
/ from ohter node, or makes attempt to connect other node after port for
/ other node is specified.
/ nodeIP is localhost for all nodes.
*/

func main() {

	thisNode = new(DHTnode)
	m = 160
	fmt.Printf("\nNew node is starting...\n")

	fmt.Printf("Port for this node: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		{
			break
		}
	}
	nodePort := scanner.Text()

	thisNode = createNode(nodePort)
	//	thisNode.updateAllFingerTables()

	//Enable listening for rpc
	thisNode.listenHTTP(nodePort)
	fmt.Printf("listenHTTP done !")
	thisNode.printRing()

	//thisNode.startNodeListener()

	fmt.Printf("\nSearch for key: ")
	for scanner.Scan() {
		{
			break
		}
	}
	testKey := scanner.Text()
	//testHash := sha1hash(testKey)

	//var err error
	fmt.Printf("\nOn port node : ")
	for scanner.Scan() {
		{
			break
		}
	}
	//port := scanner.Text()

	// Care ! lookup have to be used with nodeTarget not whith thisNode (current node)
	// But as far as successor are not set, can't test for the moment with successor (by examples)
	reply := thisNode.lookup(testKey)
	fmt.Printf("\nlookup result: %s\n", reply.Id)

	//fmt.Printf("\nListening on port %s\n", thisNode.NodePort)
	//fmt.Printf("\nConnect to remote node on port: ")
	//scanner = bufio.NewScanner(os.Stdin)
	//for scanner.Scan() {
	//	{
	//		break
	//	}
	//}
	//remoteNodePort := scanner.Text()

	//msg := createMessage(thisNodePort, thisNodePort, remoteNodePort, "a", "b")
	//sendMessage(msg, remoteNodePort)

	fmt.Printf("\n\n")

}
