package peer

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestOpts() Options {
	opts := NewOptions()
	opts.Path = "/myapp"
	opts.Host = "localhost"
	opts.Port = 9000
	opts.Secure = false
	opts.Debug = 0
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

func TestHelloWorld(t *testing.T) {

	peer1, err := NewPeer("peer1", getTestOpts())
	assert.NoError(t, err)
	defer peer1.Close()

	peer2, err := NewPeer("peer2", getTestOpts())
	assert.NoError(t, err)
	defer peer2.Close()

	// done := false
	done := make(chan bool)
	peer2.On("connection", func(data interface{}) {
		conn2 := data.(*DataConnection)
		conn2.On("data", func(data interface{}) {
			// Will print 'hi!'
			log.Printf("Received: %v\n", data)
			done <- true
		})
	})

	conn1, err := peer1.Connect("peer2", nil)
	assert.NoError(t, err)
	conn1.On("open", func(data interface{}) {
		for {
			conn1.Send([]byte("hi!"), false)
			<-time.After(time.Millisecond * 1000)
		}
	})

	// // <-time.After(time.Millisecond * 100)
	// // assert.True(t, done)
	log.Print("Waiting for event")
	<-done
	log.Print("Exiting..")
}
