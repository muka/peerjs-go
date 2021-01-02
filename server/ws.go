package server

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

//NewWebSocket create a new WebSocket
func NewWebSocket(opts Options) *WebSocket {
	ws := WebSocket{
		log:    createLogger("websocket", opts),
		cMutex: sync.Mutex{},
	}
	return &ws
}

// WebSocket wrap the websocket server
type WebSocket struct {
	conn   *websocket.Conn
	cMutex sync.Mutex
	log    *logrus.Entry
}

// Send send data to the clients
func (ws *WebSocket) Send(data []byte) error {
	return ws.conn.WriteMessage(websocket.BinaryMessage, data)
}

// Close closes the websocket connection
func (ws *WebSocket) Close() error {
	ws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return ws.conn.Close()
}
