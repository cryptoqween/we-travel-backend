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

type MultiLineString struct {
	Type        string       `json:"type"`
	Coordinates [][2]float64 `json:"coordinates"`
}

func formatResponse(nodes []Node) [][2]float64 {
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
	result := formatResponse(nodes)

	geojsonData := MultiLineString{Type: "MultiLineString", Coordinates: result}
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
