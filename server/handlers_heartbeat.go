package server

import (
	"github.com/muka/peer/models"
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
