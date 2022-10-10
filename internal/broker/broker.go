package broker

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Broker struct {
	client  mqtt.Client
	options *mqtt.ClientOptions
}

func New(clientId, host string, port int) Broker {
	options := mqtt.NewClientOptions()

	address := fmt.Sprintf("ssl://%s:%d", host, port)
	options.AddBroker(address)
	options.SetClientID(clientId)

	return Broker{
		options: options,
	}
}

// Connect establishes a connection to the client broker.
// Returns an error if a connection cannot be established.
func (b *Broker) Connect(username, password string) error {
	tlsConfig, err := newTlsConfig()
	if err != nil {
		return err
	}

	b.options.SetTLSConfig(tlsConfig)
	b.options.SetUsername(username)
	b.options.SetPassword(password)

	b.client = mqtt.NewClient(b.options)
	token := b.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Disconnect waits the specified number of milliseconds before
// closing the connection to the broker to allow existing work to finish.
func (b *Broker) Disconnect(waitTime uint) {
	b.client.Disconnect(waitTime)
}

// Returns a map of various broker connection properties.
func (b *Broker) Properties() map[string]string {
	options := b.client.OptionsReader()

	url := options.Servers()[0]
	host, port, _ := net.SplitHostPort(url.Host)

	return map[string]string{
		"host":     host,
		"port":     port,
		"clientId": options.ClientID(),
	}
}

// Configures a certificate authority cert
// to be used for establishing TLS connection
// with the MQTT broker.
func newTlsConfig() (*tls.Config, error) {
	certpool := x509.NewCertPool()
	ca, err := os.ReadFile("ca.crt")
	if err != nil {
		return nil, err
	}

	ok := certpool.AppendCertsFromPEM(ca)
	if !ok {
		return nil, errors.New("unable to parse certificate authority cert")
	}
	cfg := &tls.Config{RootCAs: certpool}
	return cfg, nil
}

// Set a handler for whenever the connection with the broker is established.
func (b *Broker) SetOnConnectHandler(onConn mqtt.OnConnectHandler) {
	b.options.OnConnect = onConn
}

// Set a handler for whenever the connection with the broker is lost.
func (b *Broker) SetOnConnectionLostHandler(onLost mqtt.ConnectionLostHandler) {
	b.options.OnConnectionLost = onLost
}

// Set a handler for whenever a message is received on an unsubscribed topic
// or when a subscribed topic has no specified handler.
func (b *Broker) SetDefaultMessageHandler(defaultHandler mqtt.MessageHandler) {
	b.options.SetDefaultPublishHandler(defaultHandler)
}

// Subscribes to the specified topics, assigning each the specified quality of service level.
// The specified callback will be called whenever a message from the topics is received.
// If no callback is specified, the default message handler will be called instead.
func (b *Broker) Subscribe(filters map[string]byte, callback mqtt.MessageHandler) error {
	token := b.client.SubscribeMultiple(filters, callback)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
