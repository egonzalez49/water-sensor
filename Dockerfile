# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.17.7-alpine

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy over the project files
COPY *.go ./

# Compile the project
RUN go build -o /water-sensor-client

CMD [ "/water-sensor-client" ]