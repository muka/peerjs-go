package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/muka/peer"
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
	ExpireTimeout   int
	AliveTimeout    int
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
}

// NewHTTPServer init a server
func NewHTTPServer(realm IRealm, opts Options) *HTTPServer {

	r := mux.NewRouter()

	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s := &HTTPServer{
		opts:           opts,
		realm:          realm,
		log:            createLogger("http", opts),
		router:         r,
		http:           srv,
		handlers:       []func(http.HandlerFunc) http.HandlerFunc{},
		messageHandler: NewMessageHandler(realm, nil, opts),
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

		payload := new(peer.Message)
		err = json.Unmarshal(body, payload)
		if err != nil {
			http.Error(w, "Failed to decode message", http.StatusInternalServerError)
			return
		}

		message := peer.Message{
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

func (h *HTTPServer) registerHandlers() {

	baseRoute := h.router.PathPrefix(h.opts.Path).Subrouter()

	// public API
	baseRoute.
		HandleFunc("/{key}/id", func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Add("content-type", "text/html")
			rw.Write([]byte(h.realm.GenerateClientID()))
		}).
		Methods("GET")

	baseRoute.
		HandleFunc("/{key}/peers", func(rw http.ResponseWriter, r *http.Request) {
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
		}).
		Methods("GET")

	paths := []string{
		"offer",
		"candidate",
		"answer",
		"leave",
	}

	for _, p := range paths {
		endpoint := fmt.Sprintf("/{key}/{id}/{token}/%s", p)
		baseRoute.
			HandleFunc(endpoint, h.handler()).
			Methods("POST")
	}

}

//AddHandlers register HTTP handlers
func (h *HTTPServer) AddHandlers(middlewares ...mux.MiddlewareFunc) {
	h.router.Use(middlewares...)
}

//Start start the HTTP server
func (h *HTTPServer) Start() error {
	h.registerHandlers()
	return h.http.ListenAndServe()
}

//Stop stops the HTTP server
func (h *HTTPServer) Stop() error {
	return h.http.Close()
}
