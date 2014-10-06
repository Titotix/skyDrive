
package main

import (
  "fmt"
)

func printNode(node *DHTnode) {

  fmt.Printf("node id: %s\n", node.NodeId)

  /*
  fmt.Printf("predecessor1 id: %s\n", node.predecessors[1].NodeId)
  fmt.Printf("predecessor0 id: %s\n", node.predecessors[0].NodeId)
  
  fmt.Printf("node id          : %s\n", node.NodeId)
  */

  //fmt.Printf("Successor0 id    : %s\n", node.Successors[0].NodeId)
  /*
  fmt.Printf("Finger   1    key: %s\n", node.fingers[0].key)   
  fmt.Printf("Finger   1 NodeId: %s\n\n", node.fingers[0].NodeId)
  fmt.Printf("Finger   2    key: %s\n", node.fingers[1].key)
  fmt.Printf("Finger   2 NodeId: %s\n\n", node.fingers[1].NodeId) 
  
  
  fmt.Printf("Finger   3    key: %s\n", node.fingers[2].key) 
  fmt.Printf("Finger   3 NodeId: %s\n\n", node.fingers[2].NodeId) 
  fmt.Printf("Finger   4    key: %s\n", node.fingers[3].key)
  fmt.Printf("Finger   4 NodeId: %s\n\n", node.fingers[3].NodeId) 
  fmt.Printf("Finger   5    key: %s\n", node.fingers[4].key)
  fmt.Printf("Finger   5 NodeId: %s\n\n", node.fingers[4].NodeId) 
  fmt.Printf("Finger   6    key: %s\n", node.fingers[5].key)  
  fmt.Printf("Finger   6 NodeId: %s\n\n", node.fingers[5].NodeId)
  fmt.Printf("Finger   7    key: %s\n", node.fingers[6].key)
  fmt.Printf("Finger   7 NodeId: %s\n\n", node.fingers[6].NodeId) 
  fmt.Printf("Finger   8    key: %s\n", node.fingers[7].key)
  fmt.Printf("Finger   8 NodeId: %s\n\n", node.fingers[7].NodeId) 
  fmt.Printf("Finger   9    key: %s\n", node.fingers[8].key)
  fmt.Printf("Finger   9 NodeId: %s\n\n", node.fingers[8].NodeId) 
  fmt.Printf("Finger  10    key: %s\n", node.fingers[9].key)
  fmt.Printf("Finger  10 NodeId: %s\n\n", node.fingers[9].NodeId) 
  

  
  fmt.Printf("Finger  80    key: %s\n", node.fingers[79].key)
  fmt.Printf("Finger  80 NodeId: %s\n\n", node.fingers[79].NodeId)
  fmt.Printf("Finger 130    key: %s\n", node.fingers[129].key)
  fmt.Printf("Finger 130 NodeId: %s\n\n", node.fingers[129].NodeId)
  fmt.Printf("Finger 155    key: %s\n", node.fingers[154].key)
  fmt.Printf("Finger 155 NodeId: %s\n\n", node.fingers[154].NodeId)
  fmt.Printf("Finger 160    key: %s\n", node.fingers[159].key)
  fmt.Printf("Finger 160 NodeId: %s\n\n", node.fingers[159].NodeId)
  */

  /*
  fmt.Printf("nodeIp: %s\n", node.nodeIp)
  fmt.Printf("NodePort: %s\n", node.NodePort)
  fmt.Printf("joinViaIp: %s\n", node.joinViaIp)
  fmt.Printf("joinViaPort: %s\n\n", node.joinViaPort)
  fmt.Printf("Successor0 id: %s\n", node.Successors[0].NodeId)
  fmt.Printf("Successor1 id: %s\n", node.Successors[1].NodeId)
  fmt.Printf("SuccessorIp: %s\n", self.Successors[0].nodeIp)
  fmt.Printf("SuccessorPort: %s\n\n", self.Successors[0].NodePort)
  */

  //fmt.Printf("\n")
}


func (self *DHTnode) printRing(){
  fmt.Println("\nNodes in ring:")
  start := self
  printNode(self)
  
  if (self.Successor != nil) {
    self = self.Successor
    for(self != start){
      printNode(self)
      self = self.Successor
    }
  } 
}