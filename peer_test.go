package peer

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/muka/peerjs-go/server"
	"github.com/muka/peerjs-go/util"
	"github.com/pion/webrtc/v3"
	"github.com/stretchr/testify/assert"
)

func rndName(name string) string {
	return fmt.Sprintf("%s_%s", name, util.RandomToken())
}

func getTestOpts(serverOpts server.Options) Options {
	opts := NewOptions()
	opts.Path = serverOpts.Path
	opts.Host = serverOpts.Host
	opts.Port = serverOpts.Port
	opts.Secure = false
	opts.Debug = 0
	return opts
}

func startServer() (*server.PeerServer, server.Options) {
	opts := server.NewOptions()
	opts.Port = 9000
	opts.Host = "localhost"
	opts.Path = "/myapp"
	return server.New(opts), opts
}

func TestNewPeer(t *testing.T) {
	peerServer, serverOpts := startServer()
	err := peerServer.Start()
	if err != nil {
		t.Logf("Server error: %s", err)
		t.FailNow()
	}
	defer peerServer.Stop()
	p, err := NewPeer("test", getTestOpts(serverOpts))
	assert.NoError(t, err)
	assert.NotEmpty(t, p.ID)
	p.Close()
}

func TestNewPeerRandomID(t *testing.T) {
	peerServer, serverOpts := startServer()
	err := peerServer.Start()
	if err != nil {
		t.Logf("Server error: %s", err)
		t.FailNow()
	}
	defer peerServer.Stop()
	p, err := NewPeer("", getTestOpts(serverOpts))
	assert.NoError(t, err)
	assert.NotEmpty(t, p.ID)
	p.Close()
}

func TestNewPeerEvents(t *testing.T) {
	peerServer, serverOpts := startServer()
	err := peerServer.Start()
	if err != nil {
		t.Logf("Server error: %s", err)
		t.FailNow()
	}
	defer peerServer.Stop()
	p, err := NewPeer(rndName("test"), getTestOpts(serverOpts))
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

func TestDuplicatedID(t *testing.T) {

	peer1Name := rndName("duplicated")

	peerServer, serverOpts := startServer()
	err := peerServer.Start()
	if err != nil {
		t.Logf("Server error: %s", err)
		t.FailNow()
	}
	defer peerServer.Stop()

	peer1, err := NewPeer(peer1Name, getTestOpts(serverOpts))
	assert.NoError(t, err)
	defer peer1.Close()

	peer2, err := NewPeer(peer1Name, getTestOpts(serverOpts))
	assert.NoError(t, err)
	defer peer2.Close()

	peer2.On("error", func(raw interface{}) {
		err := raw.(error)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "taken")
	})

	<-time.After(time.Second * 1)
}

func TestHelloWorld(t *testing.T) {

	peer1Name := rndName("peer1")
	peer2Name := rndName("peer2")

	peerServer, serverOpts := startServer()
	err := peerServer.Start()
	if err != nil {
		t.Logf("Server error: %s", err)
		t.FailNow()
	}
	defer peerServer.Stop()

	peer1, err := NewPeer(peer1Name, getTestOpts(serverOpts))
	assert.NoError(t, err)
	defer peer1.Close()

	peer2, err := NewPeer(peer2Name, getTestOpts(serverOpts))
	assert.NoError(t, err)
	defer peer2.Close()

	// done := false
	done := false
	peer2.On("connection", func(data interface{}) {
		conn2 := data.(*DataConnection)
		conn2.On("data", func(data interface{}) {
			// Will print 'hi!'
			log.Println("Received")
			done = true
		})
	})

	conn1, err := peer1.Connect(peer2Name, nil)
	assert.NoError(t, err)
	conn1.On("open", func(data interface{}) {
		for {
			conn1.Send([]byte("hi!"), false)
			<-time.After(time.Millisecond * 1000)
		}
	})

	<-time.After(time.Second * 1)
	assert.True(t, done)
}

func TestLongPayload(t *testing.T) {

	peer1Name := rndName("peer1")
	peer2Name := rndName("peer2")

	peerServer, serverOpts := startServer()
	err := peerServer.Start()
	if err != nil {
		t.Logf("Server error: %s", err)
		t.FailNow()
	}
	defer peerServer.Stop()

	peer1, err := NewPeer(peer1Name, getTestOpts(serverOpts))
	assert.NoError(t, err)
	defer peer1.Close()

	peer2, err := NewPeer(peer2Name, getTestOpts(serverOpts))
	assert.NoError(t, err)
	defer peer2.Close()

	done := make(chan bool)
	peer2.On("connection", func(data interface{}) {
		conn2 := data.(*DataConnection)
		conn2.On("data", func(data interface{}) {
			log.Printf("Received\n")
			done <- true
		})
	})

	conn1, err := peer1.Connect(peer2Name, nil)
	assert.NoError(t, err)
	if err != nil {
		t.Fatal(err)
	}
	conn1.On("open", func(data interface{}) {
		raw := bytes.NewBuffer([]byte{})
		for {
			raw.Write([]byte("test"))
			if raw.Len() > 60000 {
				log.Printf("Msg size %d\n", raw.Len())
				break
			}
		}
		conn1.Send(raw.Bytes(), false)
	})

	<-done
}

func TestMediaCall(t *testing.T) {

	peer1Name := rndName("peer1")
	peer2Name := rndName("peer2")

	peerServer, serverOpts := startServer()
	err := peerServer.Start()
	if err != nil {
		t.Logf("Server error: %s", err)
		t.FailNow()
	}
	defer peerServer.Stop()

	peer1, err := NewPeer(peer1Name, getTestOpts(serverOpts))
	assert.NoError(t, err)
	defer peer1.Close()

	peer2, err := NewPeer(peer2Name, getTestOpts(serverOpts))
	assert.NoError(t, err)
	defer peer2.Close()

	track, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", "pion")
	if err != nil {
		panic(err)
	}

	call1, err := peer1.Call(peer2Name, track, nil)
	assert.NoError(t, err)

	peer2.On("call", func(raw interface{}) {
		// Answer the call, providing our mediaStream
		call := raw.(MediaConnection)
		var mediaStream webrtc.TrackLocal

		call.Answer(mediaStream, nil)
		call.On("stream", func(raw interface{}) {
			// stream := raw.(MediaStream)
			t.Log("peer2: Received remote stream")
		})
	})

	call1.On("stream", func(raw interface{}) {
		// stream := raw.(MediaStream)
		t.Log("peer1: Received remote stream")
	})

}
