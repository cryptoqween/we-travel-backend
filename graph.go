package main

import (
	"container/heap"
	"fmt"
	"math"
	"strconv"

	//"math"
	"sync"
	//"github.com/cheekybits/genny/generic"
)

// Node a single node that composes the tree
type Node struct {
	Value Coordinate
}

func CreateNode(coords Coordinate) Node {
	var hash = strconv.FormatFloat(coords[0], 'f', -1, 64)
	hash += strconv.FormatFloat(coords[1], 'f', -1, 64)
	return Node{coords}
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.Value)
}

type Graph struct {
	nodes []*Node
	edges map[Node][]*Node
	lock  sync.RWMutex
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(n *Node) {
	g.lock.Lock()
	g.nodes = append(g.nodes, n)
	g.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *Graph) AddEdge(n1, n2 *Node) {
	g.lock.Lock()
	if g.edges == nil {
		g.edges = make(map[Node][]*Node)
	}
	g.edges[*n1] = append(g.edges[*n1], n2)
	g.edges[*n2] = append(g.edges[*n2], n1)
	g.lock.Unlock()
}

// Print graph
func (g *Graph) String() {
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

func (g *Graph) FindNode(coords Coordinate) int {
	nodes := g.nodes
	//min distance
	var foundNodeIndex int
	var minDistance float64
	//index of node
	// Probably there is a better algo for this, just doing the brute force sorry :(
	for i := 0; i < len(nodes); i++ {
		//calc distance
		value := nodes[i].Value
		dx := coords[0] - value[0]
		dy := coords[1] - value[1]
		distance := math.Sqrt(dx*dx + dy*dy)
		if i == 0 || distance < minDistance {
			foundNodeIndex = i
			minDistance = distance
		}
	}
	return foundNodeIndex
}

// A* routing
// heap has child nodes sorted by distance
func (g *Graph) FindPath(src, dest *Node) []Node {
	g.lock.RLock()
	pqueue := make(PriorityQueue, 1)
	rootPath := []Node{}
	rootItem := QueueItem{*src, rootPath}
	pqueue[0] = &Item{
		Value:    &rootItem,
		Priority: 0,
		Index:    0,
	}
	heap.Init(&pqueue)
	visited := make(map[*Node]bool)
	for {
		if pqueue.Len() == 0 {
			break
		}
		pqitem := pqueue.Pop().(*Item)
		value := pqitem.Value
		node := value.Node
		visited[&node] = true
		children := g.edges[node]

		for i := 0; i < len(children); i++ {
			child := children[i]

			if *child == *dest {
				fmt.Println("Found Dest with distance", pqitem.Priority)
				return value.Path
			}

			if !visited[child] {
				path := append(value.Path, *child)
				queueItem := QueueItem{*child, path}
				dx := (node.Value[0] - child.Value[0])
				dy := (node.Value[1] - child.Value[1])
				distance := math.Sqrt(dx*dx + dy*dy)
				newItem := Item{
					Value:    &queueItem,
					Priority: distance,
				}
				heap.Push(&pqueue, &newItem)
				visited[child] = true
			}
		}
	}
	g.lock.RUnlock()
	// No path
	return []Node{}
}
