# Golang PeerJS

A Golang port of [PeerJS](https://github.com/peers/peerjs)

## Implementation

- [x] Datachannel
- [x] Mediachannel
- [x] Test coverage > 80%
- [x] Signalling server
- [ ] Interoperability tests

## Docs

[![Go Reference](https://pkg.go.dev/badge/github.com/muka/peerjs-go.svg)](https://pkg.go.dev/github.com/muka/peerjs-go)

⚠️ _Note_: While the Javascript [PeerJS documentation](https://peerjs.com/docs/) often applies to this library, there are differences, namely:

- All methods and properties are in PascalCase.
- Enum values are represented as seperate constants.
- All peer event callback functions should take a generic interface{} parameter, and then cast the interface{} to the appropriate peerjs-go type.
- Blocked peer event callback functions will block other peerjs-go events from firing.
- Refer to the [go package docs](https://pkg.go.dev/github.com/muka/peerjs-go) whenever unsure.

### Unsupported features

- Payload de/encoding based on [js-binarypack](https://github.com/peers/js-binarypack) is not supported.
- Message chunking (should be already done in recent browser versions)

## Usage example

See [\_examples folder](./_examples)

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

## Peer server

This library includes a GO based peer server in the [/server folder](./server/)

### Documentation:

[![Go Reference](https://pkg.go.dev/badge/github.com/muka/peerjs-go/server.svg)](https://pkg.go.dev/github.com/muka/peerjs-go/server)

### Example usage

```golang
package main

import (
	"log"

	peerjsServer "github.com/muka/peerjs-go/server"
)

func main() {
	serverOptions := peerjsServer.NewOptions()
	// These are the default values NewOptions() creates:
	serverOptions.Port = 9000
	serverOptions.Host = "0.0.0.0"
	serverOptions.LogLevel = "info"
	serverOptions.ExpireTimeout = 5000
	serverOptions.AliveTimeout = 60000
	serverOptions.Key = "peerjs"
	serverOptions.Path = "/"
	serverOptions.ConcurrentLimit = 5000
	serverOptions.AllowDiscovery = false
	serverOptions.CleanupOutMsgs = 1000

	server := peerjsServer.New(serverOptions)
	defer server.Stop()

	if err := server.Start(); err != nil {
		log.Printf("Error starting peerjs server: %s", err)
	}

	select{}
}
```

### Docker

A docker image for the GO based peer server is available at [opny/peer-server](https://hub.docker.com/r/opny/peer-server) built for Raspberry Pi and PCs

### Standalone

To build a standalone GO based Peerjs server executable, run `go build ./cmd/server/main.go` in the repository folder. To set the server options, create a `peer.yaml` configuration file in the same folder as the executable with the following options:

**Available Server Options:**

- **Host** String
- **Port** Int
- **LogLevel** String
- **ExpireTimeout** Int64
- **AliveTimeout** Int64
- **Key** String
- **Path** String
- **ConcurrentLimit** Int
- **AllowDiscovery** Bool
- **CleanupOutMsgs** Int
