package emitter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
	Publish(string, string, interface{}, ...Option) Token
	PublishWithTTL(string, string, interface{}, int) Token
	Subscribe(string, string, ...Option) Token
	SubscribeWithHistory(string, string, int) Token
	Unsubscribe(string, string, ...Option) Token
	Presence(*PresenceRequest) Token
	GenerateKey(*KeyGenRequest) Token
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

	// If there's no brokers configured, configure the default one
	if o.Servers == nil {
		o.AddBroker("tcp://api.emitter.io:8080")
	}

	// Copy options to mqtt.ClientOptions
	mqttOptions := mqtt.NewClientOptions()
	mqttOptions.Servers = o.Servers
	mqttOptions.ClientID = o.ClientID
	mqttOptions.Username = o.Username
	mqttOptions.Password = o.Password
	mqttOptions.TLSConfig = o.TLSConfig
	mqttOptions.KeepAlive = int64(o.KeepAlive)
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

		if strings.HasPrefix(m.Topic(), "emitter/keygen") {
			// Invoke the keygen event
			if o.OnKeyGen != nil {
				r := &KeyGenResponse{}
				json.Unmarshal(m.Payload(), r)
				o.OnKeyGen(c, *r)
			}

		} else if strings.HasPrefix(m.Topic(), "emitter/presence") {
			// Invoke the presence event
			if o.OnPresence != nil {
				r := &PresenceEvent{}
				json.Unmarshal(m.Payload(), r)
				o.OnPresence(c, *r)
			}

		} else {
			// Invoke message handler
			if o.OnMessage != nil {
				o.OnMessage(c, m)
			}
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
func (c *emitter) Publish(key string, channel string, payload interface{}, options ...Option) Token {
	return c.conn.Publish(formatTopic(key, channel, options), 0, false, payload)
}

// PublishWithTTL publishes a message with a specified Time-To-Live option
func (c *emitter) PublishWithTTL(key string, channel string, payload interface{}, ttl int) Token {
	opt1 := Option{Key: "ttl", Value: strconv.Itoa(ttl)}
	return c.Publish(key, channel, payload, opt1)
}

// Subscribe starts a new subscription. Provide a MessageHandler to be executed when
// a message is published on the topic provided.
func (c *emitter) Subscribe(key string, channel string, options ...Option) Token {
	return c.conn.Subscribe(formatTopic(key, channel, options), 0, nil)
}

// SubscribeWithHistory performs a subscribe with an option to retrieve the specified number
// of messages that were already published in the channel.
func (c *emitter) SubscribeWithHistory(key string, channel string, last int) Token {
	opt1 := Option{Key: "last", Value: strconv.Itoa(last)}
	return c.Subscribe(key, channel, opt1)
}

// Unsubscribe will end the subscription from each of the topics provided.
// Messages published to those topics from other clients will no longer be
// received.
func (c *emitter) Unsubscribe(key string, channel string, options ...Option) Token {
	return c.conn.Unsubscribe(formatTopic(key, channel, options))
}

// GenerateKey sends a key generation request to the broker
func (c *emitter) GenerateKey(r *KeyGenRequest) Token {
	serialized, err := json.Marshal(r)
	if err != nil {
		fmt.Println("Unable to serialize keygen request.")
	}

	return c.conn.Publish("emitter/keygen/", 0, false, serialized)
}

// GenerateKey sends a key generation request to the broker
func (c *emitter) Presence(r *PresenceRequest) Token {
	serialized, err := json.Marshal(r)
	if err != nil {
		fmt.Println("Unable to serialize presence request.")
	}

	return c.conn.Publish("emitter/presence/", 0, false, serialized)
}

// OnMessageHandler is a callback type which can be set to be
// executed upon the arrival of messages published to topics
// to which the client is subscribed.
type OnMessageHandler func(Emitter, Message)

// OnKeyGenHandler is a callback type which can be set to be executed upon
// the arrival of key generation responses.
type OnKeyGenHandler func(Emitter, KeyGenResponse)

// OnPresenceHandler is a callback type which can be set to be executed upon
// the arrival of presence events.
type OnPresenceHandler func(Emitter, PresenceEvent)

// OnConnectionLostHandler is a callback type which can be set to be
// executed upon an unintended disconnection from the MQTT broker.
// Disconnects caused by calling Disconnect or ForceDisconnect will
// not cause an OnConnectionLost callback to execute.
type OnConnectionLostHandler func(Emitter, error)

// OnConnectHandler is a callback that is called when the client
// state changes from unconnected/disconnected to connected. Both
// at initial connection and on reconnection
type OnConnectHandler func(Emitter)

// Default connection lost handler, simply prints out the log
func defaultConnectionLostHandler(client Emitter, reason error) {
	fmt.Println("Connection lost:", reason.Error())
}
