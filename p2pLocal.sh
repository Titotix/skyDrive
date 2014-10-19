#!/bin/bash
i=0

while [ $i != $1 ]
do

    (( port = $i+5555 ))
    go run util.go output.go standalone_node.go node.go network.go routing.go finger.go message.go rpc.go <<jetenique
    5555
jetenique
    (( i=$i+1 ))
done
