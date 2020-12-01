package peer

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSocket(t *testing.T) {
	s := NewSocket(getTestOpts())
	done := false
	s.On(SocketEventTypeMessage, func(data interface{}) {
		ev := data.(SocketEvent)
		assert.Equal(t, ev.Type, SocketEventTypeMessage)
		log.Println("socket received")
		done = true
	})
	err := s.Start("test", "test")
	assert.NoError(t, err)
	err = s.Close()
	assert.NoError(t, err)
	<-time.After(time.Millisecond * 500)
	assert.True(t, done)
}
