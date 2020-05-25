// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package emitter

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
)

// MessageHandler is a callback type which can be set to be
// executed upon the arrival of messages published to topics
// to which the client is subscribed.
type MessageHandler func(*Client, Message)

// PresenceHandler is a callback type which can be set to be executed upon
// the arrival of presence events.
type PresenceHandler func(*Client, PresenceEvent)

// ErrorHandler is a callback type which can be set to be executed upon
// the arrival of an emitter error.
type ErrorHandler func(*Client, Error)

// DisconnectHandler is a callback type which can be set to be
// executed upon an unintended disconnection from the MQTT broker.
// Disconnects caused by calling Disconnect or ForceDisconnect will
// not cause an OnConnectionLost callback to execute.
type DisconnectHandler func(*Client, error)

// ConnectHandler is a callback that is called when the client
// state changes from unconnected/disconnected to connected. Both
// at initial connection and on reconnection
type ConnectHandler func(*Client)

// Option represents a key/value pair that can be supplied to the publish/subscribe or unsubscribe
// methods and provide ways to configure the operation.
type Option interface {
	String() string
}

// Error represents an event code which provides a more details.
type Error struct {
	Request uint16 `json:"req,omitempty"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}

// RequestID returns the request ID for the response.
func (e *Error) RequestID() uint16 {
	return e.Request
}

// ------------------------------------------------------------------------------------

// KeyGenRequest represents a request that can be sent to emitter broker
// in order to generate a new channel key.
type keygenRequest struct {
	Key     string `json:"key"`
	Channel string `json:"channel"`
	Type    string `json:"type"`
	TTL     int    `json:"ttl"`
}

// KeyGenResponse  represents a response from emitter broker which contains
// the response to the key generation request.
type keyGenResponse struct {
	Request      uint16 `json:"req,omitempty"`
	Status       int    `json:"status"`
	Key          string `json:"key"`
	Channel      string `json:"channel"`
	ErrorMessage string `json:"message"`
}

// RequestID returns the request ID for the response.
func (r *keyGenResponse) RequestID() uint16 {
	return r.Request
}

// ------------------------------------------------------------------------------------

// PresenceRequest represents a request that can be sent to emitter broker
// in order to request presence information.
type presenceRequest struct {
	Key     string `json:"key"`
	Channel string `json:"channel"`
	Status  bool   `json:"status"`
	Changes bool   `json:"changes"`
}

// presenceMessage represents a presence message, for partial unmarshal
type presenceMessage struct {
	Request uint16          `json:"req,omitempty"`
	Event   string          `json:"event"`
	Channel string          `json:"channel"`
	Time    int             `json:"time"`
	Who     json.RawMessage `json:"who"`
}

// PresenceEvent  represents a response from emitter broker which contains
// presence state or a join/leave notification.
type PresenceEvent struct {
	presenceMessage
	Who []PresenceInfo
}

// RequestID returns the request ID for the response.
func (r *PresenceEvent) RequestID() uint16 {
	return r.Request
}

// PresenceInfo represents a response from emitter broker which contains
// presence information.
type PresenceInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// ------------------------------------------------------------------------------------

// meResponse represents information about the client.
type meResponse struct {
	Request uint16            `json:"req,omitempty"`   // The corresponding request ID.
	ID      string            `json:"id"`              // The private ID of the connection.
	Links   map[string]string `json:"links,omitempty"` // The set of pre-defined channels.
}

// RequestID returns the request ID for the response.
func (r *meResponse) RequestID() uint16 {
	return r.Request
}

// ------------------------------------------------------------------------------------

// linkRequest represents a request to create a link.
type linkRequest struct {
	Name      string `json:"name"`      // The name of the shortcut, max 2 characters.
	Key       string `json:"key"`       // The key for the channel.
	Channel   string `json:"channel"`   // The channel name for the shortcut.
	Subscribe bool   `json:"subscribe"` // Specifies whether the broker should auto-subscribe.
}

// Link represents a response for the link creation.
type Link struct {
	Request uint16 `json:"req,omitempty"`
	Name    string `json:"name,omitempty"`    // The name of the shortcut, max 2 characters.
	Channel string `json:"channel,omitempty"` // The channel which was registered.
}

// RequestID returns the request ID for the response.
func (r *Link) RequestID() uint16 {
	return r.Request
}

// ------------------------------------------------------------------------------------

// KeyBanRequest represents a request that can be sent to emitter broker
// in order to ban/blacklist a channel key.
type keybanRequest struct {
	Secret string `json:"secret"` // The master key to use.
	Target string `json:"target"` // The target key to ban.
	Banned bool   `json:"banned"` // Whether the target should be banned or not.
}

// keyBanResponse  represents a response from emitter broker which contains
// the response to the key ban request.
type keyBanResponse struct {
	Request      uint16 `json:"req,omitempty"`
	Status       int    `json:"status"` // The status of the response
	Banned       bool   `json:"banned"` // Whether the target should be banned or not.
	ErrorMessage string `json:"message"`
}

// RequestID returns the request ID for the response.
func (r *keyBanResponse) RequestID() uint16 {
	return r.Request
}

// ------------------------------------------------------------------------------------

// uuid generates a simple UUID
func uuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
