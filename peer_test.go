package peer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestOpts() Options {
	opts := NewOptions()
	opts.Path = "/myapp"
	opts.Host = "localhost"
	opts.Port = 9000
	opts.Secure = false
	opts.Debug = 3
	return opts
}

func TestNewPeer(t *testing.T) {
	p, err := NewPeer("test", getTestOpts())
	assert.NoError(t, err)
	assert.NotEmpty(t, p.ID)
	p.Close()
}

func TestNewPeerEvents(t *testing.T) {
	p, err := NewPeer("test", getTestOpts())
	// done := false
	// p.On(PeerEventTypeOpen, func(data interface{}) {
	// 	done = true
	// })
	assert.NoError(t, err)
	assert.NotEmpty(t, p.ID)

	p.Close()
	// <-time.After(time.Millisecond * 1000)
	// assert.True(t, done)
}
