package server

import (
	"github.com/muka/peerjs-go/models"
)

//Handler wrap a callback
type Handler func(client IClient, message models.IMessage) bool

// IHandlersRegistry interface for HandlersRegistry
type IHandlersRegistry interface {
	RegisterHandler(messageType string, handler Handler)
	Handle(client IClient, message models.IMessage) bool
}

//NewHandlersRegistry creates a new HandlersRegistry
func NewHandlersRegistry() IHandlersRegistry {
	h := &HandlersRegistry{
		handlers: make(map[string]Handler),
	}
	return h
}

// HandlersRegistry handlers registry
type HandlersRegistry struct {
	handlers map[string]Handler
}

// RegisterHandler register an handler
func (r *HandlersRegistry) RegisterHandler(messageType string, handler Handler) {
	if _, ok := r.handlers[messageType]; ok {
		return
	}
	r.handlers[messageType] = handler
}

//Handle handles a message
func (r *HandlersRegistry) Handle(client IClient, message models.IMessage) bool {
	handler, ok := r.handlers[message.GetType()]
	if !ok {
		return false
	}
	return handler(client, message)
}
