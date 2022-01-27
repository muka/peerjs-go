# Golang PeerJS

A Golang port of [PeerJS](https://github.com/peers/peerjs)

## Implementation

- [X] Datachannel
- [X] Mediachannel
- [X] Test coverage > 80%
- [X] Signalling server
- [ ] Interoperability tests

## Usage example

See [_examples folder](./_examples)

```golang

package main

import (
	"log"
	"time"

	peer "github.com/muka/peerjs-go"
)

func main() {
	peer1, _ := peer.NewPeer("peer1", peer.NewOptions())
	defer peer1.Close()

	peer2, _ := peer.NewPeer("peer2", peer.NewOptions())
	defer peer2.Close()

	peer2.On("connection", func(data interface{}) {
		conn2 := data.(*peer.DataConnection)
		conn2.On("data", func(data interface{}) {
			// Will print 'hi!'
			log.Printf("Received: %#v: %s\n", data, data)
		})
	})

	conn1, _ := peer1.Connect("peer2", nil)
	conn1.On("open", func(data interface{}) {
		for {
			conn1.Send([]byte("hi!"), false)
			<-time.After(time.Millisecond * 1000)
		}
	})

	select {}
}
```


Further documentation can be found at: [https://pkg.go.dev/github.com/muka/peerjs-go/](https://pkg.go.dev/github.com/muka/peerjs-go/)

### Peer server

A docker image for the GO based peer server is available at [opny/peer-server](https://hub.docker.com/r/opny/peer-server) built for Raspberry Pi and PCs

The source of the GO based peer server can be found in the [server folder](./server/) which can be imported as the go package `"github.com/muka/peerjs-go/server"`
and is documentated on [pkg.go.dev](https://pkg.go.dev/github.com/muka/peerjs-go/server). Also see example usage in the [peer_test.go](peer_test.go) file.

If you want a standalone GO based Peerjs server, run `go build ./cmd/server/main.go` to get an exacutable. To set the server options, create a `peer.yaml` configuration file in the same folder as the executable.
__Available Server Options:__
- __Host__ String
- __Port__ Int
- __LogLevel__ String
- __ExpireTimeout__ Int64
- __AliveTimeout__ Int64
- __Key__ String
- __Path__ String
- __ConcurrentLimit__ Int
- __AllowDiscovery__ Bool

### Unsupported features

- Payload de/encoding based on [js-binarypack](https://github.com/peers/js-binarypack) is not supported.
- Message chunking (should be already done in recent browser versions)
