package peer

import (
	"github.com/muka/peer/enums"
	"github.com/muka/peer/models"
	"github.com/pion/webrtc/v3"
	"github.com/rs/xid"
)

// NewOptions return Peer options with defaults
func NewOptions() Options {
	return Options{
		Host:         "0.peerjs.com",
		Port:         443,
		PingInterval: 1000,
		Path:         "/",
		Secure:       true,
		Token:        xid.New().String(),
		Key:          DefaultKey,
		Configuration: webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{"stun:stun.l.google.com:19302"},
				},
			},
			SDPSemantics: webrtc.SDPSemanticsUnifiedPlan,
		},
		Debug: 0,
	}
}

//Options store Peer options
type Options struct {
	// Key API key for the cloud PeerServer. This is not used for servers other than 0.peerjs.com.
	Key string
	// Server host. Defaults to 0.peerjs.com. Also accepts '/' to signify relative hostname.
	Host string
	//Port Server port. Defaults to 443.
	Port int
	//PingInterval Ping interval in ms. Defaults to 5000.
	PingInterval int
	//Path The path where your self-hosted PeerServer is running. Defaults to '/'.
	Path string
	//Secure true if you're using SSL.
	Secure bool
	//Configuration hash passed to RTCPeerConnection. This hash contains any custom ICE/TURN server configuration. Defaults to { 'iceServers': [{ 'urls': 'stun:stun.l.google.com:19302' }], 'sdpSemantics': 'unified-plan' }
	Configuration webrtc.Configuration
	// Debug
	// Prints log messages depending on the debug level passed in. Defaults to 0.
	// 0 Prints no logs.
	// 1 Prints only errors.
	// 2 Prints errors and warnings.
	// 3 Prints all logs.
	Debug int8
	//Token a string to group peers
	Token string
}

// NewConnectionOptions return a ConnectionOptions with defaults
func NewConnectionOptions() *ConnectionOptions {
	return &ConnectionOptions{
		Serialization: enums.SerializationTypeRaw,
		Debug:         -1,
	}
}

//ConnectionOptions wrap optios for Peer Connect()
type ConnectionOptions struct {
	//ConnectionID
	ConnectionID string
	//Payload
	Payload models.Payload
	//Label A unique label by which you want to identify this data connection. If left unspecified, a label will be generated at random.
	Label string
	// Metadata associated with the connection, passed in by whoever initiated the connection.
	Metadata interface{}
	// Serialization Can be binary (default), binary-utf8, json, or none.
	Serialization string
	// Reliable whether the underlying data channels should be reliable (e.g. for large file transfers) or not (e.g. for gaming or streaming). Defaults to false.
	Reliable bool
	// Stream contains the reference to a media stream
	Stream *MediaStream
	// Originator indicate if the originator
	Originator bool
	// SDP contains SDP information
	SDP webrtc.SessionDescription
	// Debug level for debug taken. See Options
	Debug int8
	// SDPTransform transformation function for SDP message
	SDPTransform func(string) string
}

//AnswerOption wraps answer options
type AnswerOption struct {
	// SDPTransform transformation function for SDP message
	SDPTransform func(string) string
}
