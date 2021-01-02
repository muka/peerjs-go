package server

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

//NewWebSocketServer create a new WebSocketServer
func NewWebSocketServer(opts Options) *WebSocketServer {
	ws := WebSocketServer{
		upgrader: websocket.Upgrader{},
		log:      createLogger("websocket-server", opts),
	}
	return &ws
}

// WebSocketServer wrap the websocket server
type WebSocketServer struct {
	upgrader websocket.Upgrader
	clients  []*websocket.Conn
	cMutex   sync.Mutex
	log      *logrus.Entry
}

// Send send data to the clients
func (wss *WebSocketServer) Send(data []byte) {
	for _, conn := range wss.clients {
		err := conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			wss.log.Warnf("Write failed: %s", err)
		}
	}
}

// Handler expose the http handler for websocket
func (wss *WebSocketServer) Handler() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := wss.upgrader.Upgrade(w, r, nil)
			if err != nil {
				wss.log.Warnf("upgrade error: %s", err)
				next.ServeHTTP(w, r)
				return
			}
			defer c.Close()
			for {
				mt, message, err := c.ReadMessage()
				if err != nil {
					wss.log.Warnf("read: %s", err)
					break
				}
				wss.log.Infof("recv: %s", message)
				err = c.WriteMessage(mt, message)
				if err != nil {
					wss.log.Warnf("write: %s", err)
					break
				}
			}
		})
	}
}
