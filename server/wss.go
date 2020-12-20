package server

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

//NewWebSocket create a new WebSocket
func NewWebSocket(opts Options) *WebSocket {
	ws := WebSocket{
		upgrader: websocket.Upgrader{},
		log:      createLogger("websocket", opts),
	}
	return &ws
}

// WebSocket wrap the websocket server
type WebSocket struct {
	upgrader websocket.Upgrader
	clients  []*websocket.Conn
	cMutex   sync.Mutex
	log      *logrus.Entry
}

// Send send data to the clients
func (wss *WebSocket) Send(data []byte) {
	for _, conn := range wss.clients {
		err := conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			wss.log.Warnf("Write failed: %s", err)
		}
	}
}

// Handler expose the http handler for websocket
func (wss *WebSocket) Handler() mux.MiddlewareFunc {
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
