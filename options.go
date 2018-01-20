package emitter

import (
	"crypto/tls"
	"net/url"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// ClientOptions contains configurable options for an Client.
type ClientOptions struct {
	Servers              []*url.URL
	ClientID             string
	Username             string
	Password             string
	TLSConfig            *tls.Config
	KeepAlive            time.Duration
	PingTimeout          time.Duration
	ConnectTimeout       time.Duration
	MaxReconnectInterval time.Duration
	AutoReconnect        bool
	OnMessage            OnMessageHandler
	OnConnect            OnConnectHandler
	OnConnectionLost     OnConnectionLostHandler
	OnKeyGen             OnKeyGenHandler
	OnPresence           OnPresenceHandler
}

// NewClientOptions will create a new ClientClientOptions type with some default values.
func NewClientOptions() *ClientOptions {
	id, err := uuid.NewV1()
	if err != nil {
		panic(err)
	}

	// Create new client options with defaults
	o := &ClientOptions{
		Servers:              nil,
		ClientID:             id.String(),
		Username:             "",
		Password:             "",
		TLSConfig:            &tls.Config{},
		KeepAlive:            30 * time.Second,
		PingTimeout:          10 * time.Second,
		ConnectTimeout:       30 * time.Second,
		MaxReconnectInterval: 10 * time.Minute,
		AutoReconnect:        true,
		OnConnect:            nil,
		OnConnectionLost:     defaultConnectionLostHandler,
	}
	return o
}

// AddBroker adds a broker URI to the list of brokers to be used. The format should be
// scheme://host:port
// Where "scheme" is one of "tcp", "ssl", or "ws", "host" is the ip-address (or hostname)
// and "port" is the port on which the broker is accepting connections.
func (o *ClientOptions) AddBroker(server string) *ClientOptions {
	brokerURI, _ := url.Parse(server)
	o.Servers = append(o.Servers, brokerURI)
	return o
}

// SetClientID will set the client id to be used by this client when
// connecting to the MQTT broker. According to the MQTT v3.1 specification,
// a client id mus be no longer than 23 characters.
func (o *ClientOptions) SetClientID(id string) *ClientOptions {
	o.ClientID = id
	return o
}

// SetUsername will set the username to be used by this client when connecting
// to the MQTT broker. Note: without the use of SSL/TLS, this information will
// be sent in plaintext accross the wire.
func (o *ClientOptions) SetUsername(u string) *ClientOptions {
	o.Username = u
	return o
}

// SetPassword will set the password to be used by this client when connecting
// to the MQTT broker. Note: without the use of SSL/TLS, this information will
// be sent in plaintext accross the wire.
func (o *ClientOptions) SetPassword(p string) *ClientOptions {
	o.Password = p
	return o
}

// SetTLSConfig will set an SSL/TLS configuration to be used when connecting
// to an MQTT broker. Please read the official Go documentation for more
// information.
func (o *ClientOptions) SetTLSConfig(t *tls.Config) *ClientOptions {
	o.TLSConfig = t
	return o
}

// SetKeepAlive will set the amount of time (in seconds) that the client
// should wait before sending a PING request to the broker. This will
// allow the client to know that a connection has not been lost with the
// server.
func (o *ClientOptions) SetKeepAlive(k time.Duration) *ClientOptions {
	o.KeepAlive = k
	return o
}

// SetPingTimeout will set the amount of time (in seconds) that the client
// will wait after sending a PING request to the broker, before deciding
// that the connection has been lost. Default is 10 seconds.
func (o *ClientOptions) SetPingTimeout(k time.Duration) *ClientOptions {
	o.PingTimeout = k
	return o
}

// SetOnMessageHandler sets the MessageHandler that will be called when a message
// is received that does not match any known subscriptions.
func (o *ClientOptions) SetOnMessageHandler(defaultHandler OnMessageHandler) *ClientOptions {
	o.OnMessage = defaultHandler
	return o
}

// SetOnConnectHandler sets the function to be called when the client is connected. Both
// at initial connection time and upon automatic reconnect.
func (o *ClientOptions) SetOnConnectHandler(onConn OnConnectHandler) *ClientOptions {
	o.OnConnect = onConn
	return o
}

// SetOnConnectionLostHandler will set the OnConnectionLost callback to be executed
// in the case where the client unexpectedly loses connection with the MQTT broker.
func (o *ClientOptions) SetOnConnectionLostHandler(onLost OnConnectionLostHandler) *ClientOptions {
	o.OnConnectionLost = onLost
	return o
}

// SetConnectTimeout limits how long the client will wait when trying to open a connection
// to an MQTT server before timeing out and erroring the attempt. A duration of 0 never times out.
// Default 30 seconds. Currently only operational on TCP/TLS connections.
func (o *ClientOptions) SetConnectTimeout(t time.Duration) *ClientOptions {
	o.ConnectTimeout = t
	return o
}

// SetMaxReconnectInterval sets the maximum time that will be waited between reconnection attempts
// when connection is lost
func (o *ClientOptions) SetMaxReconnectInterval(t time.Duration) *ClientOptions {
	o.MaxReconnectInterval = t
	return o
}

// SetAutoReconnect sets whether the automatic reconnection logic should be used
// when the connection is lost, even if disabled the ConnectionLostHandler is still
// called
func (o *ClientOptions) SetAutoReconnect(a bool) *ClientOptions {
	o.AutoReconnect = a
	return o
}

// SetOnPresenceHandler sets the OnPresenceHandler that will be called when a presence event is received.
func (o *ClientOptions) SetOnPresenceHandler(handler OnPresenceHandler) *ClientOptions {
	o.OnPresence = handler
	return o
}

// SetOnKeyGenHandler sets the OnKeyGenHandler that will be called when a key generation response is received.
func (o *ClientOptions) SetOnKeyGenHandler(handler OnKeyGenHandler) *ClientOptions {
	o.OnKeyGen = handler
	return o
}

// Option represents a key/value pair that can be supplied to the publish/subscribe or unsubscribe
// methods and provide ways to configure the operation.
type Option struct {
	Key   string
	Value string
}

// Makes a topic name from the key/channel pair
func formatTopic(key string, channel string, options []Option) string {
	// Clean the key
	key = strings.TrimPrefix(key, "/")
	key = strings.TrimSuffix(key, "/")

	// Clean the channel name
	channel = strings.TrimPrefix(channel, "/")
	channel = strings.TrimSuffix(channel, "/")

	// Add the options
	opts := ""
	if options != nil && len(options) > 0 {
		opts += "?"
		for i, option := range options {
			opts += option.Key + "=" + option.Value
			if i+1 < len(options) {
				opts += "&"
			}
		}
	}

	// Concatenate
	return key + "/" + channel + "/" + opts
}
