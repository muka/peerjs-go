package server

import "github.com/muka/peer"

//New creates a new PeerServer
func New(opts Options) *PeerServer {

	s := new(PeerServer)
	s.Emitter = peer.NewEmitter()

	s.realm = NewRealm()
	s.auth = NewAuth(s.realm, opts)
	s.wss = NewWebSocketServer(s.realm, opts)

	s.http = NewHTTPServer(s.realm, opts)

	s.http.AddHandlers(
		s.wss.Handler(),
		s.auth.Handler(),
	)

	s.checkBrokenConnections = NewCheckBrokenConnections(
		s.realm,
		opts,
		func(client IClient) {
			s.Emit("disconnect", client)
		},
	)

	s.messageExpire = NewMessagesExpire(s.realm, opts, s.http.messageHandler)

	s.initialize()

	return s
}

//PeerServer wrap the peer server functionalities
type PeerServer struct {
	peer.Emitter
	http                   *HTTPServer
	realm                  IRealm
	auth                   *Auth
	wss                    *WebSocketServer
	checkBrokenConnections *CheckBrokenConnections
	messageExpire          IMessagesExpire
}

func (p *PeerServer) initialize() {

	p.wss.On("connection", func(data interface{}) {
		client := data.(IClient)
		mq := p.realm.GetMessageQueueByID(client.GetID())
		if mq != nil {
			for {
				message := mq.ReadMessage()
				if message == nil {
					break
				}
				p.http.messageHandler.Handle(client, message)
			}
			p.realm.ClearMessageQueue(client.GetID())
		}
		p.Emit("connection", client)
	})

	p.wss.On("message", func(data interface{}) {
		cm := data.(ClientMessage)
		p.Emit("message", cm)
		p.http.messageHandler.Handle(cm.Client, cm.Message)
	})

	p.wss.On("close", func(data interface{}) {
		client := data.(IClient)
		p.Emit("disconnect", client)
	})

	p.wss.On("error", func(data interface{}) {
		err := data.(error)
		p.Emit("error", err)
	})

	p.messageExpire.Start()
	p.checkBrokenConnections.Start()
}

// Stop stops the peer server
func (p *PeerServer) Stop() error {
	p.http.Stop()
	p.messageExpire.Stop()
	p.checkBrokenConnections.Stop()
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
