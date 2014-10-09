
package main

import (
  "fmt"
)

func printNode(node *DHTnode) {

  fmt.Printf("node id: %s\n", node.Id)

  /*
  fmt.Printf("predecessor1 id: %s\n", node.predecessors[1].Id)
  fmt.Printf("predecessor0 id: %s\n", node.predecessors[0].Id)
  
  fmt.Printf("node id          : %s\n", node.Id)
  */

  //fmt.Printf("Successor0 id    : %s\n", node.Successors[0].Id)
  /*
  fmt.Printf("Finger   1    key: %s\n", node.fingers[0].key)   
  fmt.Printf("Finger   1 Id: %s\n\n", node.fingers[0].Id)
  fmt.Printf("Finger   2    key: %s\n", node.fingers[1].key)
  fmt.Printf("Finger   2 Id: %s\n\n", node.fingers[1].Id) 
  
  
  fmt.Printf("Finger   3    key: %s\n", node.fingers[2].key) 
  fmt.Printf("Finger   3 Id: %s\n\n", node.fingers[2].Id) 
  fmt.Printf("Finger   4    key: %s\n", node.fingers[3].key)
  fmt.Printf("Finger   4 Id: %s\n\n", node.fingers[3].Id) 
  fmt.Printf("Finger   5    key: %s\n", node.fingers[4].key)
  fmt.Printf("Finger   5 Id: %s\n\n", node.fingers[4].Id) 
  fmt.Printf("Finger   6    key: %s\n", node.fingers[5].key)  
  fmt.Printf("Finger   6 Id: %s\n\n", node.fingers[5].Id)
  fmt.Printf("Finger   7    key: %s\n", node.fingers[6].key)
  fmt.Printf("Finger   7 Id: %s\n\n", node.fingers[6].Id) 
  fmt.Printf("Finger   8    key: %s\n", node.fingers[7].key)
  fmt.Printf("Finger   8 Id: %s\n\n", node.fingers[7].Id) 
  fmt.Printf("Finger   9    key: %s\n", node.fingers[8].key)
  fmt.Printf("Finger   9 Id: %s\n\n", node.fingers[8].Id) 
  fmt.Printf("Finger  10    key: %s\n", node.fingers[9].key)
  fmt.Printf("Finger  10 Id: %s\n\n", node.fingers[9].Id) 
  

  
  fmt.Printf("Finger  80    key: %s\n", node.fingers[79].key)
  fmt.Printf("Finger  80 Id: %s\n\n", node.fingers[79].Id)
  fmt.Printf("Finger 130    key: %s\n", node.fingers[129].key)
  fmt.Printf("Finger 130 Id: %s\n\n", node.fingers[129].Id)
  fmt.Printf("Finger 155    key: %s\n", node.fingers[154].key)
  fmt.Printf("Finger 155 Id: %s\n\n", node.fingers[154].Id)
  fmt.Printf("Finger 160    key: %s\n", node.fingers[159].key)
  fmt.Printf("Finger 160 Id: %s\n\n", node.fingers[159].Id)
  */

  /*
  fmt.Printf("nodeIp: %s\n", node.nodeIp)
  fmt.Printf("NodePort: %s\n", node.NodePort)
  fmt.Printf("joinViaIp: %s\n", node.joinViaIp)
  fmt.Printf("joinViaPort: %s\n\n", node.joinViaPort)
  fmt.Printf("Successor0 id: %s\n", node.Successors[0].Id)
  fmt.Printf("Successor1 id: %s\n", node.Successors[1].Id)
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