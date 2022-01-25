package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/muka/peerjs-go/emitter"
	"github.com/muka/peerjs-go/models"
	"github.com/sirupsen/logrus"
)

// ClientMessage wrap a message received by a client
type ClientMessage struct {
	Client  IClient
	Message *models.Message
}

//NewWebSocketServer create a new WebSocketServer
func NewWebSocketServer(realm IRealm, opts Options) *WebSocketServer {
	wss := WebSocketServer{
		Emitter:  emitter.NewEmitter(),
		upgrader: websocket.Upgrader{},
		log:      createLogger("websocket-server", opts),
		realm:    realm,
		opts:     opts,
	}

	wss.upgrader.CheckOrigin = func(r *http.Request) bool {
		// TODO: expose as options
		return true
	}

	return &wss
}

// WebSocketServer wrap the websocket server
type WebSocketServer struct {
	emitter.Emitter
	upgrader websocket.Upgrader
	clients  []*websocket.Conn
	cMutex   sync.Mutex
	log      *logrus.Entry
	realm    IRealm
	opts     Options
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

//onSocketConnection called when a client connect
func (wss *WebSocketServer) sendErrorAndClose(conn *websocket.Conn, msg string) error {
	err := conn.WriteJSON(models.Message{
		Type: MessageTypeError,
		Payload: models.Payload{
			Msg: msg,
		},
	})
	if err != nil {
		return err
	}
	err = conn.Close()
	if err != nil {
		return err
	}
	return nil
}

//
func (wss *WebSocketServer) configureWS(conn *websocket.Conn, client IClient) error {
	client.SetSocket(conn)

	conn.SetPingHandler(func(appData string) error {
		// wss.log.Debugf("[%s] Ping received", client.GetID())
		client.SetLastPing(getTime())
		return nil
	})

	closed := false
	conn.SetCloseHandler(func(code int, text string) error {
		// if any close error happens, stop the loop and remove the client
		wss.log.Debugf("Closed connection, cleaning up %s", client.GetID())
		if client.GetSocket() == conn {
			wss.realm.RemoveClientByID(client.GetID())
		}
		conn.Close()
		wss.Emit(WebsocketEventClose, client)
		closed = true
		return nil
	})

	go func() {
		for {
			if closed {
				return
			}
			_, raw, err := conn.ReadMessage()
			if err != nil {
				wss.log.Errorf("[%s] Read WS error: %s", client.GetID(), err)
				return
			}

			// message handling
			data, err := ioutil.ReadAll(bytes.NewReader(raw))
			if err != nil {
				wss.log.Errorf("client message read error: %s", err)
				wss.Emit(WebsocketEventError, err)
				continue
			}

			message := new(models.Message)
			err = json.Unmarshal(data, message)
			if err != nil {
				wss.log.Errorf("client message unmarshal error: %s", err)
				wss.Emit(WebsocketEventError, err)
				continue
			}

			message.Src = client.GetID()
			wss.Emit(WebsocketEventMessage, ClientMessage{client, message})
		}
	}()

	wss.Emit(WebsocketEventConnection, client)
	return nil
}

//registerClient
func (wss *WebSocketServer) registerClient(conn *websocket.Conn, id, token string) error {
	// Check concurrent limit
	clientsCount := len(wss.realm.GetClientsIds())

	if clientsCount >= wss.opts.ConcurrentLimit {
		err := wss.sendErrorAndClose(conn, ErrorConnectionLimitExceeded)
		if err != nil {
			wss.log.Errorf("[sendErrorAndClose] Error: %s", err)
		}
		return nil
	}

	client := NewClient(id, token)
	wss.realm.SetClient(client, id)

	err := conn.WriteJSON(models.Message{Type: MessageTypeOpen})
	if err != nil {
		return err
	}

	err = wss.configureWS(conn, client)
	if err != nil {
		return err
	}
	return nil
}

//onSocketConnection called when a client connect
func (wss *WebSocketServer) onSocketConnection(conn *websocket.Conn, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	key := query.Get("key")

	if id == "" || token == "" || key == "" {
		err := wss.sendErrorAndClose(conn, ErrorInvalidWSParameters)
		if err != nil {
			wss.log.Errorf("[sendErrorAndClose] Error: %s", err)
		}
		return
	}

	if key != wss.opts.Key {
		err := wss.sendErrorAndClose(conn, ErrorInvalidKey)
		if err != nil {
			wss.log.Errorf("[sendErrorAndClose] Error: %s", err)
		}
		return
	}

	client := wss.realm.GetClientByID(id)

	if client == nil {
		err := wss.registerClient(conn, id, token)
		if err != nil {
			wss.log.Errorf("[registerClient] Error: %s", err)
		}
		return
	}

	if token != client.GetToken() {
		// ID-taken, invalid token
		err := conn.WriteJSON(models.Message{
			Type: MessageTypeIDTaken,
			Payload: models.Payload{
				Msg: "ID is taken",
			},
		})
		if err != nil {
			wss.log.Errorf("[%s] Failed to write message: %s", MessageTypeIDTaken, err)
		}
		go func() {
			// wait for the client to receive the response message
			<-time.After(time.Millisecond * 100)
			conn.Close()
		}()
		return
	}

	wss.configureWS(conn, client)
}

// Handler expose the http handler for websocket
func (wss *WebSocketServer) Handler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := wss.upgrader.Upgrade(w, r, nil)
		if err != nil {
			wss.log.Warnf("upgrade error: %s", err)
			w.WriteHeader(500)
			// next.ServeHTTP(w, r)
			return
		}
		wss.onSocketConnection(c, r)
	})
}
