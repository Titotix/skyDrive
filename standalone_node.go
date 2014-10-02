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
	
	fmt.Printf("\nNew node is starting...\n")

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
	
	startNodeListener(thisNode)
	fmt.Printf("\nListening on port %s\n", thisNode.nodePort)
	fmt.Printf("\nConnect to remote node on port: ")
	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		{break}
	}
	remoteNodePort := scanner.Text()
	
	msg := createMessage("1", "2", "a", "b")
	sendMessage(msg, remoteNodePort)

	fmt.Printf("\n\n")
	fmt.Printf("thisNode port: %s\n", thisNodePort)

	
}