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
	FanSpeed *FanSpeed
	// DesiredRoomTemperature in C (5-40) times 10
	DesiredRoomTemperature *int
	// DesiredDHWTemperature in C (10-60) times 10
	DesiredDHWTemperature *int
	// DHWProductionPaused tells if DHW is switched temporaty off
	DHWProductionPaused *bool
	// DHWProductionPauseDuration is duration of DHW pause (1-180)
	DHWProductionPauseDuration *int
	// CentralHeatingPaused tells if central heating is switched temporary off
	CentralHeatingPaused *bool
	// CentralHeatingPauseDuration is duration of central heating pause (1-180)
	CentralHeatingPauseDuration *int
	// VentilationMode is either 0 (Auto), 1 (Cooling) or 2 (Heating)
	VentilationMode *int
	// VentilationOnPause is used for stopping ventilation (emergency)
	VentilationOnPause *bool
	// SetpointSupplyTemperature in C (5-50)
	SetpointSupplyTemperature *int
}

// Readings from Nilan sensors
type Readings struct {
	// RoomTemperature in C times 10
	RoomTemperature int
	// OutdoorTemperature in C times 10
	OutdoorTemperature int
	// AverageHumidity (0-100%)
	AverageHumidity int
	// ActualHumidity of air (0-100%)
	ActualHumidity int
	// DHWTankTopTempeature in C times 10
	DHWTankTopTemperature int
	// DHWTankBottomTemperature in C times 10
	DHWTankBottomTemperature int
	// SupplyFlowTemperature in C times 10
	SupplyFlowTemperature int
}
