package server

//New creates a new PeerServer
func New(opts Options) *PeerServer {

	s := new(PeerServer)

	s.realm = NewRealm()
	s.auth = NewAuth(s.realm, opts)
	s.wss = NewWebSocket(opts)

	s.http = NewHTTPServer(opts)

	s.http.AddHandlers(
		s.wss.Handler(),
		s.auth.Handler(),
	)

	return s
}

//PeerServer wrap the peer server functionalities
type PeerServer struct {
	http  *HTTPServer
	realm IRealm
	auth  *Auth
	wss   *WebSocket
}

// Stop stops the peer server
func (p *PeerServer) Stop() error {
	p.http.Stop()
	return nil
}

// Start start the peer server
func (p *PeerServer) Start() error {

	errEv := make(chan error)

	go func() {
		err := p.http.Start()
		if err != nil {
			errEv <- err
		}
	}()

	err := <-errEv
	return err
}
