package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type ArgLookup struct {
	Node DHTnode //useless
	Key  string
}

type ArgAddToRing struct {
	FirstNode DHTnode
}

/*
Abstract RPC for Lookup method
@arg : ArgLookup{nodeTarget.Successor, keyTarget}
With nodeTarget.Successor is the node which are going to respond to this rpc
keyTarget is the key which we are looking for
*/
func (self *DHTnode) callLookup(clientSocket *rpc.Client, arg *ArgLookup) *DHTnode {

	var reply DHTnode
	err := clientSocket.Call("DHTnode.Lookup", arg, &reply)
	if err != nil {
		log.Fatal("remote lookup error:", err)
	}
	fmt.Printf("reply : %s", reply.Id)
	return &reply
}

//Abstract callLookup method
// nodeTarget is the node where rpc is computed
func (nodeTarget *DHTnode) lookup(keyTarget string) *DHTnode {

	clientSocket := nodeTarget.connect(nodeTarget.Ip, nodeTarget.Port)
	arg := ArgLookup{*nodeTarget, keyTarget}
	reply := nodeTarget.callLookup(clientSocket, &arg)
	clientSocket.Close()
	return reply
}
