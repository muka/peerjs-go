package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/muka/peerjs-go/enums"
	"github.com/muka/peerjs-go/models"
	"github.com/stretchr/testify/assert"
)

func TestHTTPServerGetID(t *testing.T) {
	opts := NewOptions()
	opts.Port = 64666
	opts.Host = "localhost"
	opts.AllowDiscovery = true

	getURL := func(path string) string {
		return fmt.Sprintf("http://%s:%d/%s%s", opts.Host, opts.Port, opts.Key, path)
	}

	realm := NewRealm()
	srv := NewHTTPServer(realm, nil, nil, opts)

	go srv.Start()
	defer srv.Stop()
	// wait for server to start
	<-time.After(time.Millisecond * 200)

	resp, err := http.Get(getURL("/id"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Get(getURL("/peers"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}
func TestHTTPServerNoDiscovery(t *testing.T) {
	opts := NewOptions()
	opts.Port = 64666
	opts.Host = "localhost"
	opts.AllowDiscovery = false

	getURL := func(path string) string {
		return fmt.Sprintf("http://%s:%d/%s%s", opts.Host, opts.Port, opts.Key, path)
	}

	realm := NewRealm()
	srv := NewHTTPServer(realm, nil, nil, opts)

	go srv.Start()
	defer srv.Stop()
	// wait for server to start
	<-time.After(time.Millisecond * 200)

	resp, err := http.Get(getURL("/peers"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

}
func TestHTTPServerExchange(t *testing.T) {
	opts := NewOptions()
	opts.Port = 64666
	opts.Host = "localhost"
	opts.AllowDiscovery = false

	id := "myid"
	token := "mytoken"

	getURL := func(path string) string {
		return fmt.Sprintf("http://%s:%d%s", opts.Host, opts.Port, path)
	}

	realm := NewRealm()
	auth := NewAuth(realm, opts)
	wss := NewWebSocketServer(realm, opts)
	srv := NewHTTPServer(realm, auth, wss, opts)

	go srv.Start()
	defer srv.Stop()
	// wait for server to start
	<-time.After(time.Millisecond * 200)

	msg := models.Message{
		Type:    enums.ServerMessageTypeOffer,
		Src:     "foo",
		Dst:     "bar",
		Payload: models.Payload{},
	}

	raw, err := json.Marshal(msg)
	assert.NoError(t, err)
	url := getURL(fmt.Sprintf("/%s/%s/%s/offer", opts.Key, id, token))

	// client not found
	resp, err := http.Post(url, "application/json", bytes.NewReader(raw))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	c := NewClient(id, token)
	srv.realm.SetClient(c, id)

	resp, err = http.Post(url, "application/json", bytes.NewReader(raw))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}
