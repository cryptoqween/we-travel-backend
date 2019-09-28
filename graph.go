package main

import (
	"fmt"
	"math"

	//"math"
	"sync"

	"github.com/cheekybits/genny/generic"
	//"github.com/cheekybits/genny/generic"
)

// Item the type of the binary search tree
type Item generic.Type

// Node a single node that composes the tree
type Node struct {
	value [2]float64
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.value)
}

// ItemGraph the Items graph
type ItemGraph struct {
	nodes []*Node
	edges map[Node][]*Node
	lock  sync.RWMutex
}

// AddNode adds a node to the graph
func (g *ItemGraph) AddNode(n *Node) {
	g.lock.Lock()
	g.nodes = append(g.nodes, n)
	g.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *ItemGraph) AddEdge(n1, n2 *Node) {
	g.lock.Lock()
	if g.edges == nil {
		g.edges = make(map[Node][]*Node)
	}
	g.edges[*n1] = append(g.edges[*n1], n2)
	g.edges[*n2] = append(g.edges[*n2], n1)
	g.lock.Unlock()
}

// Print graph
func (g *ItemGraph) String() {
	g.lock.RLock()
	s := ""
	for i := 0; i < len(g.nodes); i++ {
		s += g.nodes[i].String() + " -> "
		near := g.edges[*g.nodes[i]]
		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
		}
		s += "\n"
	}
	fmt.Println(s)
	g.lock.RUnlock()
}

type QueueItem struct {
	Node Node
	Path []Node
}

type NodeQueue struct {
	items []QueueItem
	lock  sync.RWMutex
}

// New creates a new NodeQueue
func (s *NodeQueue) New() *NodeQueue {
	s.lock.Lock()
	s.items = []QueueItem{}
	s.lock.Unlock()
	return s
}

// Enqueue adds an Node to the end of the queue
func (s *NodeQueue) Enqueue(t QueueItem) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Dequeue removes an Node from the start of the queue
func (s *NodeQueue) Dequeue() *QueueItem {
	s.lock.Lock()
	item := s.items[0]
	s.items = s.items[1:len(s.items)]
	s.lock.Unlock()
	return &item
}

// Front returns the item next in the queue, without removing it
func (s *NodeQueue) Front() *QueueItem {
	s.lock.RLock()
	item := s.items[0]
	s.lock.RUnlock()
	return &item
}

// IsEmpty returns true if the queue is empty
func (s *NodeQueue) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items) == 0
}

// Size returns the number of Nodes in the queue
func (s *NodeQueue) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}

func (g *ItemGraph) FindNode(coords [2]float64) int {
	nodes := g.nodes
	//min distance
	var foundNodeIndex int
	var minDistance float64
	//index of node
	// Probably there is a better algo for this, just doing the brute force sorry :(
	for i := 0; i < len(nodes); i++ {
		//calc distance
		distance := math.Sqrt((coords[0]-nodes[i].value[0])*(coords[0]-nodes[i].value[0]) + (coords[1]-nodes[i].value[1])*(coords[1]-nodes[i].value[1]))
		if i == 0 || distance < minDistance {
			foundNodeIndex = i
			minDistance = distance
		}
	}
	return foundNodeIndex
}

// A* routing
// heap has child nodes sorted by distance
func (g *ItemGraph) FindPath(src, dest *Node) []Node {
	g.lock.RLock()
	q := NodeQueue{}
	q.New()
	item := QueueItem{*src, []Node{}}
	q.Enqueue(item)
	visited := make(map[*Node]bool)
	for {
		if q.IsEmpty() {
			break
		}
		item := q.Dequeue()
		node := item.Node
		visited[&node] = true
		near := g.edges[node]

		for i := 0; i < len(near); i++ {
			j := near[i]

			if *j == *dest {
				fmt.Println("Found Dest")
				return item.Path
			}

			if !visited[j] {
				path := append(item.Path, *j)
				item := QueueItem{*j, path}
				q.Enqueue(item)
				visited[j] = true
			}
		}
	}
	g.lock.RUnlock()
	// No path
	return []Node{}
}
