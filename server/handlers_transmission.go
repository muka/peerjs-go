package server

import (
	"encoding/json"
	"errors"

	"github.com/muka/peer"
)

//NewTransmissionHandler handles transmission of messages
func NewTransmissionHandler(realm IRealm, opts Options) func(client IClient, message peer.IMessage) bool {

	var handle func(client IClient, message peer.IMessage) bool

	handle = func(client IClient, message peer.IMessage) bool {

		log := createLogger("client:"+client.GetID(), opts)

		mtype := message.GetType()
		srcId := message.GetSrc()
		dstId := message.GetDst()

		destinationClient := realm.GetClientByID(dstId)

		// User is connected!
		if destinationClient != nil {
			socket := destinationClient.GetSocket()
			var err error
			var data []byte
			if socket != nil {
				data, err = json.Marshal(message)
				if err != nil {
					err = socket.Send(data)
				}
			} else {
				err = errors.New("Peer dead")
			}

			if err != nil {
				// This happens when a peer disconnects without closing connections and
				// the associated WebSocket has not closed.
				// Tell other side to stop trying.
				log.Warnf("Error: %s", err)
				if socket != nil {
					socket.Close()
				} else {
					realm.RemoveClientByID(destinationClient.GetID())
				}

				handle(client, peer.Message{
					Type: MessageTypeLeave,
					Src:  dstId,
					Dst:  srcId,
				})
			}

		} else {
			// Wait for this client to connect/reconnect (XHR) for important
			// messages.
			if (mtype != MessageTypeLeave && mtype != MessageTypeExpire) && dstId != "" {
				realm.AddMessageToQueue(dstId, message)
			} else if mtype == MessageTypeLeave && dstId == "" {
				realm.RemoveClientByID(srcId)
			} else {
				// Unavailable destination specified with message LEAVE or EXPIRE
				// Ignore
			}
		}

		return true
	}

	return handle
}
