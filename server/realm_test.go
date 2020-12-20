package server

import (
	"testing"

	"github.com/muka/peer"
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
	c0 := peer.Message{}
	r.AddMessageToQueue("1", c0)
	c1 := r.GetClientByID("1")
	assert.Equal(t, c0, c1)
}

func TestRealmRandomID(t *testing.T) {
	r := NewRealm()
	assert.NotEqual(t, r.GenerateClientID(), r.GenerateClientID())
}
