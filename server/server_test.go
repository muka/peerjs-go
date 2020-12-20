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
		if err != nil {
			t.Logf("Start error: %v", err)
		}
	}()

	<-time.After(time.Millisecond * 400)
	err := p.Stop()
	assert.NoError(t, err)
	<-time.After(time.Millisecond * 100)
}
