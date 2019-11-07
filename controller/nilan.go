package controller

// FanSpeed represents Nilan ventilation intensity value in range from 101 (lowest) to 104 (highest).
type FanSpeed uint16

const (
	// FanSpeedLow represents lowest fan speed aka level 1
	FanSpeedLow FanSpeed = 101
	// FanSpeedNormal represents normal fan speed aka level 2
	FanSpeedNormal FanSpeed = 102
	// FanSpeedHigh represents high fan speed aka level 3
	FanSpeedHigh FanSpeed = 103
	// FanSpeedVeryHigh represents highest fan speed aka level 4
	FanSpeedVeryHigh FanSpeed = 104
)

// Settings of Nilan system
type Settings struct {
	// FanSpeed of ventilation
	FanSpeed FanSpeed
	// DesiredRoomTemperature in C (5-40) times 10
	DesiredRoomTemperature int
}

// Readings from Nilan sensors
type Readings struct {
	// RoomTemperature in C times 10
	RoomTemperature int
	// OutdoorTemperature in C times 10
	OutdoorTemperature int
}
