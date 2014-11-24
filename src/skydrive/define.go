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

type ArgLookup struct {
	Key     string
	KeyByte []byte
}

type ArgAddToRing struct {
	FirstNode DHTnode
}

type ArgUpdateFingerTable struct {
	Node Node
	I    int
}

type ArgFirstUpdate struct {
	secondNode Node
}

type ArgUpdateFingerFromDeadOne struct {
	DeadNode BasicNode
}

type ArgEmpty struct{}
