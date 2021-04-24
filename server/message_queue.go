package server

import (
	"sync"

	"github.com/muka/peerjs-go/models"
)

// NewMessageQueue creates a new MessageQueue
func NewMessageQueue() *MessageQueue {
	mq := new(MessageQueue)
	mq.lastReadAt = getTime()
	return mq
}

// IMessageQueue message queue interface
type IMessageQueue interface {
	GetLastReadAt() int64
	AddMessage(message models.IMessage)
	ReadMessage() models.IMessage
	GetMessages() []models.IMessage
}

//MessageQueue type
type MessageQueue struct {
	lastReadAt int64
	messages   []models.IMessage
	mMutex     sync.Mutex
}

//GetLastReadAt return last message read time
func (mq *MessageQueue) GetLastReadAt() int64 {
	return mq.lastReadAt
}

//AddMessage add message to queue
func (mq *MessageQueue) AddMessage(message models.IMessage) {
	mq.mMutex.Lock()
	defer mq.mMutex.Unlock()
	mq.messages = append(mq.messages, message)
}

//ReadMessage read last message
func (mq *MessageQueue) ReadMessage() models.IMessage {
	if len(mq.messages) > 0 {
		mq.mMutex.Lock()
		defer mq.mMutex.Unlock()
		mq.lastReadAt = getTime()
		msg := mq.messages[0]
		mq.messages = mq.messages[1:]
		return msg
	}
	return nil
}

//GetMessages return all queued messages
func (mq *MessageQueue) GetMessages() []models.IMessage {
	return mq.messages
}
