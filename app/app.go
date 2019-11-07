package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"../controller"
)

func readings(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(controller.FetchReadings())
}

func settings(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(controller.FetchSettings())
}

func updateSettings(w http.ResponseWriter, r *http.Request) {
	var newSettings controller.Settings
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please verify data")
		return
	}

	json.Unmarshal(reqBody, &newSettings)

	controller.SendSettings(newSettings)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSettings)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/readings", readings).Methods("GET")
	router.HandleFunc("/settings", settings).Methods("GET")
	router.HandleFunc("/settings", updateSettings).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
