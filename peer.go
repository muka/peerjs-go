package peer

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/muka/peerjs-go/emitter"
	"github.com/muka/peerjs-go/enums"
	"github.com/muka/peerjs-go/models"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

//DefaultKey is the default API key
var DefaultKey = "peerjs"

var socketEvents = []string{
	enums.SocketEventTypeMessage,
	enums.SocketEventTypeError,
	enums.SocketEventTypeDisconnected,
	enums.SocketEventTypeClose,
}

type socketEventWrapper struct {
	Event string
	Data  interface{}
}

//NewPeer initializes a new Peer object
func NewPeer(id string, opts Options) (*Peer, error) {
	p := &Peer{
		Emitter:      emitter.NewEmitter(),
		opts:         opts,
		api:          NewAPI(opts),
		socket:       NewSocket(opts),
		lostMessages: make(map[string][]models.Message),
		connections:  make(map[string]map[string]Connection),
	}

	if id == "" {
		raw, err := p.api.RetrieveID()
		id = string(raw)
		if err != nil {
			return p, err
		}
	}

	p.log = createLogger(fmt.Sprintf("peer:%s", id), opts.Debug)

	err := p.initialize(id)
	if err != nil {
		return p, err
	}

	return p, nil
}

//Peer expose the PeerJS API
type Peer struct {
	emitter.Emitter
	ID           string
	opts         Options
	connections  map[string]map[string]Connection
	api          API
	socket       *Socket
	log          *logrus.Entry
	open         bool
	destroyed    bool
	disconnected bool
	lastServerID string
	lostMessages map[string][]models.Message
}

//GetSocket return a socket connection
func (p *Peer) GetSocket() *Socket {
	return p.socket
}

//GetOptions return options
func (p *Peer) GetOptions() Options {
	return p.opts
}

//AddConnection add the connection to the peer
func (p *Peer) AddConnection(peerID string, connection Connection) {
	if _, ok := p.connections[peerID]; !ok {
		p.connections[peerID] = make(map[string]Connection)
	}
	p.connections[peerID][connection.GetID()] = connection
}

//RemoveConnection removes the connection from the peer
func (p *Peer) RemoveConnection(connection Connection) {
	peerID := connection.GetPeerID()
	id := connection.GetID()
	if connections, ok := p.connections[peerID]; ok {
		for id := range connections {
			if id == connection.GetID() {
				delete(p.connections[peerID], id)
			}
		}
	}
	// remove lost messages
	if _, ok := p.lostMessages[id]; ok {
		delete(p.lostMessages, id)
	}
}

//GetConnection return a connection based on peerID and connectionID
func (p *Peer) GetConnection(peerID string, connectionID string) (Connection, bool) {
	_, ok := p.connections[peerID]
	if !ok {
		return nil, false
	}
	conn, ok := p.connections[peerID][connectionID]
	return conn, ok
}

func (p *Peer) messageHandler(msg SocketEvent) {
	peerID := msg.Message.GetSrc()
	payload := msg.Message.GetPayload()
	switch msg.Message.GetType() {
	case enums.ServerMessageTypeOpen:
		p.lastServerID = p.ID
		p.open = true
		p.log.Debugf("Open session with id=%s", p.ID)
		p.Emit(enums.PeerEventTypeOpen, p.ID)
		break
	case enums.ServerMessageTypeError:
		p.abort(enums.PeerErrorTypeServerError, msg.Error)
		break
	case enums.ServerMessageTypeIDTaken: // The selected ID is taken.
		p.abort(enums.PeerErrorTypeUnavailableID, fmt.Errorf("ID %s is taken", p.ID))
		break
	case enums.ServerMessageTypeInvalidKey: // The given API key cannot be found.
		p.abort(enums.PeerErrorTypeInvalidKey, fmt.Errorf("API KEY %s is invalid", p.opts.Key))
		break
	case enums.ServerMessageTypeLeave: // Another peer has closed its connection to this peer.
		peerID := msg.Message.GetSrc()
		p.log.Debugf("Received leave message from %s", peerID)
		p.cleanupPeer(peerID)
		if _, ok := p.connections[peerID]; ok {
			delete(p.connections, peerID)
		}
		break
	case enums.ServerMessageTypeExpire: // The offer sent to a peer has expired without response.
		p.EmitError(enums.PeerErrorTypePeerUnavailable, fmt.Errorf("Could not connect to peer %s", peerID))
		break
	case enums.ServerMessageTypeOffer:

		// we should consider switching this to CALL/CONNECT, but this is the least breaking option.
		connectionID := payload.ConnectionID
		connection, ok := p.GetConnection(peerID, connectionID)

		if ok {
			connection.Close()
			p.log.Warnf("Offer received for existing Connection ID %s", connectionID)
		}

		var err error
		// Create a new connection.
		if payload.Type == enums.ConnectionTypeMedia {
			connection, err = NewMediaConnection(peerID, p, ConnectionOptions{
				ConnectionID: connectionID,
				Payload:      payload,
				Metadata:     payload.Metadata,
			})
			if err != nil {
				p.log.Errorf("Cannot initialize MediaConnection: %s", err)
				return
			}
			p.AddConnection(peerID, connection)
			p.Emit(enums.PeerEventTypeCall, connection)
		} else if payload.Type == enums.ConnectionTypeData {
			connection, err = NewDataConnection(peerID, p, ConnectionOptions{
				ConnectionID:  connectionID,
				Payload:       payload,
				Metadata:      payload.Metadata,
				Label:         payload.Label,
				Serialization: payload.Serialization,
				Reliable:      payload.Reliable,
				SDP:           *payload.SDP,
			})
			if err != nil {
				p.log.Errorf("Cannot initialize DataConnection: %s", err)
				return
			}
			p.AddConnection(peerID, connection)
			p.Emit(enums.PeerEventTypeConnection, connection)
		} else {
			p.log.Warnf(`Received malformed connection type:%s`, payload.Type)
			return
		}

		// Find messages.
		messages := p.GetMessages(connectionID)
		for _, message := range messages {
			connection.HandleMessage(&message)
		}

		break
	default:

		if msg.Message == nil {
			p.log.Warnf(`You received a malformed message from %s of type %s`, peerID, msg.Type)
			return
		}

		connectionID := msg.Message.GetPayload().ConnectionID
		connection, ok := p.GetConnection(peerID, connectionID)

		if ok && connection.GetPeerConnection() != nil {
			// Pass it on.
			connection.HandleMessage(msg.Message)
		} else if connectionID != "" {
			// Store for possible later use
			p.storeMessage(connectionID, *msg.Message)
		} else {
			p.log.Warnf("You received an unrecognized message: %v", msg.Message)
		}
		break
	}
}

