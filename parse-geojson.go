package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/cheekybits/genny/generic"
)

// Item the type of the binary search tree
type Item generic.Type

// Node a single node that composes the tree
type Node struct {
	value Item
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

type Property struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Highway  string `json:"highway"`
	Access   string `json:"access"`
	Lit      string `json:"lit"`
	Sidewalk string `json:"sidewalk"`
}

type Coordinate = [2]float64

type Feature struct {
	Type       string   `json:"type"`
	ID         string   `json:"id"`
	Properties Property `json:"properties"`
	Geometry   struct {
		Type        string       `json:"type"`
		Coordinates []Coordinate `json:"coordinates"`
	} `json:"geometry"`
}

type GeoJson struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

func main() {
	jsonFile, err := os.Open("./central.geojson")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var geojson GeoJson
	json.Unmarshal(byteValue, &geojson)

	var graph ItemGraph

	for i := 0; i < len(geojson.Features); i++ {
		isLineString := geojson.Features[i].Geometry.Type == "LineString"
		isHighway := geojson.Features[i].Properties.Highway != ""
		hasSidewalk := geojson.Features[i].Properties.Sidewalk != "" || geojson.Features[i].Properties.Sidewalk != "none"
		isPath := geojson.Features[i].Properties.Highway == "path"
		isValidPath := geojson.Features[i].Properties.Highway == "path" && (geojson.Features[i].Properties.Access == "no" || geojson.Features[i].Properties.Access == "private")
		isValidAccess := !isPath || isValidPath
		isLit := geojson.Features[i].Properties.Lit != "" || geojson.Features[i].Properties.Lit == "yes"
		if isHighway && isLineString && hasSidewalk && isValidAccess && isLit {
			var feature = geojson.Features[i]
			var prev *Node
			fmt.Println("ID: " + feature.Properties.Highway)
			for j := 0; j < len(feature.Geometry.Coordinates); j++ {
				var coords = feature.Geometry.Coordinates[j]
				var hash = strconv.FormatFloat(coords[0], 'f', -1, 64)
				hash += strconv.FormatFloat(coords[1], 'f', -1, 64)
				node := Node{coords}
				graph.AddNode(&node)
				if j != 0 {
					graph.AddEdge(&node, prev)
					graph.AddEdge(prev, &node)
				}
				prev = &node
			}
		}
	}
	fmt.Println(graph.FindPath(graph.nodes[0], graph.nodes[2]))

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	fmt.Println("running")

}
