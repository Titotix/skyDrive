package main

import (
	"fmt"
	"strconv"
)

func printNode(node DHTnode) {

	fmt.Printf("node id: %s\n", node.Id)
	fmt.Printf("Finger   1    id: %s\n", node.Fingers[0].Id)
	fmt.Println("successor : " + node.Successor.Id)

}

func printRing(self DHTnode) {
	fmt.Println("\nNodes in ring:")
	start := self
	fmt.Println("\n")
	fmt.Println(" ****** Ring : *******")
	fmt.Println("node 1")
	self.print()

	if self.Fingers[0].Id != "" && self.Fingers[0].Id != self.Id {
		self.Node = self.Fingers[0].Node
		i := 2
		for self.Id != start.Id {

			fmt.Printf("node %d\n", i)
			i++
			self.print()
			next, _ := add(self.Id, 1)
			save := self
			self.Node = self.findSuccessor(next)
			fmt.Println("reponse finSucc :" + self.Node.Id + "next :" + next)
			if save.Id == self.Id {
				break
			}
		}
	} else {
		fmt.Println("ring only 1 node")
	}
}

//func (self *DHTnode) printRing() {
//	fmt.Println("\nNodes in ring:")
//	start := self
//	printNode(*self)
//	if self.Successor.Id != "" {
//		self.BasicNode = self.Successor
//		for self != start {
//			printNode(*self)
//			self.BasicNode = self.Successor
//		}
//	}
//}

func (self *BasicNode) print() {
	fmt.Println("Id :" + self.Id)
	self.printIdByte()
	fmt.Printf("\n")
}

func (self *Node) print() {
	fmt.Println("Id :" + self.Id)
	self.BasicNode.printIdByte()
	fmt.Printf(" * Successor :\n")
	self.Successor.print()
	fmt.Printf(" * Predecessor :\n")
	self.Predecessor.print()
}

func (self *Finger) print() {
	fmt.Println("finger key:" + self.key)
	printIdByte(self.keyByte)
	self.Node.print()
}

func (self *ComparableNode) print() {
	fmt.Println("Id :" + self.Id)
	fmt.Printf("\n")
}

func (self *DHTnode) printFingers() {
	m := 160
	for i := 0; i < m; i++ {
		fmt.Println("\nFinger " + strconv.Itoa(i+1))
		fmt.Printf("key :" + self.Fingers[i].key + "\n")
		self.Fingers[i].print()
	}
	fmt.Printf("\n")
}

func (self *BasicNode) printIdByte() {
	fmt.Printf("%x", self.IdByte)
	fmt.Printf("\n")
}

func printIdByte(idByte []byte) {
	fmt.Printf("%x", idByte)
	fmt.Printf("\n")
}
