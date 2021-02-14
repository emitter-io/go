package emitter

import (
	"crypto/tls"
	"net/url"
	"strconv"
	"time"
)

// WithMatcher If "mqtt", then topic matching would follow MQTT specification.
func WithMatcher(matcher string) func(*Client) {
	return func(c *Client) {
		if "mqtt" == matcher {
			c.handlers = NewTrieMQTT()
		}
	}
}

// WithBrokers configures broker URIs to connect to. The format should be scheme://host:port
// Where "scheme" is one of "tcp", "ssl", or "ws", "host" is the ip-address (or hostname)
// and "port" is the port on which the broker is accepting connections.
func WithBrokers(brokers ...string) func(*Client) {
	return func(c *Client) {
		c.opts.Servers = []*url.URL{}
		for _, broker := range brokers {
			brokerURI, err := url.Parse(broker)
			if err != nil {
				panic(err)
			}

			c.opts.Servers = append(c.opts.Servers, brokerURI)
		}
	}
}

// WithClientID will set the client id to be used by this client when
// connecting to the MQTT broker. According to the MQTT v3.1 specification,
// a client id mus be no longer than 23 characters.
func WithClientID(id string) func(*Client) {
	return func(c *Client) {
		c.opts.SetClientID(id)
	}
}

// WithUsername will set the username to be used by this client when connecting
// to the MQTT broker. Note: without the use of SSL/TLS, this information will
// be sent in plaintext accross the wire.
func WithUsername(username string) func(*Client) {
	return func(c *Client) {
		c.opts.SetUsername(username)
	}
}

// WithPassword will set the password to be used by this client when connecting
// to the MQTT broker. Note: without the use of SSL/TLS, this information will
// be sent in plaintext accross the wire.
func WithPassword(password string) func(*Client) {
	return func(c *Client) {
		c.opts.SetPassword(password)
	}
}

// WithTLSConfig will set an SSL/TLS configuration to be used when connecting
// to an MQTT broker. Please read the official Go documentation for more
// information.
func WithTLSConfig(t *tls.Config) func(*Client) {
	return func(c *Client) {
		c.opts.SetTLSConfig(t)
	}
}

// WithKeepAlive will set the amount of time (in seconds) that the client
// should wait before sending a PING request to the broker. This will
// allow the client to know that a connection has not been lost with the
// server.
func WithKeepAlive(k time.Duration) func(*Client) {
	return func(c *Client) {
		c.opts.SetKeepAlive(k)
	}
}

// WithPingTimeout will set the amount of time (in seconds) that the client
// will wait after sending a PING request to the broker, before deciding
// that the connection has been lost. Default is 10 seconds.
func WithPingTimeout(k time.Duration) func(*Client) {
	return func(c *Client) {
		c.opts.SetPingTimeout(k)
	}
}

// WithConnectTimeout limits how long the client will wait when trying to open a connection
// to an MQTT server before timeing out and erroring the attempt. A duration of 0 never times out.
// Default 30 seconds. Currently only operational on TCP/TLS connections.
func WithConnectTimeout(t time.Duration) func(*Client) {
	return func(c *Client) {
		c.opts.SetConnectTimeout(t)
	}
}

// WithMaxReconnectInterval sets the maximum time that will be waited between reconnection attempts
// when connection is lost
func WithMaxReconnectInterval(t time.Duration) func(*Client) {
	return func(c *Client) {
		c.opts.SetMaxReconnectInterval(t)
	}
}

// WithAutoReconnect sets whether the automatic reconnection logic should be used
// when the connection is lost, even if disabled the ConnectionLostHandler is still
// called
func WithAutoReconnect(a bool) func(*Client) {
	return func(c *Client) {
		c.opts.SetAutoReconnect(a)
	}
}

// option represents a key/value pair that can be supplied to the publish/subscribe or unsubscribe
// methods and provide ways to configure the operation.
type option string

const (
	withRetain = option("+r")
	withQos0   = option("+0")
	withQos1   = option("+1")
)

// String converts the option to a string.
func (o option) String() string {
	return string(o)
}

// WithoutEcho constructs an option which disables self-receiving messages if subscribed to a channel.
func WithoutEcho() Option {
	return option("me=0")
}

// WithTTL constructs an option which can be used during publish requests to set a Time-To-Live.
func WithTTL(seconds int) Option {
	return option("ttl=" + strconv.Itoa(seconds))
}

// WithLast constructs an option which can be used during subscribe requests to retrieve a message history.
func WithLast(messages int) Option {
	return option("last=" + strconv.Itoa(messages))
}

// WithRetain constructs an option which sets the message 'retain' flag to true.
func WithRetain() Option {
	return withRetain
}

// WithAtMostOnce instructs to publish at most once (MQTT QoS 0).
func WithAtMostOnce() Option {
	return withQos0
}

// WithAtLeastOnce instructs to publish at least once (MQTT QoS 1).
func WithAtLeastOnce() Option {
	return withQos1
}

func getUTCTimestamp(input time.Time) int64 {
	t := input
	if zone, _ := t.Zone(); zone != "UTC" {
		loc, _ := time.LoadLocation("UTC")
		t = t.In(loc)
	}
	return t.Unix()
}

// WithFrom request messages from a point in time.
func WithFrom(from time.Time) Option {
	return option("from=" + strconv.FormatInt(getUTCTimestamp(from), 10))
}

// WithUntil request messages until a point in time.
func WithUntil(until time.Time) Option {
	return option("until=" + strconv.FormatInt(getUTCTimestamp(until), 10))
}
