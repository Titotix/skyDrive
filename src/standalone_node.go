package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

//GLobal variable
var thisNode *DHTnode
var m int
var defaultPort string

/*
/ When finished this file should be able to replace ring.go
/ A single node is created, the node stars to lisen for connection attempts
/ from ohter node, or makes attempt to connect other node after port for
/ other node is specified.
/ nodeIP is localhost for all nodes.
*/

func main() {

	//fmt.Println("addr : " + net.InterfaceAddrs[0].String())
	//fmt.Println("addr : " + net.InterfaceAddrs.String())
	//fmt.Println("addr : " + net.InterfaceAddrs[1].String())
	thisNode = new(DHTnode)
	m = 160
	defaultPort = "9999"
	fmt.Printf("\nNew node is starting...\n")

	fmt.Printf("Port for this node: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		{
			break
		}
	}
	nodePort := scanner.Text()

	firstNode := createFirstNode("localhost", "5555")
	*thisNode = createNode(nodePort)

	thisNode.join(firstNode)

	//Enable listening for rpc
	thisNode.listenHTTP(nodePort)

	fmt.Printf("\nThis node:\n")
	thisNode.print()
	//for {
	//fmt.Println("\nSearch for key: ")
	//for scanner.Scan() {
	//	{
	//		break
	//	}
	//}
	//testKey := scanner.Text()
	////var responsibleNode DHTnode
	////_ = thisNode.FingerLookup(&ArgLookup{*thisNode, testKey}, &responsibleNode)
	//reply := thisNode.findSuccessor(testKey)
	//fmt.Printf("\nlookup result: %s\n", reply.Id)
	//}
	//testHash := sha1hash(testKey)

	//msg := createMessage(thisNodePort, thisNodePort, remoteNodePort, "a", "b")
	//sendMessage(msg, remoteNodePort)

	//Wait in put for printFingers

	//httpServer()
	for {
		for scanner.Scan() {
			{
				break
			}
		}
		input := scanner.Text()
		if input == "y" {
			thisNode.printFingers()
			fmt.Println("** MOI :")
			thisNode.print()
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("\n\n")

}

func (joinedNode *Node) setJoinedNode(ip string, port string) {

	//	fmt.Printf("IP of first node to join : (let empty for default) ")
	//	scanner := bufio.NewScanner(os.Stdin)
	//	for scanner.Scan() {
	//		{
	//			break
	//		}
	//	}
	joinedNode.Port = port
	joinedNode.Ip = ip
}
