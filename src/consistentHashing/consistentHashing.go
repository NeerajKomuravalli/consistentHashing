package consistentHashing

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

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

func (ring *Ring) AddNode(id string) {
	ring.Lock()
	defer ring.Unlock()

	newNode := NewNode(id)
	ring.Nodes = append(ring.Nodes, newNode)
	sort.Sort(ring.Nodes)
}

func (ring *Ring) RemoveNode(id string) error {
	findNode := func(i int) bool {
		return ring.Nodes[i].Id == id
	}
	ring.RLock()
	index := sort.Search(ring.Nodes.Len(), findNode)
	if index >= ring.Nodes.Len() {
		return fmt.Errorf("no node with id %v was found", id)
	}
	ring.RUnlock()
	ring.Lock()
	defer ring.Unlock()
	ring.Nodes = append(ring.Nodes[:index], ring.Nodes[index+1:]...)
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
		return nodes[i].HashId >= hashId(id)
	}
	return sort.Search(len(nodes), searchFn)
}

// Creates a new node
func NewNode(id string) *Node {
	return &Node{id, hashId(id)}
}

// Calculates hash for a given Id
func hashId(id string) uint32 {
	return crc32.ChecksumIEEE([]byte(id))
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
