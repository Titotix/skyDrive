package main

import (
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
