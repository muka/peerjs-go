package server

import (
	"time"

	"github.com/sirupsen/logrus"
)

const DefaultCheckInterval = 300

//NewCheckBrokenConnections create a new CheckBrokenConnections
func NewCheckBrokenConnections(realm IRealm, opts Options, onClose func(client IClient)) *CheckBrokenConnections {
	if onClose == nil {
		onClose = func(client IClient) {}
	}
	return &CheckBrokenConnections{
		realm:   realm,
		opts:    opts,
		onClose: onClose,
		log:     createLogger("checkBrokenConnections", opts),
		close:   make(chan bool, 1),
	}
}

//CheckBrokenConnections watch for broken connections
type CheckBrokenConnections struct {
	realm   IRealm
	opts    Options
	onClose func(IClient)
	ticker  *time.Ticker
	log     *logrus.Entry
	close   chan bool
}

func (b *CheckBrokenConnections) checkConnections() {

	clientsIds := b.realm.GetClientsIds()
	now := getTime()
	aliveTimeout := b.opts.AliveTimeout

	for _, clientID := range clientsIds {

		client := b.realm.GetClientByID(clientID)
		if client == nil {
			continue
		}

		timeSinceLastPing := now - client.GetLastPing()

		if timeSinceLastPing < aliveTimeout {
			continue
		}

		socket := client.GetSocket()
		if socket != nil {
			b.log.Infof("Closing broken connection clientID=%s", clientID)
			err := socket.Close()
			if err != nil {
				b.log.Warnf("Failed to close socket: %s", err)
			}
		}
		b.realm.ClearMessageQueue(clientID)
		b.realm.RemoveClientByID(clientID)
		client.SetSocket(nil)
		b.onClose(client)
	}
}

//Stop close the connection checker
func (b *CheckBrokenConnections) Stop() {
	if b.ticker == nil {
		return
	}
	b.close <- true
}

//Start initialize the connection checker
func (b *CheckBrokenConnections) Start() {

	b.ticker = time.NewTicker(DefaultCheckInterval * time.Millisecond)

	go func() {
		for {
			select {
			case <-b.close:
				b.ticker.Stop()
				b.ticker = nil
				return
			case <-b.ticker.C:
				b.checkConnections()
			}
		}
	}()

}
