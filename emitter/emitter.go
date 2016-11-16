package emitter

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Message defines the externals that a message implementation must support
// these are received messages that are passed to the callbacks, not internal
// messages
type Message interface {
	Topic() string
	Payload() []byte
}

//Token defines the interface for the tokens used to indicate when
//actions have completed.
type Token interface {
	Wait() bool
	WaitTimeout(time.Duration) bool
	Error() error
}

// Emitter defines the externals that a message implementation must support
// these are received messages that are passed to the callbacks, not internal
// messages
type Emitter interface {
	IsConnected() bool
	Connect() Token
	Disconnect(uint)
	Publish(string, interface{}) Token
	Subscribe(string) Token
	Unsubscribe(string) Token
}

type emitter struct {
	conn mqtt.Client
}

// NewClient will create an MQTT v3.1.1 client with all of the options specified
// in the provided ClientOptions. The client must have the Connect method called
// on it before it may be used. This is to make sure resources (such as a net
// connection) are created before the application is actually ready.
func NewClient(o *ClientOptions) Emitter {

	// Create an emitter client
	c := &emitter{}

	// Copy options to mqtt.ClientOptions
	mqttOptions := mqtt.NewClientOptions()
	mqttOptions.Servers = o.Servers
	mqttOptions.ClientID = o.ClientID
	mqttOptions.Username = o.Username
	mqttOptions.Password = o.Password
	mqttOptions.TLSConfig = o.TLSConfig
	mqttOptions.KeepAlive = o.KeepAlive
	mqttOptions.PingTimeout = o.PingTimeout
	mqttOptions.ConnectTimeout = o.ConnectTimeout
	mqttOptions.MaxReconnectInterval = o.MaxReconnectInterval
	mqttOptions.AutoReconnect = o.AutoReconnect

	// Set the mqtt handler to call out into our emitter connection handler
	mqttOptions.SetOnConnectHandler(func(_ mqtt.Client) {
		if o.OnConnect != nil {
			o.OnConnect(c)
		}
	})

	// Set the mqtt handler to call out into our emitter connection lost handler
	mqttOptions.SetConnectionLostHandler(func(_ mqtt.Client, e error) {
		if o.OnConnectionLost != nil {
			o.OnConnectionLost(c, e)
		}
	})

	// Set the mqtt handler to call out into our emitter connection lost handler
	mqttOptions.SetDefaultPublishHandler(func(_ mqtt.Client, m mqtt.Message) {
		if o.OnMessage != nil {
			o.OnMessage(c, m)
		}
	})

	// Create the underlying MQTT client and set the options
	c.conn = mqtt.NewClient(mqttOptions)
	return c
}

// IsConnected returns a bool signifying whether the client is connected or not.
func (c *emitter) IsConnected() bool {
	return c.conn.IsConnected()
}

// Connect will create a connection to the message broker
// If clean session is false, then a slice will
// be returned containing Receipts for all messages
// that were in-flight at the last disconnect.
// If clean session is true, then any existing client
// state will be removed.
func (c *emitter) Connect() Token {
	return c.conn.Connect()
}

// Disconnect will end the connection with the server, but not before waiting
// the specified number of milliseconds to wait for existing work to be
// completed.
func (c *emitter) Disconnect(waitTime uint) {
	c.conn.Disconnect(waitTime)
}

// Publish will publish a message with the specified QoS and content
// to the specified topic.
// Returns a token to track delivery of the message to the broker
func (c *emitter) Publish(topic string, payload interface{}) Token {
	return c.conn.Publish(topic, 0, false, payload)
}

// Subscribe starts a new subscription. Provide a MessageHandler to be executed when
// a message is published on the topic provided.
func (c *emitter) Subscribe(topic string) Token {
	return c.conn.Subscribe(topic, 0, nil)
}

// Unsubscribe will end the subscription from each of the topics provided.
// Messages published to those topics from other clients will no longer be
// received.
func (c *emitter) Unsubscribe(topic string) Token {
	return c.conn.Unsubscribe(topic)
}
