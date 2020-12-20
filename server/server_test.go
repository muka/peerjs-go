package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPeerServer_StartStop(t *testing.T) {
	opts := NewOptions()
	opts.Port = 64888
	p := New(opts)

	go func() {
		err := p.Start()
		assert.NoError(t, err)
	}()

	<-time.After(time.Second * 1)
	err := p.Stop()
	assert.NoError(t, err)
}
