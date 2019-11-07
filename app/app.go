package main

import (
	"../controller"
)

func main() {
	var settings = controller.FetchSettings()
	settings.FanSpeed = controller.FanSpeedLow
	controller.SendSettings(settings)
}
