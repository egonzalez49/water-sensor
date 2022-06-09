##
## Build
##

# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.17.7-alpine AS build

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy over the project files
COPY ./ ./

# Compile the project
RUN go build -o /water-sensor-service ./cmd/water-sensor-service

CMD ["/water-sensor-service"]

# ##
# ## Deploy
# ##

# FROM gcr.io/distroless/base-debian10

# WORKDIR /app

# COPY --from=build /app/water-sensor-service ./

# ENTRYPOINT ["./water-sensor-service"]