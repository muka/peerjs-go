package server

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/muka/peerjs-go"
	"github.com/stretchr/testify/assert"
)

func TestPeerServer_StartStop(t *testing.T) {
	opts := NewOptions()
	opts.Port = 9001
	p := New(opts)

	var err error
	go func() {
		err = p.Start()
	}()

	<-time.After(time.Millisecond * 500)
	if err != nil {
		t.Logf("Start error: %v", err)
		assert.NoError(t, err)
	}
	err = p.Stop()
	assert.NoError(t, err)
}

func TestPeerServer_ClientPingPong(t *testing.T) {

	opts := NewOptions()
	opts.Port = 9001
	opts.Path = "/myapp"
	server := New(opts)
	defer server.Stop()

	peerOpts := peer.NewOptions()
	peerOpts.Host = opts.Host
	peerOpts.Port = opts.Port
	peerOpts.Path = opts.Path
	peerOpts.Secure = false

	var err error
	go func() {
		err = server.Start()
	}()

	<-time.After(time.Millisecond * 1000)

	if err != nil {
		t.FailNow()
	}

	peer1, err := peer.NewPeer("peer1", peerOpts)
	assert.NoError(t, err)
	assert.NotNil(t, peer1)
	if peer1 != nil {
		defer peer1.Close()
	}

	peer2, err := peer.NewPeer("peer2", peerOpts)
	assert.NoError(t, err)
	assert.NotNil(t, peer2)
	if peer2 != nil {
		defer peer2.Close()
	}

	done := make(chan error)
	peer2.On("connection", func(data interface{}) {
		conn2 := data.(*peer.DataConnection)
		conn2.On("data", func(data interface{}) {
			log.Printf("Received\n")
			done <- nil
		})
	})

	conn1, err := peer1.Connect("peer2", nil)
	assert.NoError(t, err)
	if err != nil {
		t.Fatal(err)
	}
	conn1.On("open", func(data interface{}) {
		conn1.Send([]byte("hello world"), false)
	})

	go func() {
		<-time.After(time.Millisecond * 5000)
		done <- errors.New("Timeout")
	}()

	err = <-done
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

}
