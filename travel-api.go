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
	FromLocation Coordinate `json:"fromLocation"`
	ToLocation   Coordinate `json:"toLocation"`
}

func getCoordinates(nodes []Node) []Coordinate {
	result := make([]Coordinate, len(nodes))
	for i := 0; i < len(nodes); i++ {
		result[i] = nodes[i].Value
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

	geometry := Geometry{Type: "LineString", Coordinates: coords}
	featureOut := Feature{Type: "Feature", ID: "1234", Properties: Property{}, Geometry: geometry}
	features := make([]Feature, 1)

	features[0] = featureOut

	geojsonData := GeoJson{
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

	// terminate if killed
	// go func() {
	//     if err := server.ListenAndServe(":8080", router); err != nil {
	// 		// handle err
	// 		log.Fatal(err)
	//     }
	// }()

	// // Setting up signal capturing
	// stop := make(chan os.Signal, 1)
	// signal.Notify(stop, os.Interrupt)

	// // Waiting for SIGINT (pkill -2)
	// <-stop

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := server.Shutdown(ctx); err != nil {
	//     // handle err
	// }

}
