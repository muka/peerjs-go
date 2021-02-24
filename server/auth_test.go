package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuth(t *testing.T) {
	a := NewAuth(NewRealm(), NewOptions())
	assert.NotNil(t, a)
}

func testGetAuth() *Auth {
	o := NewOptions()
	r := NewRealm()
	a := NewAuth(r, o)
	return a
}

func TestAuthHandlerWrongKey(t *testing.T) {

	a := testGetAuth()
	id := a.realm.GenerateClientID()
	token := a.realm.GenerateClientID()

	// wrong key
	err := a.checkRequest("wrong key", id, token)
	assert.Equal(t, err, errInvalidKey)
}
func TestAuthHandlerClientValid(t *testing.T) {

	a := testGetAuth()
	id := a.realm.GenerateClientID()
	token := a.realm.GenerateClientID()

	var err error

	// valid client
	client := NewClient(id, token)
	a.realm.SetClient(client, id)
	err = a.checkRequest(a.opts.Key, id, token)
	assert.Nil(t, err)
}

func TestAuthHandlerClientIDEmpty(t *testing.T) {

	a := testGetAuth()
	token := a.realm.GenerateClientID()

	var err error

	err = a.checkRequest(a.opts.Key, "", token)
	assert.Equal(t, err, errUnauthorized)

}
func TestAuthHandlerClientTokenEmpty(t *testing.T) {

	a := testGetAuth()
	id := a.realm.GenerateClientID()

	var err error

	err = a.checkRequest(a.opts.Key, id, "")
	assert.Equal(t, err, errUnauthorized)

}
func TestAuthHandlerClientTokenInvalid(t *testing.T) {

	a := testGetAuth()
	id := a.realm.GenerateClientID()
	token := a.realm.GenerateClientID()

	var err error

	client := NewClient(id, token)
	a.realm.SetClient(client, id)
	err = a.checkRequest(a.opts.Key, id, token)

	err = a.checkRequest(a.opts.Key, id, "wrong")
	assert.Equal(t, err, errUnauthorized)

}
func TestAuthHandlerClientRemoved(t *testing.T) {

	a := testGetAuth()
	id := a.realm.GenerateClientID()
	token := a.realm.GenerateClientID()

	var err error

	// valid client
	client := NewClient(id, token)
	a.realm.SetClient(client, id)
	err = a.checkRequest(a.opts.Key, id, token)
	assert.Nil(t, err)

	a.realm.RemoveClientByID(id)
	assert.Nil(t, a.realm.GetClientByID(id))

	err = a.checkRequest(a.opts.Key, id, token)
	assert.Error(t, err)

}
