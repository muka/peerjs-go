package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
	opts     Options
	router   *mux.Router
	http     *http.Server
	handlers []func(http.HandlerFunc) http.HandlerFunc
}

// NewHTTPServer init a server
func NewHTTPServer(opts Options) *HTTPServer {
	r := mux.NewRouter()

	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s := &HTTPServer{
		opts:     opts,
		router:   r,
		http:     srv,
		handlers: []func(http.HandlerFunc) http.HandlerFunc{},
	}

	return s
}

func (h *HTTPServer) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h *HTTPServer) registerHandlers() {

	paths := []string{
		"/offer",
		"/candidate",
		"/answer",
		"/leave",
	}

	for _, p := range paths {
		h.router.Path(p).HandlerFunc(h.handler())
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
