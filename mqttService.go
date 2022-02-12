package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var topics = [...]string{"sensor/water"}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s.\n", msg.Payload(), msg.Topic())
	sendSms()
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected.")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v.\n", err)
}

func initMqtt() {
	var broker = os.Getenv("MOSQUITTO_BROKER")
	var clientId = os.Getenv("MOSQUITTO_CLIENT_ID")
	var username = os.Getenv("MOSQUITTO_USERNAME")
	var password = os.Getenv("MOSQUITTO_PASSWORD")
	port, err := strconv.Atoi(os.Getenv("MOSQUITTO_PORT"))
	if err != nil {
		panic("Invalid Mosquitto port.")
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%d", broker, port))
	tlsConfig := newTlsConfig()
	opts.SetTLSConfig(tlsConfig)
	opts.SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)
	publish(client)

	// client.Disconnect(1000 * 1000 * 1000)
}

func publish(client mqtt.Client) {
	num := 1
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d.", i)
		token := client.Publish("topic/fillmypihole", 0, false, text)
		token.Wait()
		// time.Sleep(time.Second)
	}
}

func sub(client mqtt.Client) {
	for _, topic := range topics {
		token := client.Subscribe(topic, 1, nil)
		token.Wait()
		fmt.Printf("Subscribed to topic: %s\n", topic)
	}
}

func newTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		panic(err.Error())
	}

	certpool.AppendCertsFromPEM(ca)
	return &tls.Config{RootCAs: certpool}
}
