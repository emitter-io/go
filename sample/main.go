package main

import (
	"fmt"
	"time"

	emitter "github.com/emitter-io/go/v2"
)

func main() {
	clientA()
	clientB()

	// stop after 10 seconds
	time.Sleep(1 * time.Second)
}

func clientA() {
	const key = "RUvY5GTEOUmIqFs_zfpJcfTqBUIKBhfs" // read on sdk-integration-test/#/

	// Create the client and connect to the broker
	c, _ := emitter.Connect("", func(_ *emitter.Client, msg emitter.Message) {
		fmt.Printf("[emitter] -> [A] received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})

	// Subscribe to sdk-integration-test channel
	fmt.Println("[emitter] <- [A] subscribing to 'sdk-integration-test/...'")
	c.Subscribe(key, "sdk-integration-test/", nil)
}

func clientB() {
	const key = "pGrtRRL6RrjAdExSArkMzBZOoWr2eB3L" // everything on sdk-integration-test/

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
}
