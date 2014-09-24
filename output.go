
package main

import (
  "fmt"
)

func printNode(node *DHTnode) {

  fmt.Printf("node id: %s\n", node.nodeId)

  /*
  fmt.Printf("predecessor1 id: %s\n", node.predecessors[1].nodeId)
  fmt.Printf("predecessor0 id: %s\n", node.predecessors[0].nodeId)
  
  fmt.Printf("node id          : %s\n", node.nodeId)
  */

  fmt.Printf("successor0 id    : %s\n", node.successors[0].nodeId)
  fmt.Printf("Finger   1    key: %s\n", node.fingers[0].key)  
  fmt.Printf("Finger   1 nodeId: %s\n", node.fingers[0].nodeId)
  fmt.Printf("Finger   3    key: %s\n", node.fingers[2].key)
  fmt.Printf("Finger   3 nodeId: %s\n", node.fingers[2].nodeId) 
  fmt.Printf("Finger  80    key: %s\n", node.fingers[79].key)
  fmt.Printf("Finger  80 nodeId: %s\n", node.fingers[79].nodeId)
  fmt.Printf("Finger 130    key: %s\n", node.fingers[129].key)
  fmt.Printf("Finger 130 nodeId: %s\n", node.fingers[129].nodeId)
  fmt.Printf("Finger 155    key: %s\n", node.fingers[154].key)
  fmt.Printf("Finger 155 nodeId: %s\n", node.fingers[154].nodeId)
  fmt.Printf("Finger 160    key: %s\n", node.fingers[159].key)
  fmt.Printf("Finger 160 nodeId: %s\n", node.fingers[159].nodeId)
  
  /*
  fmt.Printf("nodeIp: %s\n", node.nodeIp)
  fmt.Printf("nodePort: %s\n", node.nodePort)
  fmt.Printf("joinViaIp: %s\n", node.joinViaIp)
  fmt.Printf("joinViaPort: %s\n\n", node.joinViaPort)
  fmt.Printf("successor0 id: %s\n", node.successors[0].nodeId)
  fmt.Printf("successor1 id: %s\n", node.successors[1].nodeId)
  fmt.Printf("successorIp: %s\n", self.successors[0].nodeIp)
  fmt.Printf("successorPort: %s\n\n", self.successors[0].nodePort)
  */

  fmt.Printf("\n")
}


func (self *DHTnode) printRing(){
	fmt.Println("\nNodes in ring:")
 	start := self
 	  
    printNode(self)

  	self = self.successor
  	for(self != start){

      printNode(self)

      self = self.successor
  	}
 }