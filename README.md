# nilan-rest-api

Nilan API server is a REST interface for controlling Nilan heatpump. Server communicates with Nilan CTS700 using
[Modbus TCP protocol](https://www.nilan.dk/Admin/Public/DWSDownload.aspx?File=%2fFiles%2fFiler%2fDownload%2fDanish%2fDokumentation%2fSoftware+vejledninger%2fModbus%2fCTS700_Modbus_protokol.pdf).
It exposes heatpump settings and readings in JSON format.

Server is written in Go.

## Hardware setup

Tested on
[Nilan Compact P AIR 9](https://www.nilan.dk/da-dk/forside/loesninger/boligloesninger/kompaktloesning/compact-p-air-9) heatpump 
and [Raspbery Pi 3](https://static.raspberrypi.org/files/product-briefs/Raspberry-Pi-Model-Bplus-Product-Brief.pdf) host.

Connect Nilan to the host computer using Ethernet cable. Ethernet cable should hanging inside the heatpump.

## Software setup

### Prerequisites

1. Host computer and Nilan heatpump must be on the same network. By default, heatpump has this IP address: 192.168.5.107.
Make sure that host computer and heatpump are running on the same network (subnet) by adjusting ethernet port IP address.
2. Install either golang or docker.

### Set up server using Docker

Build and run Docker container as follows:
```
docker build -t nilan .
docker run -it --rm nilan
```

### Set up server using Go

Compile and run server as follows:
```
go get -d -v ./...
go build -o nilanapp app/app.go
./nilanapp
```

