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
	// WaterAfterHeaterTemperatureRegister is ID of register holding T9 water after heater temperature
	WaterAfterHeaterTemperatureRegister Register = 20298
	// DHWTopTankTemperatureRegister is ID of register holding T21 top DHW tank temperature
	DHWTopTankTemperatureRegister Register = 20580
	// DHWBottomTankTemperatureRegister is ID of register holding T21 bottom DHW tank temperature
	DHWBottomTankTemperatureRegister Register = 20582
)

// FetchSettings of Nilan
func FetchSettings() Settings {
	registers := []Register{FanSpeedRegister, DesiredRoomTemperatureRegister}
	registerValues := fetchRegisterValues(1, registers)

	fanSpeed := new(FanSpeed)
	*fanSpeed = FanSpeed(registerValues[FanSpeedRegister])

	desiredRoomTemperature := new(int)
	*desiredRoomTemperature = int(registerValues[DesiredRoomTemperatureRegister])

	settings := Settings{FanSpeed: fanSpeed, DesiredRoomTemperature: desiredRoomTemperature}
	log.Printf("Settings: %+v\n", settings)
	return settings
}

// SendSettings of Nilan
func SendSettings(settings Settings) {
	settingsStr := spew.Sprintf("%+v", settings)
	log.Printf("Sending new settings to Nialn (<nil> values will be ignored): %+v\n", settingsStr)
	registerValues := make(map[Register]uint16)

	if settings.FanSpeed != nil {
		fanSpeed := new(uint16)
		*fanSpeed = uint16(*settings.FanSpeed)
		registerValues[FanSpeedRegister] = *fanSpeed
	}

	if settings.DesiredRoomTemperature != nil {
		desiredRoomTemperature := new(uint16)
		*desiredRoomTemperature = uint16(*settings.DesiredRoomTemperature)
		registerValues[DesiredRoomTemperatureRegister] = *desiredRoomTemperature
	}

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

	registers := []Register{roomTemperatureRegister,
		OutdoorTemperatureRegister,
		AverageHumidityRegister,
		ActualHumidityRegister,
		WaterAfterHeaterTemperatureRegister,
		DHWTopTankTemperatureRegister,
		DHWBottomTankTemperatureRegister}
	readingsRaw := fetchRegisterValues(1, registers)

	roomTemperature := int(readingsRaw[roomTemperatureRegister])
	outdoorTemperature := int(readingsRaw[OutdoorTemperatureRegister])
	averageHumidity := int(readingsRaw[AverageHumidityRegister])
	actualHumidity := int(readingsRaw[ActualHumidityRegister])
	waterAfterHeaterTemperature := int(readingsRaw[WaterAfterHeaterTemperatureRegister])
	dhwTopTemperature := int(readingsRaw[DHWTopTankTemperatureRegister])
	dhwBottomTemperature := int(readingsRaw[DHWBottomTankTemperatureRegister])

	readings := Readings{
		RoomTemperature:             roomTemperature,
		OutdoorTemperature:          outdoorTemperature,
		AverageHumidity:             averageHumidity,
		ActualHumidity:              actualHumidity,
		WaterAfterHeaterTemperature: waterAfterHeaterTemperature,
		DHWTankTopTemperature:       dhwTopTemperature,
		DHWTankBottomTemperature:    dhwBottomTemperature}
	log.Printf("Readings: %+v\n", readings)
	return readings
}
