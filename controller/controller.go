package controller

import (
	"encoding/binary"
	"log"
	"time"

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
	} else {
		panic("Cannot read register value")
	}
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
	// FanSpeedRegister is a register with desired FanSpeed value
	FanSpeedRegister Register = 20148
	// DesiredRoomTemperatureRegister is a register with desired room temperature in C times 10.
	// Example: 23.5 C is stored as 235.
	DesiredRoomTemperatureRegister Register = 20260
	// MasterTemperatureSensorSettingRegister holds either 0 (read temperature from T3) or 1
	// (read temperature from Text)
	MasterTemperatureSensorSettingRegister Register = 20263
	// T3ExtractAirTemperatureRegister holds room temperature value when
	// MasterTemperatureSensorSettingRegister is 0
	T3ExtractAirTemperatureRegister Register = 20286
	// TextRoomTemperatureRegister holds room temperature value when
	// MasterTemperatureSensorSettingRegister is 1
	TextRoomTemperatureRegister Register = 20280
	// OutdoorTemperatureRegister holds outdoor temperature
	OutdoorTemperatureRegister Register = 20282
)

// FetchSettings of Nilan
func FetchSettings() Settings {
	registers := []Register{FanSpeedRegister, DesiredRoomTemperatureRegister}
	registerValues := fetchRegisterValues(1, registers)

	fanSpeed := FanSpeed(registerValues[FanSpeedRegister])
	desiredRoomTemperature := int(registerValues[DesiredRoomTemperatureRegister])

	settings := Settings{FanSpeed: fanSpeed, DesiredRoomTemperature: desiredRoomTemperature}
	log.Printf("Settings: %+v\n", settings)
	return settings
}

// SendSettings of Nilan
func SendSettings(settings Settings) {
	log.Printf("New settings: %+v\n", settings)
	registerValues := make(map[Register]uint16)

	fanSpeed := uint16(settings.FanSpeed)
	desiredRoomTemperature := uint16(settings.DesiredRoomTemperature)

	registerValues[FanSpeedRegister] = fanSpeed
	registerValues[DesiredRoomTemperatureRegister] = desiredRoomTemperature

	setRegisterValues(1, registerValues)
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

	registers := []Register{roomTemperatureRegister, OutdoorTemperatureRegister}
	readingsRaw := fetchRegisterValues(1, registers)

	roomTemperature := int(readingsRaw[roomTemperatureRegister])
	outdoorTemperature := int(readingsRaw[OutdoorTemperatureRegister])

	readings := Readings{RoomTemperature: roomTemperature, OutdoorTemperature: outdoorTemperature}
	log.Printf("Readings: %+v\n", readings)
	return readings
}
