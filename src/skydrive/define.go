package main

type BasicNode struct {
	Id     string
	Ip     string
	Port   string
	IdByte []byte
}

type Node struct {
	BasicNode
	Successor   BasicNode
	Predecessor BasicNode
}

type DHTnode struct {
	Node
	Fingers []*Finger
}

type Finger struct {
	Node
	number  int
	key     string
	keyByte []byte
}
