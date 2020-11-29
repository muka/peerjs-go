# Golang PeerJS 

A Golang implementation of [PeerJS](https://github.com/peers/peerjs)

This project is in __early stage development__!

## Implementation notes

- [X] 
- [X] Datachannel
- [ ] Mediachannel
- [ ] Test coverage > 80%

### Unsupported features:

- Message chunking (implemented but untested)
- Payload de/encoding based on [js-binarypack](https://github.com/peers/js-binarypack) is not supported. An incomplete attempt in branch `binarypack`

## Usage example

See [_examples folder](./_examples)

```golang

	peer1, _ := NewPeer("peer1", getTestOpts())
	defer peer1.Close()

	peer2, _ := NewPeer("peer2", getTestOpts())
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