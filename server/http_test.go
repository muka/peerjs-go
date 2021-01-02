package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/muka/peer"
	"github.com/stretchr/testify/assert"
)

func TestHTTPServerGetID(t *testing.T) {
	opts := NewOptions()
	opts.Port = 64666
	opts.Host = "localhost"
	opts.AllowDiscovery = true

	getURL := func(path string) string {
		return fmt.Sprintf("http://%s:%d%s", opts.Host, opts.Port, path)
	}

	realm := NewRealm()
	srv := NewHTTPServer(realm, opts)

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
		return fmt.Sprintf("http://%s:%d%s", opts.Host, opts.Port, path)
	}

	realm := NewRealm()
	srv := NewHTTPServer(realm, opts)

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

	getURL := func(path string) string {
		return fmt.Sprintf("http://%s:%d%s", opts.Host, opts.Port, path)
	}

	realm := NewRealm()
	srv := NewHTTPServer(realm, opts)

	go srv.Start()
	defer srv.Stop()
	// wait for server to start
	<-time.After(time.Millisecond * 200)

	msg := peer.Message{
		Type:    peer.ServerMessageTypeOffer,
		Src:     "foo",
		Dst:     "bar",
		Payload: peer.Payload{},
	}

	raw, err := json.Marshal(msg)
	assert.NoError(t, err)
	resp, err := http.Post(getURL("/offer"), "application/json", bytes.NewReader(raw))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = http.Post(getURL(fmt.Sprintf("/offer/")), "application/json", bytes.NewReader(raw))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}
