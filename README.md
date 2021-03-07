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

	peer1, _ := NewPeer("peer1", NewOptions())
	defer peer1.Close()

	peer2, _ := NewPeer("peer2", NewOptions())
	defer peer2.Close()

	peer2.On("connection", func(data interface{}) {
		conn2 := data.(*DataConnection)
		conn2.On("data", func(data interface{}) {
			// Will print 'hi!'
			log.Printf("Received: %v\n", data)
		})
	})

	conn1, _ := peer1.Connect("peer2", nil)
	conn1.On("open", func(data interface{}) {
		for {
			conn1.Send([]byte("hi!"), false)
			<-time.After(time.Millisecond * 1000)
		}
	})

	select{}
```

### Peer server

A docker image for the GO based peer server is available at [opny/peer-server](https://hub.docker.com/r/opny/peer-server) built for Raspberry Pi and PCs

### Unsupported features

- Payload de/encoding based on [js-binarypack](https://github.com/peers/js-binarypack) is not supported.
- Message chunking (should be already done in recent browser versions)
