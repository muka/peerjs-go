package peer

import (
	"fmt"

	"github.com/muka/peerjs-go/enums"
	"github.com/muka/peerjs-go/models"
	"github.com/muka/peerjs-go/util"
	"github.com/pion/webrtc/v3"
)

//MediaChannelIDPrefix the media channel connection id prefix
const MediaChannelIDPrefix = "mc_"

//NewMediaConnection create new MediaConnection
func NewMediaConnection(id string, peer *Peer, opts ConnectionOptions) (*MediaConnection, error) {

	m := &MediaConnection{
		BaseConnection: newBaseConnection(enums.ConnectionTypeMedia, peer, opts),
	}

	m.peerID = id

	m.id = opts.ConnectionID
	if m.id == "" {
		m.id = fmt.Sprintf("%s%s", MediaChannelIDPrefix, util.RandomToken())
	}

	m.localStream = opts.Stream

	m.negotiator = NewNegotiator(m, opts)
	var err error
	if m.localStream != nil {
		opts.Originator = true
		err = m.negotiator.StartConnection(opts)
	}

	return m, err
}

// MediaConnection track a connection with a remote Peer
type MediaConnection struct {
	BaseConnection
	open         bool
	remoteStream *MediaStream
	localStream  *MediaStream
}

// GetLocalStream returns the local stream
func (m *MediaConnection) GetLocalStream() *MediaStream {
	return m.localStream
}

// GetRemoteStream returns the remote stream
func (m *MediaConnection) GetRemoteStream() *MediaStream {
	return m.remoteStream
}

// AddStream adds a stream to the MediaConnection
func (m *MediaConnection) AddStream(tr *webrtc.TrackRemote) {
	m.log.Debugf("Receiving stream: %v", tr)
	m.remoteStream = NewMediaStreamWithTrack([]MediaStreamTrack{tr})
	m.Emit(enums.ConnectionEventTypeStream, tr)
}

func (m *MediaConnection) handleMessage(message models.Message) {
	mtype := message.GetType()
	payload := message.GetPayload()
	switch message.GetType() {
	case enums.ServerMessageTypeAnswer:
		// Forward to negotiator
		m.negotiator.handleSDP(message.GetType(), *payload.SDP)
		m.open = true
		break
	case enums.ServerMessageTypeCandidate:
		m.negotiator.HandleCandidate(payload.Candidate)
		break
	default:
		m.log.Warnf("Unrecognized message type:%s from peer:%s", mtype, m.peerID)
		break
	}
}

//Answer open the media connection with the remote peer
func (m *MediaConnection) Answer(tl webrtc.TrackLocal, options *AnswerOption) {

	if m.localStream != nil {
		m.log.Warn("Local stream already exists on this MediaConnection. Are you answering a call twice?")
		return
	}

	stream := NewMediaStreamWithTrack([]MediaStreamTrack{tl})
	m.localStream = stream

	if options != nil && options.SDPTransform != nil {
		m.BaseConnection.opts.SDPTransform = options.SDPTransform
	}

	connOpts := m.GetOptions()
	connOpts.Stream = stream
	m.negotiator.StartConnection(connOpts)
	// Retrieve lost messages stored because PeerConnection not set up.
	messages := m.GetProvider().GetMessages(m.GetID())

	for _, message := range messages {
		m.handleMessage(message)
	}

	m.open = true
}

//Close allows user to close connection
func (m *MediaConnection) Close() error {
	if m.negotiator != nil {
		m.negotiator.Cleanup()
		m.negotiator = nil
	}

	m.localStream = nil
	m.remoteStream = nil

	if m.GetProvider() != nil {
		m.GetProvider().RemoveConnection(m)
		m.BaseConnection.Provider = nil
	}

	if m.BaseConnection.opts.Stream != nil {
		m.BaseConnection.opts.Stream = nil
	}

	if !m.open {
		return nil
	}

	m.open = false

	m.Emit(enums.ConnectionEventTypeClose, nil)
	return nil
}
