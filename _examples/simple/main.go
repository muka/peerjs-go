package main

import (
	"log"

	"github.com/muka/peer"
)

func main() {

	peer1, err := peer.NewPeer("peer1", peer.NewOptions())
	if err != nil {
		log.Fatal(err)
	}

	peer2, err := peer.NewPeer("peer2", peer.NewOptions())
	if err != nil {
		log.Fatal(err)
	}

	conn1, err := peer1.Connect("peer2", peer.NewConnectionOptions())
	if err != nil {
		log.Fatal(err)
	}
	conn1.On("open", func(data interface{}) {
		conn1.Send([]byte("hi!"), false)
	})

	peer2.On("connection", func(data interface{}) {
		conn2 := data.(peer.DataConnection)
		conn2.On("data", func(data interface{}) {
			// Will print 'hi!'
			log.Printf("Received: %v\n", data)
		})
	})

}
