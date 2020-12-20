package server

import (
	"sync"

	"github.com/muka/peer"
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
	AddMessage(message peer.IMessage)
	ReadMessage() peer.IMessage
	GetMessages() []peer.IMessage
}

//MessageQueue type
type MessageQueue struct {
	lastReadAt int64
	messages   []peer.IMessage
	mMutex     sync.Mutex
}

//GetLastReadAt return last message read time
func (mq *MessageQueue) GetLastReadAt() int64 {
	return mq.lastReadAt
}

//AddMessage add message to queue
func (mq *MessageQueue) AddMessage(message peer.IMessage) {
	mq.mMutex.Lock()
	defer mq.mMutex.Unlock()
	mq.messages = append(mq.messages, message)
}

//ReadMessage read last message
func (mq *MessageQueue) ReadMessage() peer.IMessage {
	if len(mq.messages) > 0 {
		mq.lastReadAt = getTime()
		msg := mq.messages[0]
		mq.messages = mq.messages[1:]
		return msg
	}
	return nil
}

//GetMessages return all queued messages
func (mq *MessageQueue) GetMessages() []peer.IMessage {
	return mq.messages
}
