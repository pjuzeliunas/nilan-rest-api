package controller

import (
	"encoding/binary"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/goburrow/modbus"
)

func getHandler(slaveID byte) *modbus.TCPClientHandler {
	// Modbus TCP
	handler := modbus.NewTCPClientHandler("192.168.5.107:502")
	handler.Timeout = 10 * time.Second
	handler.SlaveId = slaveID
	err := handler.Connect()

	if err != nil {
		panic(err)
	}

	return handler
}

func fetchValue(slaveID byte, register Register) uint16 {
	handler := getHandler(slaveID)
	defer handler.Close()
	client := modbus.NewClient(handler)
	resultBytes, _ := client.ReadHoldingRegisters(uint16(register), 1)
	if len(resultBytes) == 2 {
		return binary.BigEndian.Uint16(resultBytes)
	}
	panic("Cannot read register value")
}

func fetchRegisterValues(slaveID byte, registers []Register) map[Register]uint16 {
	m := make(map[Register]uint16)

	handler := getHandler(slaveID)
	defer handler.Close()
	client := modbus.NewClient(handler)

	for _, register := range registers {
		resultBytes, _ := client.ReadHoldingRegisters(uint16(register), 1)
		if len(resultBytes) == 2 {
			resultWord := binary.BigEndian.Uint16(resultBytes)
			m[register] = resultWord
		}
	}

	return m
}

func setRegisterValues(slaveID byte, values map[Register]uint16) {
	handler := getHandler(slaveID)
	defer handler.Close()
	client := modbus.NewClient(handler)

	for register, value := range values {
		client.WriteSingleRegister(uint16(register), value)
	}
}

// Register is address of register on client
type Register uint16

const (
	// FanSpeedRegister is ID of register holding desired FanSpeed value
	FanSpeedRegister Register = 20148
	// DesiredRoomTemperatureRegister is ID of register holding desired room temperature in C times 10.
	// Example: 23.5 C is stored as 235.
	DesiredRoomTemperatureRegister Register = 20260
	// MasterTemperatureSensorSettingRegister is ID of register holding either 0 (read temperature from T3)
	// or 1 (read temperature from Text)
	MasterTemperatureSensorSettingRegister Register = 20263
	// T3ExtractAirTemperatureRegister is ID of register holding room temperature value when
	// MasterTemperatureSensorSettingRegister is 0
	T3ExtractAirTemperatureRegister Register = 20286
	// TextRoomTemperatureRegister is ID of register holding room temperature value when
	// MasterTemperatureSensorSettingRegister is 1
	TextRoomTemperatureRegister Register = 20280
	// OutdoorTemperatureRegister is ID of register outdoor temperature
	OutdoorTemperatureRegister Register = 20282
	// AverageHumidityRegister is ID of register holding average humidity value
	AverageHumidityRegister Register = 20164
	// ActualHumidityRegister is ID of register holding actual humidity value
	ActualHumidityRegister Register = 21776
	// DHWTopTankTemperatureRegister is ID of register holding T11 top DHW tank temperature
	DHWTopTankTemperatureRegister Register = 20520
	// DHWBottomTankTemperatureRegister is ID of register holding T11 bottom DHW tank temperature
	DHWBottomTankTemperatureRegister Register = 20522
	// DHWSetPointRegister is ID of register holding desired DHW temperature
	DHWSetPointRegister Register = 20460
	// DHWPauseRegister is ID of register holding DHW pause flag
	DHWPauseRegister Register = 20440
	// DHWPauseDurationRegister is ID of register holding DHW pause duration value
	DHWPauseDurationRegister Register = 20441
	// CentralHeatingPauseRegister is ID of register holding central heating pause flag
	CentralHeatingPauseRegister Register = 20600
	// CentralHeatingPauseDurationRegister is ID of register holding central heating pause duration value
	CentralHeatingPauseDurationRegister Register = 20601
	// VentilationModeRegister is ID of register holding ventilation mode value (0, 1 or 2).
	VentilationModeRegister Register = 20120
	// VentilationPauseRegister is ID of register holding ventilation pause flag
	VentilationPauseRegister Register = 20100
	// SetpointSupplyTemperatureRegisterAIR9 is ID of register holding setpoint supply temperature
	// on AIR9 models
	SetpointSupplyTemperatureRegisterAIR9 Register = 20680
	// SetpointSupplyTemperatureRegisterGEO is ID of register holding setpoint supply temperature
	// on GEO models
	SetpointSupplyTemperatureRegisterGEO Register = 20640
	// DeviceTypeGEOReigister is ID of register that holds number 8 on GEO models
	DeviceTypeGEOReigister Register = 21839
	// DeviceTypeAIR9Register is ID of register that holds number 9 on AIR9 models
	DeviceTypeAIR9Register Register = 21899
	// T18ReadingRegisterGEO is ID of register holding T18 supply flow temperature reading
	// on GEO models
	T18ReadingRegisterGEO Register = 20653
	// T18ReadingRegisterAIR9 is ID of register holding T18 supply flow temperature reading
	// on AIR9 models
	T18ReadingRegisterAIR9 Register = 20686
)

