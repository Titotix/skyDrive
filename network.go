package main

import (
	"fmt"
	"strconv"
	"net"
	//"encoding/hex"
	"encoding/json"
	//"math/big"
	//"bytes"
	//"crypto/sha1"
)

func startNodeListener(node *DHTnode) bool {
	go runListener(node)
	return true
}

/*
func startNodeSender(node *DHTnode) bool {
	//dst := node.successorIp + ":" + node.successorPort
	//text := "hej"
	//data := "hhh"
	//fmt.Printf("dstAdress: %s\n", dst)
	//fmt.Printf("content: %s\n", text)
	msg := &Msg{From:"1", To:"2", Data:"a",Mess:"b"}
	send(msg, node)
	msg = &Msg{From:"3", To:"4", Data:"c",Mess:"d"}
	send(msg, node)
	return true
}
*/

func createMessage(origin string, from string, to string, action string, data string) *Msg {

	msg := &Msg{Origin: origin, From: from, To: to, Action: action, Data: data}
	return msg
}

func sendMessage(msg *Msg, port string) {

	//udpAddr, err := net.ResolveUDPAddr("udp", dhtMsg.Dst)

	url := "localhost" + ":" + port

	fmt.Printf("\nConnecting to url (ip:port): %s\n", url)

	//udpAddr, _ := net.ResolveUDPAddr("udp", "localhost:1112")
	udpAddr, _ := net.ResolveUDPAddr("udp", url)

	//conn, err := net.DialUDP("udp", nil, udpAddr)
	conn, _ := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()

	encodedMsg, _ := json.Marshal(msg)

	writeResult, _ := conn.Write([]byte(encodedMsg))
	fmt.Printf("Bytes sent: %d\n", writeResult)
}



func runListener(node *DHTnode) {

	//fmt.Println("\nListener started")
	port, _ := strconv.Atoi(node.nodePort)
	//fmt.Printf("Listening on port %d\n", port)
	addr := net.UDPAddr{
        Port: port,
        IP: net.ParseIP("localhost"),
    }
    conn, _ := net.ListenUDP("udp", &addr)

	/*
	bindAddress := node.nodeIp + ":" + node.nodePort
	fmt.Printf ("Bindadress: %s\n", bindAddress)
	//udpAddr, err := net.ResolveUDPAddr("udp", bindAddress)
	udpAddr, _ := net.ResolveUDPAddr("udp", transport.bindAddress)
	//conn, err := net.ListenUDP("udp", udpAddr)
	*/

	//conn, _ := net.ListenUDP("udp", addr)
	
	//fmt.Println("\nconnection listener established\n\n")

	defer conn.Close()
	dec := json.NewDecoder(conn)	
	//fmt.Println("Decoding started\n")
	for i := 0; i < 2; i++ {
	//for {
		//fmt.Println("\nWaiting for message")
		msg := Msg{}
		//err := json.Unmarshal(dec, &msg)
		err := dec.Decode(&msg)
		fmt.Printf("\n\nErrors when receiving : ")
		fmt.Println(err)
		//dec.Decode(&msg)
		//fmt.Println(dec)
		
		fmt.Printf("Origin: %s\n", msg.Origin)
		fmt.Printf("From: %s\n", msg.From)
		fmt.Printf("To: %s\n", msg.To)
		fmt.Printf("Action: %s\n", msg.Action)
		fmt.Printf("Data: %s\n", msg.Data)
		// we got a message
	}
	fmt.Println("Listener finished")
}


/*
func runNode(node *DHTnode) {
	fmt.Printf("\nNode %s has started\n", node.nodeId)
	for i :=0; i < 10; i++ {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("Node is running")
	}
	
	fmt.Printf("\nNode %s has finished\n", node.nodeId)
} 
*/