package main

import (
	"bufio"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		{
			break
		}
	}
	nbNode, _ := strconv.Atoi(scanner.Text())
	for i := 0; i < nbNode; i++ {
		oneNode()
	}
}
