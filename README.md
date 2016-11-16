# Emitter Golang SDK
This repository contains Go/Golang client for [Emitter](https://emitter.io) (see also on [Emitter GitHub](https://github.com/emitter-io/emitter)). Emitter is an **open-source** real-time communication service for connecting online devices. Infrastructure and APIs for IoT, gaming, apps and real-time web. At its core, emitter.io is a distributed, scalable and fault-tolerant publish-subscribe messaging platform based on MQTT protocol and featuring message storage.

This library provides a nicer MQTT interface fine-tuned and extended with specific features provided by [Emitter](https://emitter.io). The code uses the [Eclipse Paho MQTT Go Client](https://github.com/eclipse/paho.mqtt.golang) for handling all the network communication and MQTT protocol, and is released under the same license (EPL v1).

## Usage

This library aims to be as simple and straighforward as possible. First thing you'll need to do is to import it.

```go
import emitter "github.com/emitter-io/go"
```

Then, you can use the functions exposed by `Emitter` type - they are simple methods such as `Connect`, `Publish`, `Subscribe`, `Unsubscribe`, `GenerateKey`, `Presence`, etc. See the example below.

```go
func main() {
	// Create the options with default values
	o := emitter.NewClientOptions()

	// Set the message handler
	o.SetOnMessageHandler(func(client emitter.Emitter, msg emitter.Message) {
		fmt.Printf("Received message: %s\n", msg.Payload())
	})

	// Set the presence notification handler
	o.SetOnPresenceHandler(func(_ emitter.Emitter, p emitter.PresenceEvent) {
		fmt.Printf("Occupancy: %v\n", p.Occupancy)
	})

	// Create a new emitter client and connect to the broker
	c := emitter.NewClient(o)
	sToken := c.Connect()
	if sToken.Wait() && sToken.Error() != nil {
		panic("Error on Client.Connect(): " + sToken.Error().Error())
	}

	// Subscribe to the presence demo channel
	c.Subscribe("X4-nUeHjiAygHMdN8wst82S3c2KcCMn7", "presence-demo/1")

	// Publish to the channel
	c.Publish("X4-nUeHjiAygHMdN8wst82S3c2KcCMn7", "presence-demo/1", "hello")

	// Ask for presence
	r := emitter.NewPresenceRequest()
	r.Key = "X4-nUeHjiAygHMdN8wst82S3c2KcCMn7"
	r.Channel = "presence-demo/1"
	c.Presence(r)
}
```

## Installation and Build

This client, similarly to the Eclipse Paho client is designed to work with the standard Go tools, so installation is as easy as:

```go
go get github.com/emitter-io/go
```

The client depends on Eclipse Paho MQTT Go Client and Google's websockets package, also easily installed with the command:

```go
go get github.com/eclipse/paho.mqtt.golang
go get golang.org/x/net/websocket
go get github.com/satori/go.uuid
```

## License

Licensed with EPL 1.0, similarly to Eclipse Paho MQTT Client.
