package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"../controller"
)

func readings(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(controller.FetchReadings())
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/readings", readings)
	log.Fatal(http.ListenAndServe(":8080", router))
}

// func main() {
// 	var settings = controller.FetchSettings()
// 	settings.FanSpeed = controller.FanSpeedLow
// 	controller.SendSettings(settings)
// }
