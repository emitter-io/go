package emitter

// KeyGenRequest represents a request that can be sent to emitter broker
// in order to generate a new channel key.
type KeyGenRequest struct {
	Key     string `json:"key"`
	Channel string `json:"channel"`
	Type    string `json:"type"`
	TTL     int    `json:"ttl"`
}

// KeyGenResponse  represents a response from emitter broker which contains
// the response to the key generation request.
type KeyGenResponse struct {
	Status       int    `json:"status"`
	Key          string `json:"key"`
	Channel      string `json:"channel"`
	ErrorMessage string `json:"message"`
}

// PresenceRequest represents a request that can be sent to emitter broker
// in order to request presence information.
type PresenceRequest struct {
	Key     string `json:"key"`
	Channel string `json:"channel"`
	Status  bool   `json:"status"`
	Changes bool   `json:"changes"`
}

// PresenceEvent  represents a response from emitter broker which contains
// presence state or a join/leave notification.
type PresenceEvent struct {
	Event   string         `json:"event"`
	Channel string         `json:"channel"`
	Time    int            `json:"time"`
	Who     []PresenceInfo `json:"who"`
}

// PresenceInfo represents a response from emitter broker which contains
// presence information.
type PresenceInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// LinkRequest represents a request to create a link.
type linkRequest struct {
	Name      string `json:"name"`      // The name of the shortcut, max 2 characters.
	Key       string `json:"key"`       // The key for the channel.
	Channel   string `json:"channel"`   // The channel name for the shortcut.
	Subscribe bool   `json:"subscribe"` // Specifies whether the broker should auto-subscribe.
	Private   bool   `json:"private"`   // Specifies whether the broker should generate a private link.
}

// LinkResponse represents a response for the link creation.
type LinkResponse struct {
	Status  int    `json:"status"`            // The status of the response.
	Name    string `json:"name,omitempty"`    // The name of the shortcut, max 2 characters.
	Channel string `json:"channel,omitempty"` // The channel which was registered.
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
