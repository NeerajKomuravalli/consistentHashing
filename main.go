package main

import (
	"fmt"

	"github.com/NeerajKomuravalli/consistentHashing/src/consistentHashing"
)

func print(n consistentHashing.Nodes) {
	if len(n) == 0 {
		fmt.Println("Sorted Nodes : {}")
	} else {
		fmt.Printf("Sorted Nodes : ")
		for _, v := range n {
			fmt.Printf("{Id : %v, HashId : %v}, ", v.Id, v.HashId)
		}
		fmt.Println()
	}
}

func main() {
	fmt.Println("Initializing new ring")
	ring := consistentHashing.NewRing()
	fmt.Println("Adding Node with id : 1234")
	ring.AddNode("1234")
	fmt.Println("Adding Node with id : 12344")
	ring.AddNode("12344")
	fmt.Println("Adding Node with id : 12346")
	ring.AddNode("12346")
	print(ring.Nodes)
	fmt.Println("Removing a node 1234")
	ring.RemoveNode("1234")
	print(ring.Nodes)
	fmt.Println("Removing a non existing node : 123467889")
	err := ring.RemoveNode("123467889")
	if err != nil {
		fmt.Println(err)
	} else {
		print(ring.Nodes)
	}
}
