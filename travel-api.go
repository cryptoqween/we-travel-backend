package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	//"strconv"

	"github.com/gorilla/mux"
)

type PathRequestBody struct {
	FromLocation [2]float64 `json:"fromLocation"`
	ToLocation   [2]float64 `json:"toLocation"`
}

type PropertyOutput struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type GeometryOutput struct {
	Type        string       `json:"type"`
	Coordinates [][2]float64 `json:"coordinates"`
}

type FeatureOutput struct {
	Type       string         `json:"type"`
	ID         string         `json:"id"`
	Properties PropertyOutput `json:"properties"`
	Geometry   GeometryOutput `json:"geometry"`
}

type GeoJsonOutput struct {
	Type     string          `json:"type"`
	Features []FeatureOutput `json:"features"`
}

func getCoordinates(nodes []Node) []Coordinate {
	result := make([][2]float64, len(nodes))
	for i := 0; i < len(nodes); i++ {
		result[i] = nodes[i].value
	}
	return result
}

func findpathHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	var newRequestBody PathRequestBody
	if err != nil {
		fmt.Fprintf(w, "Bad input")
	}

	json.Unmarshal(reqBody, &newRequestBody)

	nodes := calculatePath(newRequestBody.FromLocation, newRequestBody.ToLocation)
	coords := getCoordinates(nodes)

	geometry := GeometryOutput{Type: "LineString", Coordinates: coords}
	featureOut := FeatureOutput{Type: "Feature", ID: "1234", Properties: PropertyOutput{}, Geometry: geometry}
	features := make([]FeatureOutput, 1)

	features[0] = featureOut

	geojsonData := GeoJsonOutput{
		Type:     "FeatureCollection",
		Features: features,
	}

	geojsonDataInJson, _ := json.Marshal(&geojsonData)
	w.Write(geojsonDataInJson)
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Lightpath backend is running")
}

func main() {
	loadGeoJSON()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/findpath", findpathHandler).Methods("POST")

	fmt.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
