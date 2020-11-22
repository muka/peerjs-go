package peer

import "github.com/pion/webrtc/v3"

//NewMediaConnection create new MediaConnection
func NewMediaConnection(id string, peer *Peer, opts ConnectionOptions) (*MediaConnection, error) {
	return &MediaConnection{
		BaseConnection: newBaseConnection(ConnectionTypeMedia, peer, opts),
	}, nil
}

// MediaConnection track a connection with a remote Peer
type MediaConnection struct {
	BaseConnection
}

// AddStream adds a stream to the MediaConnection
func (m *MediaConnection) AddStream(tr *webrtc.TrackRemote) {
	panic("TODO")
}
