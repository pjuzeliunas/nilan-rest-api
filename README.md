# nilan-rest-api

Nilan API server is a REST interface for controlling Nilan heatpump. Server communicates with Nilan CTS700 using
[Modbus TCP protocol](https://www.nilan.dk/Admin/Public/DWSDownload.aspx?File=%2fFiles%2fFiler%2fDownload%2fDanish%2fDokumentation%2fSoftware+vejledninger%2fModbus%2fCTS700_Modbus_protokol.pdf).
It exposes heatpump settings and readings in JSON format.

Server is written in Go.


## Usage

REST API server is running on port 8080 and supports the following HTTP methods:

#### GET `/readings`

Returns basic readings of Nilan heatpump in JSON format.

```
{
    "RoomTemperature": 242,
    "OutdoorTemperature": 83,
    "AverageHumidity": 38,
    "ActualHumidity": 37,
    "DHWTankTopTemperature": 553,
    "DHWTankBottomTemperature": 511,
    "SupplyFlowTemperature": 346
}
```

Temperature is in ℃ times 10.

#### GET `/settings`

Returns basic settings of Nilan heatpump in JSON format.

```
{
    "FanSpeed": 102,
    "DesiredRoomTemperature": 230,
    "DesiredDHWTemperature": 530,
    "DHWProductionPaused": false,
    "DHWProductionPauseDuration": 3,
    "CentralHeatingPaused": false,
    "CentralHeatingPauseDuration": 1,
    "VentilationMode": 0,
    "VentilationOnPause": false,
    "SetpointSupplyTemperature": 350
}
```

Fan speed can be 101 (level 1), 102 (level 2), 103 (level 3) or 104 (level 4).

Temperature is in ℃ times 10.

Ventilation mode can be 0 (auto), 1 (cooling) or 2 (heating).

#### PUT `/settings`

Sets settings of Nilan heatpump. Body must be in the same JSON format as one from GET request.
Nil values or absence of them tells the server to keep the existing setting unchanged,
thus making it possible to send a single setting change in a short form. For example, to change just
the fan speed, the following request body is valid:
```
{ "FanSpeed": 103 }
```


## Hardware setup

Tested on
[Nilan Compact P AIR 9](https://www.nilan.dk/da-dk/forside/loesninger/boligloesninger/kompaktloesning/compact-p-air-9) heatpump 
and [Raspbery Pi 3](https://static.raspberrypi.org/files/product-briefs/Raspberry-Pi-Model-Bplus-Product-Brief.pdf) host.

Connect Nilan to the host computer using Ethernet cable. Ethernet cable should hanging inside the heatpump.

## Software setup

### Prerequisites

1. Host computer and Nilan heatpump must be on the same network. By default, heatpump has this IP address: 192.168.5.107(:502).
Make sure that host computer and heatpump are running on the same network (subnet) by adjusting ethernet port IP address.
2. Install either golang or docker.

### Set up server using Docker

Build and run Docker container as follows:
```
docker build -t nilan .
docker run -e NILAN_ADDRESS=<IP and port of Nilan> -it --rm -p 8080:8080 nilan
```

For more sophisticated setup refer to Docker documentation.

### Set up server using Go

Compile and run server as follows:
```
go get -d -v ./...
go build -o nilanapp app/app.go
NILAN_ADDRESS=<IP and port of Nilan> ./nilanapp
```

## Disclaimer

This initial version of API server is developed by home automation enthusiast (outside Nilan company) and is based
on open Nilan CTS700 Modbus protocol.

Nilan is a registered trademark and belongs to [Nilan A/S](https://www.nilan.dk/Default.aspx).
