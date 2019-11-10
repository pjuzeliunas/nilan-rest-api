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

func fetchValue(register Register) uint16 {
	handler := getHandler(1)
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

	handler := getHandler(1)
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
	handler := getHandler(1)
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
	// VentilationModelRegister is ID of register holding ventilation mode value (0, 1 or 2).
	VentilationModelRegister Register = 20120
	// VentilationPauseRegister is ID of register holding ventilation pause flag
	VentilationPauseRegister Register = 20100
)

// FetchSettings of Nilan
func FetchSettings() Settings {
	client1Registers := []Register{
		FanSpeedRegister,
		DesiredRoomTemperatureRegister,
		DHWSetPointRegister,
		DHWPauseRegister,
		DHWPauseDurationRegister,
		VentilationModelRegister,
		VentilationPauseRegister}
	client4Registers := []Register{
		CentralHeatingPauseRegister,
		CentralHeatingPauseDurationRegister}

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
	*ventilationMode = int(client1RegisterValues[VentilationModelRegister])

	ventilationPause := new(bool)
	*ventilationPause = client1RegisterValues[VentilationPauseRegister] == 1

	settings := Settings{FanSpeed: fanSpeed,
		DesiredRoomTemperature:      desiredRoomTemperature,
		DesiredDHWTemperature:       desiredDHWTemperature,
		DHWProductionPaused:         dhwPaused,
		DHWProductionPauseDuration:  dhwPauseDuration,
		CentralHeatingPaused:        centralHeatingPaused,
		CentralHeatingPauseDuration: centralHeatingPauseDuration,
		VentilationMode:             ventilationMode,
		VentilationOnPause:          ventilationPause}

	log.Printf("Settings: %+v\n", settings)
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
		client4RegisterValues[VentilationModelRegister] = ventilationModeVal
	}

	if settings.VentilationOnPause != nil {
		if *settings.VentilationOnPause {
			client1RegisterValues[VentilationPauseRegister] = uint16(1)
		} else {
			client1RegisterValues[VentilationPauseRegister] = uint16(0)
		}
	}

	setRegisterValues(1, client1RegisterValues)
	setRegisterValues(4, client4RegisterValues)
}

// FetchReadings of Nilan sensors
func FetchReadings() Readings {
	var roomTemperatureRegister Register
	// Room temperature is taken from one of two sensors depending on the flag value
	masterTemperatureSensorSetting := fetchValue(MasterTemperatureSensorSettingRegister)
	if masterTemperatureSensorSetting == 0 {
		roomTemperatureRegister = T3ExtractAirTemperatureRegister
	} else {
		roomTemperatureRegister = TextRoomTemperatureRegister
	}

	registers := []Register{roomTemperatureRegister,
		OutdoorTemperatureRegister,
		AverageHumidityRegister,
		ActualHumidityRegister,
		DHWTopTankTemperatureRegister,
		DHWBottomTankTemperatureRegister}
	readingsRaw := fetchRegisterValues(1, registers)

	roomTemperature := int(readingsRaw[roomTemperatureRegister])
	outdoorTemperature := int(readingsRaw[OutdoorTemperatureRegister])
	averageHumidity := int(readingsRaw[AverageHumidityRegister])
	actualHumidity := int(readingsRaw[ActualHumidityRegister])
	dhwTopTemperature := int(readingsRaw[DHWTopTankTemperatureRegister])
	dhwBottomTemperature := int(readingsRaw[DHWBottomTankTemperatureRegister])

	readings := Readings{
		RoomTemperature:          roomTemperature,
		OutdoorTemperature:       outdoorTemperature,
		AverageHumidity:          averageHumidity,
		ActualHumidity:           actualHumidity,
		DHWTankTopTemperature:    dhwTopTemperature,
		DHWTankBottomTemperature: dhwBottomTemperature}
	log.Printf("Readings: %+v\n", readings)
	return readings
}
