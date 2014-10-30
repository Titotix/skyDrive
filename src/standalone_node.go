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

	//fmt.Println("addr : " + net.InterfaceAddrs[0].String())
	//fmt.Println("addr : " + net.InterfaceAddrs.String())
	//fmt.Println("addr : " + net.InterfaceAddrs[1].String())
	thisNode = new(DHTnode)
	//thisNode.StorageInit()

	/*
		//dataStored := false;
		//argStore := &ArgStorage{sha1hash("testkey"), "testdata", "node"}
		//err := thisNode.StoreData(argStore, &dataStored)

		// Testing to list data stored on local node
		dataListed := false;
		argList := &ArgListing{"node"}
		err = thisNode.ListStoredData(argList, &dataListed)
		if err != nil {
			log.Fatal(err)
		}
	*/

	defaultPort = "9999"
	var nodePort string
	fmt.Printf("\nNew node is starting...\n")

	fmt.Printf("IP for this node. Keep empty for localhost: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		{
			break
		}
	}
	nodeIp := scanner.Text()
	if nodeIp == "" {
		fmt.Println("Creating a localhost node")
		nodeIp = "localhost"
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

	firstNode := createFirstNode("localhost", "5555")
	*thisNode = makeDHTNode(nodeIp, nodePort)

	thisNode.join(firstNode)

	//Enable listening for rpc
	thisNode.listenHTTP(nodePort)
	go thisNode.checkFingers(1000)

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
