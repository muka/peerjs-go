package server

import (
	"sync"

	"github.com/muka/peerjs-go/models"
	"github.com/muka/peerjs-go/util"
)

// ClientIDGenerator default hash generator
var ClientIDGenerator = func() string {
	return util.RandomToken()
}

// NewRealm creates a new Realm
func NewRealm() *Realm {
	r := &Realm{
		clients:       map[string]IClient{},
		messageQueues: map[string]IMessageQueue{},
		cMutex:        sync.Mutex{},
		mMutex:        sync.Mutex{},
	}
	return r
}

//IRealm interface for Realm
type IRealm interface {
	GetClientsIds() []string
	GetClientByID(clientID string) IClient
	GetClientsIdsWithQueue() []string
	SetClient(client IClient, id string)
	RemoveClientByID(id string) bool
	GetMessageQueueByID(id string) IMessageQueue
	AddMessageToQueue(id string, message models.IMessage)
	ClearMessageQueue(id string)
	GenerateClientID() string
}

// Realm implementation of a realm
type Realm struct {
	clients       map[string]IClient
	messageQueues map[string]IMessageQueue
	cMutex        sync.Mutex
	mMutex        sync.Mutex
}

//GetClientsIds return the list of client id
func (r *Realm) GetClientsIds() []string {
	keys := []string{}
	for key := range r.clients {
		keys = append(keys, key)
	}
	return keys
}

//GetClientByID return client by id
func (r *Realm) GetClientByID(clientID string) IClient {
	c, ok := r.clients[clientID]
	if !ok {
		return nil
	}
	return c
}

// GetClientsIdsWithQueue retur clients with queue
func (r *Realm) GetClientsIdsWithQueue() []string {
	keys := []string{}
	for key := range r.messageQueues {
		keys = append(keys, key)
	}
	return keys
}

// SetClient set a client
func (r *Realm) SetClient(client IClient, id string) {
	r.cMutex.Lock()
	defer r.cMutex.Unlock()
	r.clients[id] = client
}

//RemoveClientByID remove a client by id
func (r *Realm) RemoveClientByID(id string) bool {
	client := r.GetClientByID(id)
	if client == nil {
		return false
	}
	r.cMutex.Lock()
	defer r.cMutex.Unlock()
	delete(r.clients, id)
	return true
}

// GetMessageQueueByID get message by queue id
func (r *Realm) GetMessageQueueByID(id string) IMessageQueue {
	m, ok := r.messageQueues[id]
	if !ok {
		return nil
	}
	return m
}

// AddMessageToQueue add message to queue
func (r *Realm) AddMessageToQueue(id string, message models.IMessage) {
	if r.GetMessageQueueByID(id) == nil {
		r.mMutex.Lock()
		r.messageQueues[id] = NewMessageQueue()
		r.mMutex.Unlock()
	}

	m := r.GetMessageQueueByID(id)
	if m != nil {
		m.AddMessage(message)
	}
}

// ClearMessageQueue clear message queue
func (r *Realm) ClearMessageQueue(id string) {
	r.mMutex.Lock()
	defer r.mMutex.Unlock()
	delete(r.messageQueues, id)
}

// GenerateClientID generate a client id
func (r *Realm) GenerateClientID() string {
	return ClientIDGenerator()
}
