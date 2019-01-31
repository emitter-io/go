package main

import (
	"fmt"
	"time"

	emitter "github.com/emitter-io/go"
)

// demo license: zT83oDV0DWY5_JysbSTPTDr8KB0AAAAAAAAAAAAAAAI
// demo secret key: kBCZch5re3Ue-kpG1Aa8Vo7BYvXZ3UwR
func main() {
	clientA()
	clientB()

	// stop after 10 seconds
	time.Sleep(10 * time.Second)
}

func clientA() {
	const key = "8XR27YmSfeJXj4jU6vh22uGLqLhxGfU5" // read on demo/#/

	// Create client options
	o := emitter.NewClientOptions()
	o.AddBroker("tcp://127.0.0.1:8080")
	o.SetOnMessageHandler(func(client emitter.Emitter, msg emitter.Message) {
		fmt.Printf("[emitter] -> [A] received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})

	// Create a new emitter client and connect to the broker
	c := emitter.NewClient(o)
	c.Connect().WaitTimeout(time.Second)

	// Subscribe to demo channel
	fmt.Println("[emitter] <- [A] subscribing to 'demo/...'")
	c.Subscribe(key, "demo/").WaitTimeout(time.Second)
}

func clientB() {
	const key = "tGxkIMapzyQx5Cc7koX5NVtQV7tA8tMw" // everything on demo/

	// Create the options with default values
	o := emitter.NewClientOptions()
	o.AddBroker("tcp://127.0.0.1:8080")

	// Set the message handler
	o.SetOnMessageHandler(func(client emitter.Emitter, msg emitter.Message) {
		fmt.Printf("[emitter] -> [B] received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})

	// Set the presence notification handler
	o.SetOnPresenceHandler(func(_ emitter.Emitter, p emitter.PresenceEvent) {
		fmt.Printf("[emitter] -> [B] received presence: %v\n", len(p.Who))
	})

	// Set the link notification handler
	o.SetOnLinkHandler(func(_ emitter.Emitter, p emitter.LinkResponse) {
		fmt.Printf("[emitter] -> [B] created private link '%v' for '%v'\n", p.Name, p.Channel)
	})

	// Create a new emitter client and connect to the broker
	c := emitter.NewClient(o)
	c.Connect().WaitTimeout(time.Second)

	// Subscribe to demo channel
	fmt.Println("[emitter] <- [B] subscribing to 'demo/'")
	c.Subscribe(key, "demo/").WaitTimeout(time.Second)

	// Publish to the channel
	fmt.Println("[emitter] <- [B] publishing to 'demo/'")
	c.Publish(key, "demo/", "hello").WaitTimeout(time.Second)

	// Ask for presence
	r := emitter.NewPresenceRequest()
	r.Key = key
	r.Channel = "demo/"
	fmt.Println("[emitter] <- [B] asking for presence on 'demo/'")
	c.Presence(r).WaitTimeout(time.Second)

	// Ask to create a private link
	fmt.Println("[emitter] <- [B] creating a private link")
	c.CreatePrivateLink(key, "demo/", "1", true).WaitTimeout(time.Second)

	// Publish to the private link
	fmt.Println("[emitter] <- [B] publishing to private link")
	c.PublishWithLink("1", "hi from private link").WaitTimeout(time.Second)
}
