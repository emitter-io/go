package emitter

// KeyGenRequest represents a request that can be sent to emitter broker
// in order to generate a new channel key.
type KeyGenRequest struct {
	Key     string
	Channel string
	Type    string
	TTL     int
}

// KeyGenResponse  represents a response from emitter broker which contains
// the response to the key generation request.
type KeyGenResponse struct {
	Status  int
	Key     string
	Channel string
}

// PresenceRequest represents a request that can be sent to emitter broker
// in order to request presence information.
type PresenceRequest struct {
	Key     string
	Channel string
	Status  bool
	Changes bool
}

// PresenceEvent  represents a response from emitter broker which contains
// presence state or a join/leave notification.
type PresenceEvent struct {
	Event     string
	Channel   string
	Occupancy int
	Time      int
	Who       []PresenceInfo
}

// PresenceInfo represents a response from emitter broker which contains
// presence information.
type PresenceInfo struct {
	ID       string
	Username string
}

// NewKeyGenRequest creates a new KeyGenRequest type with some default values.
func NewKeyGenRequest() *KeyGenRequest {
	o := &KeyGenRequest{
		Key:     "",
		Channel: "",
		Type:    "",
		TTL:     0,
	}
	return o
}

// NewPresenceRequest creates a new PresenceRequest type with some default values.
func NewPresenceRequest() *PresenceRequest {
	o := &PresenceRequest{
		Key:     "",
		Channel: "",
		Status:  true,
		Changes: false,
	}
	return o
}
