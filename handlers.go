package emitter

import "fmt"

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
