package peer

import (
	"github.com/muka/peerjs-go/emitter"
	"github.com/muka/peerjs-go/models"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

//Connection shared interface
type Connection interface {
	GetType() string
	GetID() string
	GetPeerID() string
	GetProvider() *Peer
	GetMetadata() interface{}
	GetPeerConnection() *webrtc.PeerConnection
	SetPeerConnection(pc *webrtc.PeerConnection)
	Close() error
	HandleMessage(*models.Message) error
	Emit(string, interface{})
	GetOptions() ConnectionOptions
}

func newBaseConnection(connType string, peer *Peer, opts ConnectionOptions) BaseConnection {
	return BaseConnection{
		Emitter:    emitter.NewEmitter(),
		Type:       connType,
		Provider:   peer,
		log:        createLogger(connType, opts.Debug),
		opts:       opts,
		negotiator: nil,
	}
}

// BaseConnection shared base connection
type BaseConnection struct {
	emitter.Emitter
	// id connection ID
	id string
	// peerID peer ID of the connection
	peerID string
	//Provider is the peer instance
	Provider *Peer
	// DataChannel A reference to the RTCDataChannel object associated with the connection.
	DataChannel *webrtc.DataChannel
	// The optional label passed in or assigned by PeerJS when the connection was initiated.
	Label string
	// Metadata Any type of metadata associated with the connection, passed in by whoever initiated the connection.
	Metadata interface{}
	// Open This is true if the connection is open and ready for read/write.
	Open bool
	// PeerConnection A reference to the RTCPeerConnection object associated with the connection.
	PeerConnection *webrtc.PeerConnection
	// Reliable Whether the underlying data channels are reliable; defined when the connection was initiated.
	Reliable bool
	// Serialization The serialization format of the data sent over the connection. Can be binary (default), binary-utf8, json, or none.
	Serialization string
	// Type defines the type for connections
	Type string
	// BufferSize The number of messages queued to be sent once the browser buffer is no longer full.
	BufferSize int
	opts       ConnectionOptions
	log        *logrus.Entry
	negotiator *Negotiator
}

// GetOptions return the connection configuration
func (c *BaseConnection) GetOptions() ConnectionOptions {
	return c.opts
}

// GetMetadata return the connection metadata
func (c *BaseConnection) GetMetadata() interface{} {
	return c.Metadata
}

// GetPeerConnection return the underlying WebRTC PeerConnection
func (c *BaseConnection) GetPeerConnection() *webrtc.PeerConnection {
	return c.PeerConnection
}

// SetPeerConnection set the underlying WebRTC PeerConnection
func (c *BaseConnection) SetPeerConnection(pc *webrtc.PeerConnection) {
	c.PeerConnection = pc
	c.log.Debugf("%v", c.PeerConnection)
}

// GetID return the connection ID
func (c *BaseConnection) GetID() string {
	return c.id
}

// GetPeerID return the connection peer ID
func (c *BaseConnection) GetPeerID() string {
	return c.peerID
}

// Close closes the data connection
func (c *BaseConnection) Close() error {
	panic("Not implemented!")
}

// HandleMessage handles incoming messages
func (c *BaseConnection) HandleMessage(msg *models.Message) error {
	panic("Not implemented!")
}

// GetType return the connection type
func (c *BaseConnection) GetType() string {
	return c.Type
}

// GetProvider return the peer provider
func (c *BaseConnection) GetProvider() *Peer {
	return c.Provider
}
