package consistentHashing

import (
	"fmt"
	"sort"
	"sync"
)

var hashFuncMaps = GetHashFuncsMap()

type Node struct {
	Id     string
	HashId uint32
}

type Nodes []*Node

type Ring struct {
	Nodes Nodes
	sync.RWMutex
}

func NewRing() *Ring {
	return &Ring{}
}

func (ring *Ring) AddNode(id string, noOfVNodes int) {
	ring.Lock()
	defer ring.Unlock()
	i := 0
	for hashFuncKey := range hashFuncMaps {
		newNode := NewNode(id, hashFuncKey)
		ring.Nodes = append(ring.Nodes, newNode)
		i++
		if i == noOfVNodes {
			break
		}
	}
	sort.Sort(ring.Nodes)
}

func (ring *Ring) RemoveNode(id string) error {
	// This implementation is very bad as it has time complexity of O(N) but i am not
	// able to find a better solution yet.
	var indexList []int
	ring.RLock()
	for index, node := range ring.Nodes {
		if id == node.Id {
			indexList = append(indexList, index)
		}
	}
	ring.RUnlock()
	if len(indexList) == 0 {
		return fmt.Errorf("no node with id %v was found", id)
	} else {
		ring.Lock()
		for _, index := range indexList {
			ring.Nodes = append(ring.Nodes[:index], ring.Nodes[index+1:]...)
		}
	}
	ring.Unlock()
	return nil
}

func (ring *Ring) GetNodeId(id string) string {
	ring.RLock()
	index := ring.Nodes.search(id)
	ring.RUnlock()
	if index >= ring.Nodes.Len() {
		index = 0
	}
	return ring.Nodes[index].Id
}

func (nodes Nodes) search(id string) int {
	searchFn := func(i int) bool {
		return nodes[i].HashId >= hashFuncMaps["hashId"](id)
	}
	return sort.Search(len(nodes), searchFn)
}

// Creates a new node
func NewNode(id string, hashFuncKey string) *Node {
	return &Node{id, hashFuncMaps[hashFuncKey](id)}
}

// Implementing Sort method on Nodes : Len function
func (nodes Nodes) Len() int {
	return len(nodes)
}

// Implementing Sort method on Nodes : Swap function
func (nodes Nodes) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}

// Implementing Sort method on Nodes : Less function
func (node Nodes) Less(i, j int) bool {
	return node[i].HashId < node[j].HashId
}
