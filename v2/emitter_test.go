package emitter

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEndToEnd(t *testing.T) {
	clientA(t)
	clientB(t)

	// stop after 10 seconds
	time.Sleep(1 * time.Second)
}

func clientA(t *testing.T) {
	const key = "RUvY5GTEOUmIqFs_zfpJcfTqBUIKBhfs" // read on sdk-integration-test/#/

	// Create the client and connect to the broker
	c, _ := Connect("", func(_ *Client, msg Message) {
		fmt.Printf("[emitter] -> [A] received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})

	// Subscribe to demo channel
	fmt.Println("[emitter] <- [A] subscribing to 'demo/...'")
	err := c.Subscribe(key, "sdk-integration-test/", nil)
	assert.NoError(t, err)
}

func clientB(t *testing.T) {
	const key = "pGrtRRL6RrjAdExSArkMzBZOoWr2eB3L" // everything on sdk-integration-test/

	// Create the client and connect to the broker
	c, _ := Connect("", func(_ *Client, msg Message) {
		fmt.Printf("[emitter] -> [B] received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})

	c.OnPresence(func(_ *Client, ev PresenceEvent) {
		fmt.Printf("[emitter] -> [B] presence event: %d subscriber(s) at topic: '%s'\n", len(ev.Who), ev.Channel)
	})

	fmt.Println("[emitter] <- [B] querying own name")
	id := c.ID()
	assert.NotEmpty(t, id)

	// Subscribe to demo channel
	c.Subscribe(key, "sdk-integration-test/", func(_ *Client, msg Message) {
		fmt.Printf("[emitter] -> [B] received on specific handler: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})

	// Ask for presence
	fmt.Println("[emitter] <- [B] asking for presence on 'sdk-integration-test/'")
	err := c.Presence(key, "sdk-integration-test/", true, false)
	assert.NoError(t, err)

	// Publish to the channel
	fmt.Println("[emitter] <- [B] publishing to 'sdk-integration-test/'")
	err = c.Publish(key, "sdk-integration-test/", "hello")
	assert.NoError(t, err)

	// Ask to create a private link
	fmt.Println("[emitter] <- [B] creating a private link")
	link, err := c.CreatePrivateLink(key, "sdk-integration-test/", "1", true)
	assert.NoError(t, err)
	assert.NotNil(t, link)

	// Publish to the private link
	fmt.Println("[emitter] <- [B] publishing to private link")
	c.PublishWithLink("1", "hi from private link")
}
