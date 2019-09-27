package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type event struct {
	ID          int    `json: "id"`
	Title       string `json: "title"`
	Description string `json: "description"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          1,
		Title:       "Intro to golang",
		Description: "BLablablabla balbalbala",
	},
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Bad input")
	}

	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}

func getEvent(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
	// TODO: handle error when no event found for ID
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1])
			fmt.Fprintf(w, "Event deleted: %v", eventID)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome Home")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getEvent).Methods("GET")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
