package server

import "github.com/muka/peer"

//NewHeartbeatHandler handles a heartbeat
func NewHeartbeatHandler(opts Options) func(client IClient, message peer.IMessage) bool {
	return func(client IClient, message peer.IMessage) bool {
		if client != nil {
			client.SetLastPing(getTime())
		}

		return true
	}
}
