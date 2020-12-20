package server

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestNewAuth(t *testing.T) {
	a := NewAuth(NewRealm(), NewOptions())
	assert.NotNil(t, a)
}

func TestAuthHandler(t *testing.T) {

	o := NewOptions()
	r := NewRealm()
	a := NewAuth(r, o)

	id := xid.New().String()
	token := xid.New().String()

	var err error

	// wrong key
	err = a.checkRequest("wrong key", id, token)
	assert.Equal(t, err, errInvalidKey)

	// valid client
	client := NewClient(id, token)
	r.SetClient(client, id)
	err = a.checkRequest(o.Key, id, token)
	assert.Nil(t, err)

	err = a.checkRequest(o.Key, "", token)
	assert.Equal(t, err, errUnauthorized)

	err = a.checkRequest(o.Key, id, "")
	assert.Equal(t, err, errUnauthorized)

}
