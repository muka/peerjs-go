package peer

import "github.com/pion/webrtc/v3"

// MessagePayload wraps a message payload
type MessagePayload struct {
	Type          string      `json:"type"`
	ConnectionID  string      `json:"connectionId"`
	Metadata      interface{} `json:"metadata,omitempty"`
	Label         string      `json:"label,omitempty"`
	Serialization string      `json:"serialization,omitempty"`
	Reliable      bool        `json:"reliable,omitempty"`
	Candidate     string      `json:"candidate,omitempty"`
}

// Message shared message interface
type Message interface {
	GetPayload() MessagePayload
	GetType() string
	GetSrc() string
	GetDst() string
}

// BaseMessage Message implementation
type BaseMessage struct {
	Type    string         `json:"type"`
	Payload MessagePayload `json:"payload"`
	Src     string         `json:"src"`
	Dst     string         `json:"dst,omitempty"`
}

// GetPayload returns the message payload
func (m BaseMessage) GetPayload() MessagePayload {
	return m.Payload
}

// GetSrc returns the message src
func (m BaseMessage) GetSrc() string {
	return m.Src
}

// GetDst returns the message dst
func (m BaseMessage) GetDst() string {
	return m.Dst
}

// GetType returns the message payload
func (m BaseMessage) GetType() string {
	return m.Type
}

// Exchange message

// ExchangeMessage ANSWER Message structure
type ExchangeMessage struct {
	BaseMessage
	Payload ExchangePayload `json:"payload"`
}

//ExchangePayload wrap information to exchange session information
type ExchangePayload struct {
	MessagePayload
	SDP     webrtc.SessionDescription `json:"sdp"`
	Browser string                    `json:"browser"`
}
