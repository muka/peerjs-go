package main

import (
	"log"
	"time"

	"github.com/muka/peerjs-go"
)

func fail(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	// ensure to run your own peerjs-server
	// docker run --rm --name peerjs-server -p 9000:9000 -d peerjs/peerjs-server

	opts := peer.NewOptions()
	opts.Debug = 3
	opts.Path = "/myapp"
	opts.Host = "127.0.0.1"
	opts.Port = 9000
	opts.Secure = false

	peer1, err := peer.NewPeer("peer1", opts)
	fail(err)
	defer peer1.Close()

	peer2, err := peer.NewPeer("peer2", opts)
	fail(err)
	defer peer2.Close()

	peer2.On("connection", func(data interface{}) {
		conn2 := data.(*peer.DataConnection)
		conn2.On("data", func(data interface{}) {
			// Will print 'hi!'
			log.Printf("Received: %v\n", data)
		})
	})

	conn1, err := peer1.Connect("peer2", nil)
	fail(err)
	conn1.On("open", func(data interface{}) {
		for {
			conn1.Send([]byte("hi!"), false)
			<-time.After(time.Millisecond * 1000)
		}
	})

	select {}

}
