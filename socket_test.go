package peer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSocket(t *testing.T) {
	s := NewSocket(getTestOpts())
	done := false
	go func() {
		ev := <-s.Events
		assert.Equal(t, ev.Type, ConnectionEventTypeOpen)
		done = true
	}()
	err := s.Start("test", "test")
	assert.NoError(t, err)
	err = s.Close()
	assert.NoError(t, err)
	<-time.After(time.Millisecond * 100)
	assert.True(t, done)
}
