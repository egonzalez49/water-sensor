# water-sensor [![Tests](https://github.com/egonzalez49/water-sensor/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/egonzalez49/water-sensor/actions/workflows/test.yml)
Containerized service to receive water sensor events and send text message alerts.  
  
Runs an MQTT broker listening for emitted events from an RTL-SDR device using the open-source [rtl_433](https://github.com/merbanan/rtl_433) library,  
as well as a Go subscriber client for processing the events and sending texts via Twilio.  
  
Redis is used as a caching mechanism to prevent text message spam in case of repeat signals within a short time span.


## Getting Started

To generate a Mosquitto user, run:
```bash
mosquitto_passwd -c ./mosquitto/config/pwfile user
```

To delete a Mosquitto user, run:
```bash
mosquitto_passwd -D ./mosquitto/config/pwfile user
```

Generate the necessary self-signed certificates for Mosquitto to use TLS connections.
```bash
./scripts/gen-keys.sh
```

Move the `ca.crt`, `server.crt`, and `server.key` to `./mosquitto/config/certs`.  
Create a copy of the `ca.crt` file into the root of the project.

Make sure to fill out the required environment variables inside a `.env` file at the root of the project.

Finally, run the following to boot up the containers:
```bash
docker-compose up --build
```
