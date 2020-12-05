package peer

import (
	"github.com/pion/webrtc/v3"
)

//MediaStreamTrack interaface that wraps together TrackLocal and TrackRemote
type MediaStreamTrack interface {
	// ID is the unique identifier for this Track. This should be unique for the
	// stream, but doesn't have to globally unique. A common example would be 'audio' or 'video'
	// and StreamID would be 'desktop' or 'webcam'
	ID() string

	// StreamID is the group this track belongs too. This must be unique
	StreamID() string

	// Kind controls if this TrackLocal is audio or video
	Kind() webrtc.RTPCodecType
}

// NewMediaStreamWithTrack create a mediastream with tracks
func NewMediaStreamWithTrack(tracks []MediaStreamTrack) *MediaStream {
	m := new(MediaStream)
	m.tracks = tracks
	return m
}

// MediaStream A stream of media content. A stream consists of several tracks
// such as video or audio tracks. Each track is specified as an instance
// of MediaStreamTrack.
type MediaStream struct {
	tracks []MediaStreamTrack
}

// GetTracks returns a list of tracks
func (m *MediaStream) GetTracks() []MediaStreamTrack {
	return m.tracks
}

// AddTrack add a track
func (m *MediaStream) AddTrack(t MediaStreamTrack) {
	m.tracks = append(m.tracks, t)
}

// RemoveTrack remove a track
func (m *MediaStream) RemoveTrack(t MediaStreamTrack) {
	tracks := []MediaStreamTrack{}
	for i, t1 := range m.tracks {
		if t1 == t {
			m.tracks[i] = nil
			continue
		}
		tracks = append(tracks, t1)
	}
}
