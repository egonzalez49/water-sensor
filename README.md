# water-sensor

## Getting Started

To generate a Mosquitto user, run:
```bash
mosquitto_passwd -c passwordfile user
```

To add a password for the user, run:
```bash
mosquitto_passwd -b passwordfile user password
```

To delete a Mosquitto user, run:
```bash
mosquitto_passwd -D passwordfile user
```