func supplyFlowSetpointTemperatureRegister() Register {
	switch {
	case fetchValue(4, DeviceTypeGEOReigister) == 8:
		return SetpointSupplyTemperatureRegisterGEO
	case fetchValue(4, DeviceTypeAIR9Register) == 9:
		return SetpointSupplyTemperatureRegisterAIR9
	default:
		panic("Cannot determine device type")
	}
}

// FetchSettings of Nilan
func FetchSettings() Settings {
	supplyTemperatureRegister := supplyFlowSetpointTemperatureRegister()

	client1Registers := []Register{
		FanSpeedRegister,
		DesiredRoomTemperatureRegister,
		DHWSetPointRegister,
		DHWPauseRegister,
		DHWPauseDurationRegister,
		VentilationModeRegister,
		VentilationPauseRegister}
	client4Registers := []Register{
		CentralHeatingPauseRegister,
		CentralHeatingPauseDurationRegister,
		supplyTemperatureRegister}

	client1RegisterValues := fetchRegisterValues(1, client1Registers)
	client4RegisterValues := fetchRegisterValues(4, client4Registers)

	fanSpeed := new(FanSpeed)
	*fanSpeed = FanSpeed(client1RegisterValues[FanSpeedRegister])

	desiredRoomTemperature := new(int)
	*desiredRoomTemperature = int(client1RegisterValues[DesiredRoomTemperatureRegister])

	desiredDHWTemperature := new(int)
	*desiredDHWTemperature = int(client1RegisterValues[DHWSetPointRegister])

	dhwPaused := new(bool)
	*dhwPaused = client1RegisterValues[DHWPauseRegister] == 1

	dhwPauseDuration := new(int)
	*dhwPauseDuration = int(client1RegisterValues[DHWPauseDurationRegister])

	centralHeatingPaused := new(bool)
	*centralHeatingPaused = client4RegisterValues[CentralHeatingPauseRegister] == 1

	centralHeatingPauseDuration := new(int)
	*centralHeatingPauseDuration = int(client4RegisterValues[CentralHeatingPauseDurationRegister])

	ventilationMode := new(int)
	*ventilationMode = int(client1RegisterValues[VentilationModeRegister])

	ventilationPause := new(bool)
	*ventilationPause = client1RegisterValues[VentilationPauseRegister] == 1

	setpointTemperature := new(int)
	*setpointTemperature = int(client4RegisterValues[supplyTemperatureRegister])

	settings := Settings{FanSpeed: fanSpeed,
		DesiredRoomTemperature:      desiredRoomTemperature,
		DesiredDHWTemperature:       desiredDHWTemperature,
		DHWProductionPaused:         dhwPaused,
		DHWProductionPauseDuration:  dhwPauseDuration,
		CentralHeatingPaused:        centralHeatingPaused,
		CentralHeatingPauseDuration: centralHeatingPauseDuration,
		VentilationMode:             ventilationMode,
		VentilationOnPause:          ventilationPause,
		SetpointSupplyTemperature:   setpointTemperature}

	settingsStr := spew.Sprintf("%+v", settings)
	log.Printf("Settings: %+v\n", settingsStr)
	return settings
}