// handles socket events
func (p *Peer) socketEventHandler(data interface{}) {
	ev := data.(SocketEvent)
	switch ev.Type {
	case enums.SocketEventTypeMessage:
		p.messageHandler(ev)
		break
	case enums.SocketEventTypeError:
		p.abort(enums.PeerErrorTypeSocketError, ev.Error)
		break
	case enums.SocketEventTypeDisconnected:
		if p.disconnected {
			return
		}
		p.EmitError(enums.PeerErrorTypeNetwork, errors.New("Lost connection to server"))
		p.disconnect()
		break
	case enums.SocketEventTypeClose:
		if p.disconnected {
			return
		}
		p.abort(enums.PeerErrorTypeSocketClosed, errors.New("Underlying socket is already closed"))
		break
	}
}

func (p *Peer) unregisterSocketHandlers() {
	for _, messageType := range socketEvents {
		p.socket.Off(messageType, p.socketEventHandler)
	}
}

func (p *Peer) registerSocketHandlers() {
	for _, messageType := range socketEvents {
		p.socket.On(messageType, p.socketEventHandler)
	}
}

// Stores messages without a set up connection, to be claimed later
func (p *Peer) storeMessage(connectionID string, message models.Message) {
	if _, ok := p.lostMessages[connectionID]; !ok {
		p.lostMessages[connectionID] = []models.Message{}
	}
	p.lostMessages[connectionID] = append(p.lostMessages[connectionID], message)
}

//GetMessages Retrieve messages from lost message store
func (p *Peer) GetMessages(connectionID string) []models.Message {
	if messages, ok := p.lostMessages[connectionID]; ok {
		delete(p.lostMessages, connectionID)
		return messages
	}
	return []models.Message{}
}

//Close closes the peer instance
func (p *Peer) Close() {
	if p.lastServerID != "" {
		p.destroy()
	} else {
		p.disconnect()
	}
}

//Connect returns a DataConnection to the specified peer. See documentation for a
//complete list of options.
func (p *Peer) Connect(peerID string, opts *ConnectionOptions) (*DataConnection, error) {

	if opts == nil {
		opts = NewConnectionOptions()
	}

	if p.disconnected {
		p.log.Warn(`
	  You cannot connect to a new Peer because you called .disconnect() on this Peer 
	  and ended your connection with the server. You can create a new Peer to reconnect, 
	  or call reconnect on this peer if you believe its ID to still be available`)
		err := errors.New("Cannot connect to new Peer after disconnecting from server")
		p.EmitError(
			enums.PeerErrorTypeDisconnected,
			err,
		)
		return nil, err
	}

	// indicate we are starting the connection
	opts.Originator = true

	if opts.Debug == -1 {
		opts.Debug = p.opts.Debug
	}

	dataConnection, err := NewDataConnection(peerID, p, *opts)
	if err != nil {
		return dataConnection, err
	}

	p.AddConnection(peerID, dataConnection)
	return dataConnection, nil
}

