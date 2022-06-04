package subscriber

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/egonzalez49/water-sensor/notifier"
	"github.com/go-redis/redis/v8"
)

var topics = [...]string{"sensor/water"}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s.\n", msg.Payload(), msg.Topic())

	type data struct {
		time  string
		model string
		id    string
		event string
		code  string
		mic   string
	}

	var sensorData data
	parseMessage(msg.Payload(), sensorData)

	_, err := rdb.Get(ctx, sensorData.id).Result()
	if err == redis.Nil {
		// Key does not exist in cache
		rdb.Set(ctx, sensorData.id, true, 5*time.Minute)
		notifier.Notify()
	} else if err != nil {
		panic(err)
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected.")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v.\n", err)
}

func Subscribe() {
	initRedis()

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

func parseMessage(bytes []byte, addr interface{}) {
	if err := json.Unmarshal(bytes, &addr); err != nil {
		fmt.Println(err)
	}
}
