package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/muka/peerjs-go/models"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

//NewOptions create default options
func NewOptions() Options {
	return Options{
		Port:            9000,
		Host:            "0.0.0.0",
		LogLevel:        "info",
		ExpireTimeout:   5000,
		AliveTimeout:    60000,
		Key:             "peerjs",
		Path:            "/",
		ConcurrentLimit: 5000,
		AllowDiscovery:  false,
		CleanupOutMsgs:  1000,
	}
}

//Options peer server options
type Options struct {
	Port            int
	Host            string
	LogLevel        string
	ExpireTimeout   int64
	AliveTimeout    int64
	Key             string
	Path            string
	ConcurrentLimit int
	AllowDiscovery  bool
	CleanupOutMsgs  int
}

//HTTPServer peer server
type HTTPServer struct {
	opts           Options
	realm          IRealm
	log            *logrus.Entry
	messageHandler *MessageHandler
	router         *mux.Router
	http           *http.Server
	handlers       []func(http.HandlerFunc) http.HandlerFunc
	auth           *Auth
	wss            *WebSocketServer
}

// NewHTTPServer init a server
func NewHTTPServer(realm IRealm, auth *Auth, wss *WebSocketServer, opts Options) *HTTPServer {

	r := mux.NewRouter()

	s := &HTTPServer{
		opts:   opts,
		realm:  realm,
		log:    createLogger("http", opts),
		router: r,
		// http:           srv,
		handlers:       []func(http.HandlerFunc) http.HandlerFunc{},
		messageHandler: NewMessageHandler(realm, nil, opts),
		auth:           auth,
		wss:            wss,
	}

	return s
}

func (h *HTTPServer) handler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Error(w, "Missing client id", http.StatusBadRequest)
			return
		}

		client := h.realm.GetClientByID(id)
		if client == nil {
			http.Error(w, fmt.Sprintf("Client %s not found", id), http.StatusNotFound)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
			return
		}

		payload := new(models.Message)
		err = json.Unmarshal(body, payload)
		if err != nil {
			http.Error(w, "Failed to decode message", http.StatusInternalServerError)
			return
		}

		message := models.Message{
			Type:    payload.Type,
			Src:     id,
			Dst:     payload.Dst,
			Payload: payload.Payload,
		}

		h.messageHandler.Handle(client, message)

		w.WriteHeader(200)
		w.Write([]byte{})
	})
}

func (h *HTTPServer) peersHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if !h.opts.AllowDiscovery {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte{})
			return
		}

		rw.Header().Add("content-type", "application/json")
		raw, err := json.Marshal(h.realm.GetClientsIds())
		if err != nil {
			h.log.Warnf("/peers: Marshal error %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte{})
			return
		}
		rw.Write(raw)
	}
}

func (h *HTTPServer) registerHandlers() error {

	baseRoute := h.router.PathPrefix(h.opts.Path).Subrouter()
	h.log.Debugf("Path prefix: %s", h.opts.Path)

	err := baseRoute.
		HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Add("content-type", "application/json")
			rw.Write([]byte(`{
  "name": "PeerJS Server",
  "description": "A server side element to broker connections between PeerJS clients.",
  "website": "https://github.com/muka/peerjs-go/tree/main/server"
}`))
		}).
		Methods("GET").GetError()
	if err != nil {
		return err
	}

	// public API
	err = baseRoute.
		HandleFunc("/{key}/id", func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Add("content-type", "text/html")
			rw.Write([]byte(h.realm.GenerateClientID()))
		}).
		Methods("GET").GetError()
	if err != nil {
		return err
	}

	err = baseRoute.
		Path("/{key}/peers").
		Handler(h.auth.HTTPHandler(h.peersHandler())).
		Methods("GET").GetError()
	if err != nil {
		return err
	}

	// handle WS route
	err = baseRoute.
		Path(fmt.Sprintf("/%s", h.opts.Key)).
		Handler(h.auth.WSHandler(h.wss.Handler())).
		Methods("GET").GetError()
	if err != nil {
		return err
	}

	paths := []string{
		"offer",
		"candidate",
		"answer",
		"leave",
	}

	for _, p := range paths {
		endpoint := fmt.Sprintf("/{key}/{id}/{token}/%s", p)
		err := baseRoute.
			Path(endpoint).
			Handler(h.auth.HTTPHandler(h.handler())).
			Methods("POST").GetError()
		if err != nil {
			return err
		}
	}

	return nil
}

//Start start the HTTP server
func (h *HTTPServer) Start() error {

	err := h.registerHandlers()
	if err != nil {
		return err
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(h.router)

	h.http = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", h.opts.Host, h.opts.Port),
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return h.http.ListenAndServe()
}

//Stop stops the HTTP server
func (h *HTTPServer) Stop() error {
	return h.http.Close()
}
