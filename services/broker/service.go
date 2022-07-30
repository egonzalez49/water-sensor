package broker

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/egonzalez49/water-sensor/config"
	"github.com/egonzalez49/water-sensor/logging"
)

type Broker struct {
	client mqtt.Client
	opts   *mqtt.ClientOptions
	Logger *logging.Logger
}

func NewBroker(cfg *config.Config, logger *logging.Logger) (*Broker, error) {
	opts, err := initOpts(cfg)
	if err != nil {
		return nil, err
	}

	return &Broker{opts: opts, Logger: logger}, nil
}

func initOpts(cfg *config.Config) (*mqtt.ClientOptions, error) {
	host := cfg.Mqtt.Host
	port := cfg.Mqtt.Port

	clientId := cfg.Mqtt.ClientId
	username := cfg.Mqtt.Username
	password := cfg.Mqtt.Password

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%s", host, port))

	tlsConfig, err := newTlsConfig()
	if err != nil {
		return nil, err
	}
	opts.SetTLSConfig(tlsConfig)

	opts.SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(password)

	return opts, nil
}

func newTlsConfig() (*tls.Config, error) {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		return nil, err
	}

	certpool.AppendCertsFromPEM(ca)
	cfg := &tls.Config{RootCAs: certpool}
	return cfg, nil
}

func (b *Broker) SetConnectionHandlers(onConn mqtt.OnConnectHandler, onLost mqtt.ConnectionLostHandler) {
	b.opts.OnConnect = onConn
	b.opts.OnConnectionLost = onLost
}

func (b *Broker) SetDefaultPublishHandler(defaultHandler mqtt.MessageHandler) {
	b.opts.SetDefaultPublishHandler(defaultHandler)
}

func (b *Broker) Connect() error {
	client := mqtt.NewClient(b.opts)
	b.client = client

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (b *Broker) Subscribe(filters map[string]byte, callback mqtt.MessageHandler) {
	if token := b.client.SubscribeMultiple(filters, callback); token.Wait() && token.Error() != nil {
		b.Logger.Errorf("error when subscribing to topics: %v\n", token.Error())
	} else {
		b.Logger.Info("successfully subscribed to topics.")
	}
}
