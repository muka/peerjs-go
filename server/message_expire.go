package server

import (
	"fmt"
	"time"

	"github.com/muka/peerjs-go/models"
	"github.com/sirupsen/logrus"
)

//IMessagesExpire MessagesExpire interface
type IMessagesExpire interface {
	Start()
	Stop()
}

func NewMessagesExpire(realm IRealm, opts Options, messageHandler IMessageHandler) *MessagesExpire {
	return &MessagesExpire{
		realm:          realm,
		opts:           opts,
		messageHandler: messageHandler,
		log:            createLogger("messageExpire", opts),
		close:          make(chan bool, 1),
	}
}

//MessagesExpire check for expired messages
type MessagesExpire struct {
	realm          IRealm
	opts           Options
	messageHandler IMessageHandler
	ticker         *time.Ticker
	log            *logrus.Entry
	close          chan bool
}

func (b *MessagesExpire) pruneOutstanding() {
	destinationClientsIds := b.realm.GetClientsIdsWithQueue()

	now := getTime()
	maxDiff := b.opts.ExpireTimeout

	seen := map[string]bool{}

	for _, destinationClientID := range destinationClientsIds {

		messageQueue := b.realm.GetMessageQueueByID(destinationClientID)
		if messageQueue == nil {
			continue
		}

		lastReadDiff := now - messageQueue.GetLastReadAt()
		if lastReadDiff < maxDiff {
			continue
		}

		messages := messageQueue.GetMessages()
		for _, message := range messages {
			seenKey := fmt.Sprintf("%s_%s", message.GetSrc(), message.GetDst())

			if _, ok := seen[seenKey]; !ok {
				b.messageHandler.Handle(nil, models.Message{
					Type: MessageTypeExpire,
					Src:  message.GetDst(),
					Dst:  message.GetSrc(),
				})

				seen[seenKey] = true
			}
		}

		b.realm.ClearMessageQueue(destinationClientID)
	}
}

//Start the message expire check
func (b *MessagesExpire) Start() {

	b.ticker = time.NewTicker(DefaultCheckInterval * time.Millisecond)

	go func() {
		for {
			select {
			case <-b.close:
				b.ticker.Stop()
				b.ticker = nil
				return
			case <-b.ticker.C:
				b.pruneOutstanding()
			}
		}
	}()

}

//Stop the message expire check
func (b *MessagesExpire) Stop() {
	if b.ticker == nil {
		return
	}
	b.close <- true
}
