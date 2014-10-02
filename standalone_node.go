package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

/*
/ When finished this file should be able to replace ring.go
/ A single node is created, the node stars to lisen for connection attempts
/ from ohter node, or makes attempt to connect other node after port for 
/ other node is specified.
/ nodeIP is localhost for all nodes.
*/

func main() {
	
	fmt.Printf("New node is starting...\n")

	fmt.Printf("Port for this node: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		{break}
	}
	thisNodePort := scanner.Text()
	
	port,_ := strconv.Atoi(thisNodePort)
	
	thisNode := createNode(port)

	//thisNode.updateAllFingerTables()

	thisNode.printRing()
	
	// start Listener

	fmt.Printf("Connect to node on port: ")
	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		{break}
	}
	remoteNodePort := scanner.Text()

	fmt.Printf("\n\n")
	fmt.Printf("Remote port:%s\n", remoteNodePort)
	fmt.Printf("This node port: %s\n", thisNodePort)

	// sendMessage(join, localhost, remoteNodePort)


	
}