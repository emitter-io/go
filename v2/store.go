package emitter

import (
	"fmt"
	"sync"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

// In-memory storage implementation
type store struct {
	sync.RWMutex
	messages map[string]*packet
}

// Response represents a generic response sent by the broker.
type Response interface {
	RequestID() uint16
}

// The control packet with an optional error and a response.
type packet struct {
	packets.ControlPacket
	callback chan Response
}

// newStore creates a new message storage layer.
func newErrorStore() *store {
	store := &store{
		messages: make(map[string]*packet),
	}
	return store
}

// Open initializes a MemoryStore instance.
func (store *store) Open() {}

// Put takes a key and a pointer to a Message and stores the message.
func (store *store) Put(key string, message packets.ControlPacket) {
	store.Lock()
	defer store.Unlock()

	store.messages[key] = &packet{
		ControlPacket: message,
	}
}

// Get takes a key and looks in the store for a matching Message
// returning either the Message pointer or nil.
func (store *store) Get(key string) packets.ControlPacket {
	store.RLock()
	defer store.RUnlock()
	return store.messages[key]
}

// All returns a slice of strings containing all the keys currently
// in the MemoryStore.
func (store *store) All() []string {
	store.RLock()
	defer store.RUnlock()

	keys := []string{}
	for k := range store.messages {
		keys = append(keys, k)
	}
	return keys
}

// Del takes a key, searches the MemoryStore and if the key is found
// deletes the Message pointer associated with it.
func (store *store) Del(key string) {
	store.Lock()
	defer store.Unlock()

	m := store.messages[key]
	if m != nil && m.callback == nil {
		delete(store.messages, key)
	}
}

// Close will disallow modifications to the state of the store.
func (store *store) Close() {}

// Reset eliminates all persisted message data in the store.
func (store *store) Reset() {
	store.Lock()
	defer store.Unlock()
	store.messages = make(map[string]*packet)
}

// PutCallback adds a callback channel to a message.
func (store *store) PutCallback(id uint16) <-chan Response {
	store.Lock()
	defer store.Unlock()

	key := outboundKeyFromMID(id)
	if m, ok := store.messages[key]; ok && m != nil {
		m.callback = make(chan Response, 1)
		return m.callback
	}
	return nil
}

// NotifyResponse notifies a response on a callback (if exists)
func (store *store) NotifyResponse(id uint16, response Response) bool {
	store.RLock()
	defer store.RUnlock()

	key := outboundKeyFromMID(id)
	if m, ok := store.messages[key]; ok && m != nil {
		m.callback <- response
		close(m.callback)
		delete(store.messages, key)
		return true
	}
	return false
}

// Return a string of the form "o.[id]"
func outboundKeyFromMID(id uint16) string {
	return fmt.Sprintf("o.%d", id)
}
