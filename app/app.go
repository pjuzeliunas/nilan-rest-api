package main

import (
	"fmt"

	"../controller"
)

func main() {
	controller.FetchReadings()
	fmt.Println("Readings: ", controller.FetchReadings())
	fmt.Println("Settings: ", controller.FetchSettings())
}