// SendSettings of Nilan
func SendSettings(settings Settings) {
	settingsStr := spew.Sprintf("%+v", settings)
	log.Printf("Sending new settings to Nialn (<nil> values will be ignored): %+v\n", settingsStr)
	client1RegisterValues := make(map[Register]uint16)
	client4RegisterValues := make(map[Register]uint16)

	if settings.FanSpeed != nil {
		fanSpeed := uint16(*settings.FanSpeed)
		client1RegisterValues[FanSpeedRegister] = fanSpeed
	}

	if settings.DesiredRoomTemperature != nil {
		desiredRoomTemperature := uint16(*settings.DesiredRoomTemperature)
		client1RegisterValues[DesiredRoomTemperatureRegister] = desiredRoomTemperature
	}

	if settings.DesiredDHWTemperature != nil {
		desiredDHWTemperature := uint16(*settings.DesiredDHWTemperature)
		client1RegisterValues[DHWSetPointRegister] = desiredDHWTemperature
	}

	if settings.DHWProductionPaused != nil {
		if *settings.DHWProductionPaused {
			client1RegisterValues[DHWPauseRegister] = uint16(1)
		} else {
			client1RegisterValues[DHWPauseRegister] = uint16(0)
		}
	}

	if settings.DHWProductionPauseDuration != nil {
		pauseDuration := uint16(*settings.DHWProductionPauseDuration)
		client1RegisterValues[DHWPauseDurationRegister] = pauseDuration
	}

	if settings.CentralHeatingPaused != nil {
		if *settings.CentralHeatingPaused {
			client4RegisterValues[CentralHeatingPauseRegister] = uint16(1)
		} else {
			client4RegisterValues[CentralHeatingPauseRegister] = uint16(0)
		}
	}

	if settings.CentralHeatingPauseDuration != nil {
		pauseDuration := uint16(*settings.CentralHeatingPauseDuration)
		client4RegisterValues[CentralHeatingPauseDurationRegister] = pauseDuration
	}

	if settings.VentilationMode != nil {
		ventilationMode := *settings.VentilationMode
		if ventilationMode != 0 && ventilationMode != 1 && ventilationMode != 2 {
			panic("Unsupported VentilationMode value")
			// TODO: Think of validation pattern
		}
		ventilationModeVal := uint16(ventilationMode)
		client1RegisterValues[VentilationModeRegister] = ventilationModeVal
	}

	if settings.VentilationOnPause != nil {
		if *settings.VentilationOnPause {
			client1RegisterValues[VentilationPauseRegister] = uint16(1)
		} else {
			client1RegisterValues[VentilationPauseRegister] = uint16(0)
		}
	}

	if settings.SetpointSupplyTemperature != nil {
		setpointTempeature := uint16(*settings.SetpointSupplyTemperature)
		client4RegisterValues[SetpointSupplyTemperatureRegisterAIR9] = setpointTempeature
		client4RegisterValues[SetpointSupplyTemperatureRegisterGEO] = setpointTempeature
	}

	setRegisterValues(1, client1RegisterValues)
	setRegisterValues(4, client4RegisterValues)
}

func roomTemperatureRegister() Register {
	if fetchValue(1, MasterTemperatureSensorSettingRegister) == 0 {
		return T3ExtractAirTemperatureRegister
	} else {
		return TextRoomTemperatureRegister
	}
}

func t18ReadingRegister() Register {
	switch {
	case fetchValue(4, DeviceTypeGEOReigister) == 8:
		return T18ReadingRegisterGEO
	case fetchValue(4, DeviceTypeAIR9Register) == 9:
		return T18ReadingRegisterAIR9
	default:
		panic("Cannot determine device type")
	}
}

// FetchReadings of Nilan sensors
func FetchReadings() Readings {
	roomTemperatureRegister := roomTemperatureRegister()
	t18Register := t18ReadingRegister()

	client1Registers := []Register{roomTemperatureRegister,
		OutdoorTemperatureRegister,
		AverageHumidityRegister,
		ActualHumidityRegister,
		DHWTopTankTemperatureRegister,
		DHWBottomTankTemperatureRegister}

	client4Registers := []Register{t18Register}

	client1ReadingsRaw := fetchRegisterValues(1, client1Registers)
	client4ReadingsRaw := fetchRegisterValues(4, client4Registers)

	roomTemperature := int(client1ReadingsRaw[roomTemperatureRegister])
	outdoorTemperature := int(client1ReadingsRaw[OutdoorTemperatureRegister])
	averageHumidity := int(client1ReadingsRaw[AverageHumidityRegister])
	actualHumidity := int(client1ReadingsRaw[ActualHumidityRegister])
	dhwTopTemperature := int(client1ReadingsRaw[DHWTopTankTemperatureRegister])
	dhwBottomTemperature := int(client1ReadingsRaw[DHWBottomTankTemperatureRegister])
	supplyFlowTemperature := int(client4ReadingsRaw[t18Register])

	readings := Readings{
		RoomTemperature:          roomTemperature,
		OutdoorTemperature:       outdoorTemperature,
		AverageHumidity:          averageHumidity,
		ActualHumidity:           actualHumidity,
		DHWTankTopTemperature:    dhwTopTemperature,
		DHWTankBottomTemperature: dhwBottomTemperature,
		SupplyFlowTemperature:    supplyFlowTemperature}
	log.Printf("Readings: %+v\n", readings)
	return readings
}
