package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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

var graph ItemGraph

func loadGeoJSON() {
	jsonFile, err := os.Open("./central.geojson") //GeoJSON for central london around highbury islington
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
		isValidAccess := !isPath || isValidPath
		isLit := geojson.Features[i].Properties.Lit != "" || geojson.Features[i].Properties.Lit == "yes"
		if isHighway && isLineString && hasSidewalk && isValidAccess && isLit {
			var feature = geojson.Features[i]
			var prev *Node
			//fmt.Println("ID: " + feature.Properties.Highway)
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

	defer jsonFile.Close()
	fmt.Println("geojson explorer running")
}

func calculatePath(startCoords [2]float64, endCoords [2]float64) []Node {
	nodeStart := graph.FindNode(startCoords)
	nodeEnd := graph.FindNode(endCoords)
	pathFound := graph.FindPath(graph.nodes[nodeStart], graph.nodes[nodeEnd])
	fmt.Println(pathFound)
	return pathFound
}
