# Emitter Golang SDK [![api documentation](http://b.repl.ca/v1/api-documentation-green.png)](https://godoc.org/github.com/emitter-io/go)
This repository contains Go/Golang client for [Emitter](https://emitter.io) (see also on [Emitter GitHub](https://github.com/emitter-io/emitter)). Emitter is an **open-source** real-time communication service for connecting online devices. At its core, emitter.io is a distributed, scalable and fault-tolerant publish-subscribe messaging platform based on MQTT protocol and featuring message storage.

This library provides a nicer MQTT interface fine-tuned and extended with specific features provided by [Emitter](https://emitter.io). The code uses the [Eclipse Paho MQTT Go Client](https://github.com/eclipse/paho.mqtt.golang) for handling all the network communication and MQTT protocol, and is released under the same license (EPL v1). 

## Usage

This library aims to be as simple and straighforward as possible. First thing you'll need to do is to import it.

```go
import emitter "github.com/emitter-io/go/v2"
```

Then, you can use the functions exposed by `Emitter` type - they are simple methods such as `Connect`, `Publish`, `Subscribe`, `Unsubscribe`, `GenerateKey`, `Presence`, etc. See the example below.

```go
func main() {


	// Create the client and connect to the broker
	c, _ := emitter.Connect("", func(_ *emitter.Client, msg emitter.Message) {
		fmt.Printf("[emitter] -> [B] received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})

	// Set the presence handler
	c.OnPresence(func(_ *emitter.Client, ev emitter.PresenceEvent) {
		fmt.Printf("[emitter] -> [B] presence event: %d subscriber(s) at topic: '%s'\n", len(ev.Who), ev.Channel)
	})

	fmt.Println("[emitter] <- [B] querying own name")
	id := c.ID()
	fmt.Println("[emitter] -> [B] my name is " + id)

	// Subscribe to sdk-integration-test channel
	fmt.Println("[emitter] <- [B] subscribing to 'sdk-integration-test/'")
	c.Subscribe(key, "sdk-integration-test/", func(_ *emitter.Client, msg emitter.Message) {
		fmt.Printf("[emitter] -> [B] received on specific handler: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})

	// Ask for presence
	fmt.Println("[emitter] <- [B] asking for presence on 'sdk-integration-test/'")
	c.Presence(key, "sdk-integration-test/", true, false)

	// Publish to the channel
	fmt.Println("[emitter] <- [B] publishing to 'sdk-integration-test/'")
	c.Publish(key, "sdk-integration-test/", "hello")

	select {}  // wait forever without busy loop - Ctrl-c to stop
}
```

## Installation and Build

This client, similarly to the Eclipse Paho client is designed to work with the standard Go tools, so installation is as easy as:

```go
go get -u github.com/emitter-io/go/v2
```

For usage, please refer to the `sample` sub-folder in this repository which provides a sample application on how to use the API.

## API Documentation

The full API documentation of exported members is available on [godoc.org/github.com/emitter-io/go/v2](https://godoc.org/github.com/emitter-io/go/v2).

## License

Licensed with EPL 1.0, similarly to Eclipse Paho MQTT Client.
