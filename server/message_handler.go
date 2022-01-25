package server

import (
	"github.com/muka/peerjs-go/models"
)

//IMessageHandler interface for MessageHandler
type IMessageHandler interface {
	Handle(client IClient, message models.IMessage) bool
}

//NewMessageHandler creates a new MessageHandler
func NewMessageHandler(realm IRealm, handlersRegistry IHandlersRegistry, opts Options) *MessageHandler {

	if handlersRegistry == nil {
		handlersRegistry = NewHandlersRegistry()
	}

	m := &MessageHandler{
		realm:            realm,
		handlersRegistry: handlersRegistry,
	}

	transmissionHandler := NewTransmissionHandler(realm, opts)
	heartbeatHandler := NewHeartbeatHandler(opts)

	handleHeartbeat := func(client IClient, message models.IMessage) bool {
		return heartbeatHandler(client, message)
	}

	handleTransmission := func(client IClient, message models.IMessage) bool {
		return transmissionHandler(client, message)
	}

	m.handlersRegistry.RegisterHandler(MessageTypeHeartbeat, handleHeartbeat)
	m.handlersRegistry.RegisterHandler(MessageTypeOffer, handleTransmission)
	m.handlersRegistry.RegisterHandler(MessageTypeAnswer, handleTransmission)
	m.handlersRegistry.RegisterHandler(MessageTypeCandidate, handleTransmission)
	m.handlersRegistry.RegisterHandler(MessageTypeLeave, handleTransmission)
	m.handlersRegistry.RegisterHandler(MessageTypeExpire, handleTransmission)

	return m
}

//MessageHandler wrap the message handler
type MessageHandler struct {
	realm            IRealm
	handlersRegistry IHandlersRegistry
}

//Handle handles a message
func (m *MessageHandler) Handle(client IClient, message models.IMessage) bool {
	return m.handlersRegistry.Handle(client, message)
}
