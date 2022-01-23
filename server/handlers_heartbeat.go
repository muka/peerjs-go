package server

import (
	"github.com/KW-M/peerjs-go/models"
)

//NewHeartbeatHandler handles a heartbeat
func NewHeartbeatHandler(opts Options) func(client IClient, message models.IMessage) bool {
	return func(client IClient, message models.IMessage) bool {
		if client != nil {
			client.SetLastPing(getTime())
		}

		return true
	}
}
