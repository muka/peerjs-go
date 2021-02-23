package server

import (
	"errors"

	"github.com/muka/peer"
)

//NewTransmissionHandler handles transmission of messages
func NewTransmissionHandler(realm IRealm, opts Options) func(client IClient, message peer.IMessage) bool {

	var handle func(client IClient, message peer.IMessage) bool

	handle = func(client IClient, message peer.IMessage) bool {

		log := createLogger("client:"+client.GetID(), opts)

		mtype := message.GetType()
		srcID := message.GetSrc()
		dstID := message.GetDst()

		destinationClient := realm.GetClientByID(dstID)

		// User is connected!
		if destinationClient != nil {
			socket := destinationClient.GetSocket()
			var err error
			if socket != nil {
				err = socket.WriteJSON(message)
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
					Src:  dstID,
					Dst:  srcID,
				})
			}

		} else {
			// Wait for this client to connect/reconnect (XHR) for important
			// messages.
			if (mtype != MessageTypeLeave && mtype != MessageTypeExpire) && dstID != "" {
				realm.AddMessageToQueue(dstID, message)
			} else if mtype == MessageTypeLeave && dstID == "" {
				realm.RemoveClientByID(srcID)
			} else {
				// Unavailable destination specified with message LEAVE or EXPIRE
				// Ignore
			}
		}

		return true
	}

	return handle
}
