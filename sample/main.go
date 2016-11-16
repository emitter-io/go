package main

import (
	"fmt"
	"time"

	emitter "github.com/emitter-io/go"
)

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

	// stop after 10 seconds
	time.Sleep(10 * time.Second)
}
