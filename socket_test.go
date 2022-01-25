package peer

import (
	"log"
	"testing"
	"time"

	"github.com/muka/peerjs-go/enums"
	"github.com/stretchr/testify/assert"
)

func TestNewSocket(t *testing.T) {
	srv, srvOpts := startServer()
	srv.Start()
	defer srv.Stop()
	s := NewSocket(getTestOpts(srvOpts))
	done := false
	s.On(enums.SocketEventTypeMessage, func(data interface{}) {
		ev := data.(SocketEvent)
		assert.Equal(t, ev.Type, enums.SocketEventTypeMessage)
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
