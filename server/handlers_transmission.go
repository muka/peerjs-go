package server

import (
	"errors"

	"github.com/muka/peerjs-go/models"
)

//NewTransmissionHandler handles transmission of messages
func NewTransmissionHandler(realm IRealm, opts Options) func(client IClient, message models.IMessage) bool {

	var handle func(client IClient, message models.IMessage) bool

	handle = func(client IClient, message models.IMessage) bool {

		clientID := "<no-id>"
		if client != nil {
			clientID = client.GetID()
		}

		log := createLogger("client:"+clientID, opts)

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

				handle(client, models.Message{
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
