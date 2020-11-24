# Golang PeerJS 

A golang implementation of [PeerJS](https://github.com/peers/peerjs)

This project is in __early stage development__!

## Implementation status

- [X] Datachannel
- [ ] Mediachannel
- [ ] Test coverage > 80%

## Usage example

```golang


	peer1, err := NewPeer("peer1", getTestOpts())
	assert.NoError(t, err)
	defer peer1.Close()

	peer2, err := NewPeer("peer2", getTestOpts())
	assert.NoError(t, err)
	defer peer2.Close()

	peer2.On("connection", func(data interface{}) {
		conn2 := data.(*DataConnection)
		conn2.On("data", func(data interface{}) {
			// Will print 'hi!'
			log.Printf("Received: %v\n", data)
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

	select{}
```