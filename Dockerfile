# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.13

# Add Maintainer Info
LABEL maintainer="Povilas Juzeliunas <pjuzeliunas@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
#COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
#RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go get -d -v ./...
RUN go build -o nilanapp app/app.go

# Expose port 8080 to the outside world
# EXPOSE 8080

# Command to run the executable
CMD ["./nilanapp"]