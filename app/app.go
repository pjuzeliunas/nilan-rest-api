package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"../controller"
)

func readings(w http.ResponseWriter, r *http.Request) {
	log.Printf("Processing readings GET request from %v\n", r.RemoteAddr)
	readings := controller.FetchReadings()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(readings)
}

func settings(w http.ResponseWriter, r *http.Request) {
	log.Printf("Processing settings GET request from %v\n", r.RemoteAddr)
	settings := controller.FetchSettings()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func updateSettings(w http.ResponseWriter, r *http.Request) {
	log.Printf("Processing settings update request from %v\n", r.RemoteAddr)
	var newSettings controller.Settings
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Bad request: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &newSettings)
	if err != nil {
		log.Printf("Bad request: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	controller.SendSettings(newSettings)
	w.WriteHeader(http.StatusOK)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/readings", readings).Methods("GET")
	router.HandleFunc("/settings", settings).Methods("GET")
	router.HandleFunc("/settings", updateSettings).Methods("PUT")
	log.Println("Listening at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
