package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Property struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Highway  string `json:"highway"`
	Access   string `json:"access"`
	Lit      string `json:"lit"`
	Sidewalk string `json:"sidewalk"`
}

type Coordinate = [2]float64

type Geometry struct {
	Type        string       `json:"type"`
	Coordinates []Coordinate `json:"coordinates"`
}

type Feature struct {
	Type       string   `json:"type"`
	ID         string   `json:"id"`
	Properties Property `json:"properties"`
	Geometry   Geometry `json:"geometry"`
}

type GeoJson struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

var graph Graph

func loadGeoJSON() {
	// GeoJSON for central london around highbury islington
	// jsonFile, err := os.Open("./central.geojson")
	// GeoJSON for central london around highbury islington
	jsonFile, err := os.Open("./data/greater-london-latest.geojson")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened geojson")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var geojson GeoJson
	json.Unmarshal(byteValue, &geojson)

	for i := 0; i < len(geojson.Features); i++ {
		isLineString := geojson.Features[i].Geometry.Type == "LineString"
		isHighway := geojson.Features[i].Properties.Highway != ""
		hasSidewalk := geojson.Features[i].Properties.Sidewalk != "" || geojson.Features[i].Properties.Sidewalk != "none"
		isPath := geojson.Features[i].Properties.Highway == "path"
		isValidPath := geojson.Features[i].Properties.Highway == "path" && (geojson.Features[i].Properties.Access == "no" || geojson.Features[i].Properties.Access == "private")
		isNotPathOrIsValidPath := !isPath || isValidPath
		isLit := geojson.Features[i].Properties.Lit != "" || geojson.Features[i].Properties.Lit == "yes"
		if isHighway && isLineString && hasSidewalk && isNotPathOrIsValidPath && isLit {
			var feature = geojson.Features[i]
			var prev *Node
			for j := 0; j < len(feature.Geometry.Coordinates); j++ {
				var coords = feature.Geometry.Coordinates[j]
				node := CreateNode(coords)
				graph.AddNode(&node)
				if j != 0 {
					graph.AddEdge(&node, prev)
					graph.AddEdge(prev, &node)
				}
				prev = &node
			}
		}
	}

	defer jsonFile.Close()
	fmt.Println("geojson Graph created with %d nodes", len(graph.nodes))
}

func calculatePath(startCoords Coordinate, endCoords Coordinate) []Node {
	nodeStart := graph.FindNode(startCoords)
	nodeEnd := graph.FindNode(endCoords)
	pathFound := graph.FindPath(graph.nodes[nodeStart], graph.nodes[nodeEnd])
	fmt.Println(pathFound)
	return pathFound
}