//Call returns a MediaConnection to the specified peer. See documentation for a
//complete list of options.
func (p *Peer) Call(peerID string, track webrtc.TrackLocal, opts *ConnectionOptions) (*MediaConnection, error) {

	if opts == nil {
		opts = NewConnectionOptions()
	}

	if p.disconnected {
		p.log.Warn("You cannot connect to a new Peer because you called .disconnect() on this Peer and ended your connection with the server. You can create a new Peer to reconnect")
		err := errors.New("Cannot connect to new Peer after disconnecting from server")
		p.EmitError(
			enums.PeerErrorTypeDisconnected,
			err,
		)
		return nil, err
	}

	if track == nil {
		err := errors.New("To call a peer, you must provide a stream")
		p.log.Error(err)
		return nil, err
	}

	opts.Stream = NewMediaStreamWithTrack([]MediaStreamTrack{track})

	mediaConnection, err := NewMediaConnection(peerID, p, *opts)
	if err != nil {
		p.log.Errorf("Failed to create a MediaConnection: %s", err)
		return nil, err
	}
	p.AddConnection(peerID, mediaConnection)
	return mediaConnection, nil
}

func (p *Peer) abort(errType string, err error) error {
	p.log.Error("Aborting!")
	p.EmitError(errType, err)
	p.Close()
	return err
}

//EmitError emits an error
func (p *Peer) EmitError(errType string, err error) {
	p.log.Errorf("Error: %s", err)
	p.Emit(enums.PeerEventTypeError, err)
}

func (p *Peer) initialize(id string) error {
	p.log.Debugf("Initializing id=%s", id)
	p.ID = id
	//register event handler
	p.registerSocketHandlers()
	return p.socket.Start(id, p.opts.Token)
}

// destroys the Peer: closes all active connections as well as the connection
// to the server.
// Warning: The peer can no longer create or accept connections after being
// destroyed.
func (p *Peer) destroy() {

	if p.destroyed {
		return
	}

	p.log.Debugf(`Destroy peer with ID:%s`, p.ID)

	p.disconnect()
	p.cleanup()

	p.destroyed = true

	p.Emit(enums.PeerEventTypeClose, nil)
}

// cleanup Disconnects every connection on this peer.
func (p *Peer) cleanup() {
	for peerID := range p.connections {
		p.cleanupPeer(peerID)
		delete(p.connections, peerID)
	}

	err := p.socket.Close()
	p.socket = nil
	if err != nil {
		p.log.Warnf("Failed to close socket: %s", err)
	}
}

// cleanupPeer Closes all connections to this peer.
func (p *Peer) cleanupPeer(peerID string) {
	connections, ok := p.connections[peerID]
	if !ok {
		return
	}
	for _, connection := range connections {
		connection.Close()
	}
}

// disconnect disconnects the Peer's connection to the PeerServer. Does not close any
// active connections.
// Warning: The peer can no longer create or accept connections after being
// disconnected. It also cannot reconnect to the server.
func (p *Peer) disconnect() {
	if p.disconnected {
		return
	}

	currentID := p.ID

	p.log.Debugf("Disconnect peer with ID:%s", currentID)

	p.disconnected = true
	p.open = false

	// remove registered handlers
	p.unregisterSocketHandlers()
	p.socket.Close()

	p.lastServerID = currentID
	p.ID = ""

	p.Emit(enums.PeerEventTypeDisconnected, currentID)
}

// reconnect Attempts to reconnect with the same ID
func (p *Peer) reconnect() error {

	if p.disconnected && !p.destroyed {
		p.log.Debugf(`Attempting reconnection to server with ID %s`, p.lastServerID)
		p.disconnected = false
		p.initialize(p.lastServerID)
		return nil
	}

	if p.destroyed {
		return errors.New("This peer cannot reconnect to the server. It has already been destroyed")
	}

	if !p.disconnected && !p.open {
		// Do nothing. We're still connecting the first time.
		p.log.Error("In a hurry? We're still trying to make the initial connection!")
		return nil
	}

	return fmt.Errorf(`Peer %s cannot reconnect because it is not disconnected from the server`, p.ID)
}

// ListAllPeers Get a list of available peer IDs. If you're running your own server, you'll
// want to set allow_discovery: true in the PeerServer options. If you're using
// the cloud server, email team@peerjs.com to get the functionality enabled for
// your key.
func (p *Peer) ListAllPeers() ([]string, error) {

	peers := []string{}
	raw, err := p.api.ListAllPeers()
	if err != nil {
		return peers, p.abort(enums.PeerErrorTypeServerError, err)
	}

	err = json.Unmarshal(raw, &peers)
	if err != nil {
		return peers, p.abort(enums.PeerErrorTypeServerError, err)
	}

	return peers, nil
}
