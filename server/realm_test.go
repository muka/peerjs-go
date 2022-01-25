package server

import (
	"testing"

	"github.com/muka/peerjs-go/models"
	"github.com/stretchr/testify/assert"
)

type xMessage struct{}

func TestRealmClients(t *testing.T) {
	r := NewRealm()
	c0 := NewClient("1", "test")
	r.SetClient(c0, "1")
	c1 := r.GetClientByID("1")
	assert.Equal(t, c0, c1)
}

func TestRealmMessage(t *testing.T) {
	r := NewRealm()
	c0 := NewClient("1", "a")
	r.SetClient(c0, c0.GetID())
	m0 := models.Message{}
	r.AddMessageToQueue("1", m0)
	c1 := r.GetClientByID("1")
	assert.Equal(t, c0, c1)
}

func TestRealmRandomID(t *testing.T) {
	r := NewRealm()
	assert.NotEqual(t, r.GenerateClientID(), r.GenerateClientID())
}

func TestRealmQueueByID(t *testing.T) {
	r := NewRealm()
	c0 := NewClient("1", "a")
	r.SetClient(c0, c0.GetID())
	m0 := models.Message{}
	r.AddMessageToQueue("1", m0)

	c1 := r.GetClientsIds()
	assert.Equal(t, len(c1), 1)

	c2 := r.GetClientsIdsWithQueue()
	assert.Equal(t, len(c2), 1)

	q := r.GetMessageQueueByID("1")
	assert.Equal(t, len(q.GetMessages()), 1)

	r.ClearMessageQueue("1")
	q = r.GetMessageQueueByID("1")
	assert.Nil(t, q)

	m1 := models.Message{}
	r.AddMessageToQueue("1", m1)
	q = r.GetMessageQueueByID("1")
	assert.NotNil(t, q)
	assert.NotZero(t, q.GetLastReadAt())
	m2 := q.ReadMessage()
	assert.Equal(t, m1, m2)
	assert.NotZero(t, q.GetLastReadAt())
	assert.Empty(t, q.ReadMessage())

}
func TestRealmRemove(t *testing.T) {
	r := NewRealm()
	c0 := NewClient("1", "a")
	r.SetClient(c0, c0.GetID())

	r.RemoveClientByID("1")
	r.RemoveClientByID("1")

	assert.Equal(t, len(r.GetClientsIds()), 0)
}
