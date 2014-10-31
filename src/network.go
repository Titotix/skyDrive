package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

//TODO delete parameter port
func (self *DHTnode) listenHTTP(port string) {
	rpc.Register(self)
	rpc.HandleHTTP()
	socket, e := net.Listen("tcp", ":"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(socket, nil)
}

// Current node is going to connect to remote http server (@host, @port)
func connect(host string, port string) *rpc.Client {
	client, err := rpc.DialHTTP("tcp", host+":"+port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	return client
}

//Just an abstraction of method connect
func connectToNode(nodeTarget BasicNode) *rpc.Client {
	return connect(nodeTarget.Ip, nodeTarget.Port)
}

func isAlive(node BasicNode) bool {
	client, err := rpc.DialHTTP("tcp", node.Ip+":"+node.Port)
	if err != nil {
		fmt.Println(err)
		return false
	}
	client.Close()
	return true
}
