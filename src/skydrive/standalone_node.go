package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	//"log"
)

//GLobal variable
var thisNode *DHTnode
var m int = 160
var defaultPort string

/*
/ When finished this file should be able to replace ring.go
/ A single node is created, the node stars to lisen for connection attempts
/ from ohter node, or makes attempt to connect other node after port for
/ other node is specified.
/ nodeIP is localhost for all nodes.
*/

func main() {

	thisNode = new(DHTnode)

	defaultPort = "9999"
	var nodePort string
	var firstNodeIp string = "172.30.0.154"
	fmt.Printf("\nNew node is starting...\n")

	fmt.Printf("\nIs first Node ? \"yes\" if so.\n")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		{
			break
		}
	}
	first := scanner.Text()
	var isFirst bool
	if first == "yes" {
		isFirst = true
		nodePort = "5555"
	} else {
		isFirst = false
	}

	if isFirst == false {
		fmt.Printf("IP for this node. Keep empty for localhost: ")
		for scanner.Scan() {
			{
				break
			}
		}
		thisNode.Ip = scanner.Text()
		if thisNode.Ip == "" {
			fmt.Println("Creating a localhost node")
			thisNode.Ip = "localhost"
			firstNodeIp = "localhost"
			fmt.Printf("Port for this node in localhost : ")
			for scanner.Scan() {
				{
					break
				}
			}
			nodePort = scanner.Text()
		} else {
			nodePort = defaultPort
		}
	} else {
		fmt.Println("Do you want to create a localhost first node ? \"yes\" if so.")
		for scanner.Scan() {
			{
				break
			}
		}
		res := scanner.Text()
		if res == "yes" {
			firstNodeIp = "localhost"
			thisNode.Ip = "localhost"
		} else {
			thisNode.Ip = firstNodeIp

		}
		thisNode.Ip = firstNodeIp
	}

	firstNode := createFirstNode(firstNodeIp, "5555")
	*thisNode = makeDHTNode(thisNode.Ip, nodePort)

	thisNode.join(firstNode)

	//Enable listening for rpc
	thisNode.listenHTTP(nodePort)
	//go thisNode.checkFingers(1000)

	fmt.Printf("\nThis node:\n")
	thisNode.print()

	//httpServer()
	for {
		fmt.Printf("\n *** What do you want to do ? ***\n\n")
		fmt.Println("1) Print fingers table of current node ?")
		fmt.Println("2) Look for a responsible node of a key ?")
		for scanner.Scan() {
			{
				break
			}
		}
		input := scanner.Text()
		switch input {
		case "1":
			thisNode.printFingers()
			fmt.Println("** MOI :")
			thisNode.print()
		case "2":
			fmt.Println("Which key do you look for ?")
			for scanner.Scan() {
				{
					break
				}
			}
			key := scanner.Text()
			result := thisNode.findSuccessor(key)
			fmt.Printf("The responsible node for \"%s\" is \"%s\"\n", key, result.Id)

		default:
			fmt.Println("Unsupported input")
		}
		if input == "y" {
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
