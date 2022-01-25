package peer

import (
	"encoding/json"
	"fmt"

	"github.com/muka/peerjs-go/enums"
	"github.com/muka/peerjs-go/models"
	"github.com/pion/webrtc/v3"
	"github.com/sirupsen/logrus"
)

// DefaultBrowser is the browser name
const DefaultBrowser = "peerjs-go"

func newWebrtcAPI() *webrtc.API {
	mediaEngine := new(webrtc.MediaEngine)
	mediaEngine.RegisterDefaultCodecs()
	return webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))
}

//NewNegotiator initiate a new negotiator
func NewNegotiator(conn Connection, opts ConnectionOptions) *Negotiator {
	return &Negotiator{
		connection: conn,
		log:        createLogger("negotiator", opts.Debug),
		webrtc:     newWebrtcAPI(),
	}
}

// Negotiator manages all negotiations between Peers
type Negotiator struct {
	connection Connection
	log        *logrus.Entry
	webrtc     *webrtc.API
}

//StartConnection Returns a PeerConnection object set up correctly (for data, media). */
func (n *Negotiator) StartConnection(opts ConnectionOptions) error {

	peerConnection, err := n.startPeerConnection()
	if err != nil {
		return err
	}

	// Set the connection's PC.
	n.connection.SetPeerConnection(peerConnection)

	if n.connection.GetType() == enums.ConnectionTypeMedia && opts.Stream != nil {
		for _, track := range opts.Stream.GetTracks() {
			peerConnection.AddTrack(track.(webrtc.TrackLocal))
		}
	}

	// What do we need to do now?
	if opts.Originator {
		if n.connection.GetType() == enums.ConnectionTypeData {

			dataConnection := n.connection.(*DataConnection)

			config := &webrtc.DataChannelInit{
				Ordered: &opts.Reliable,
			}

			dataChannel, err := peerConnection.CreateDataChannel(dataConnection.Label, config)
			if err != nil {
				return err
			}

			dataConnection.Initialize(dataChannel)
		}

		n.makeOffer()
	} else {
		// OFFER
		err = n.handleSDP(enums.ServerMessageTypeOffer, opts.SDP)
		if err != nil {
			return err
		}
	}

	return nil
}

// Start a PC
func (n *Negotiator) startPeerConnection() (*webrtc.PeerConnection, error) {

	n.log.Debug("Creating RTCPeerConnection")

	// peerConnection = webrtc.PeerConnection(this.connection.provider.options.config);
	c := n.connection.GetProvider().GetOptions().Configuration
	peerConnection, err := n.webrtc.NewPeerConnection(c)
	if err != nil {
		return nil, err
	}

	n.setupListeners(peerConnection)

	return peerConnection, nil
}

// Set up various WebRTC listeners
func (n *Negotiator) setupListeners(peerConnection *webrtc.PeerConnection) {

	peerID := n.connection.GetPeerID()
	connectionID := n.connection.GetID()
	provider := n.connection.GetProvider()

	n.log.Debug("Listening for ICE candidates.")
	peerConnection.OnICECandidate(func(evt *webrtc.ICECandidate) {

		peerID := n.connection.GetPeerID()
		connectionID := n.connection.GetID()
		connectionType := n.connection.GetType()
		provider := n.connection.GetProvider()

		if evt == nil {
			n.log.Debugf("ICECandidate gathering completed for peer=%s conn=%s", peerID, connectionID)
			return
		}

		candidate := evt.ToJSON()

		if candidate.Candidate == "" {
			return
		}

		n.log.Debugf("Received ICE candidates for %s: %s", peerID, candidate.Candidate)

		msg := models.Message{
			Type: enums.ServerMessageTypeCandidate,
			Payload: models.Payload{
				Candidate:    &candidate,
				Type:         connectionType,
				ConnectionID: connectionID,
			},
			Dst: peerID,
		}

		res, err := json.Marshal(msg)
		if err != nil {
			n.log.Errorf("OnICECandidate: Failed to serialize message: %s", err)
		}

		err = provider.GetSocket().Send(res)
		if err != nil {
			n.log.Errorf("OnICECandidate: Failed to send message: %s", err)
		}

	})

	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		switch peerConnection.ICEConnectionState() {
		case webrtc.ICEConnectionStateFailed:
			n.log.Debugf("iceConnectionState is failed, closing connections to %s", peerID)
			n.connection.Emit(
				enums.ConnectionEventTypeError,
				fmt.Errorf("Negotiation of connection to %s failed", peerID),
			)
			n.connection.Close()
			break
		case webrtc.ICEConnectionStateClosed:
			n.log.Debugf("iceConnectionState is closed, closing connections to %s", peerID)
			n.connection.Emit(enums.ConnectionEventTypeError, fmt.Errorf("Connection to %s closed", peerID))
			n.connection.Close()
			break
		case webrtc.ICEConnectionStateDisconnected:
			n.log.Debugf("iceConnectionState changed to disconnected on the connection with %s", peerID)
			break
		case webrtc.ICEConnectionStateCompleted:
			peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
				// noop
			})
			break
		}

		n.connection.Emit(enums.ConnectionEventTypeIceStateChanged, peerConnection.ICEConnectionState())
	})

	// DATACONNECTION.
	n.log.Debug("Listening for data channel")

	// Fired between offer and answer, so options should already be saved in the options hash.
	peerConnection.OnDataChannel(func(dataChannel *webrtc.DataChannel) {
		n.log.Debug("Received data channel")
		conn, ok := provider.GetConnection(peerID, connectionID)
		if ok {
			connection := conn.(*DataConnection)
			connection.Initialize(dataChannel)
		}
	})

	// MEDIACONNECTION
	n.log.Debug("Listening for remote stream")

	peerConnection.OnTrack(func(tr *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		n.log.Debug("Received remote stream")

		connection, ok := provider.GetConnection(peerID, connectionID)
		if ok {
			if connection.GetType() == enums.ConnectionTypeMedia {
				mediaConnection := connection.(*MediaConnection)
				n.log.Debugf("add stream %s to media connection %s", tr.ID(), mediaConnection.GetID())
				mediaConnection.AddStream(tr)
			}
		}
	})
}

// Cleanup clean up the negotiatior internal state
func (n *Negotiator) Cleanup() {
	n.log.Debugf("Cleaning up PeerConnection to %s", n.connection.GetPeerID())

	peerConnection := n.connection.GetPeerConnection()
	if peerConnection == nil {
		return
	}

	n.connection.SetPeerConnection(nil)

	//unsubscribe from all PeerConnection's events
	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {})
	peerConnection.OnICEConnectionStateChange(func(is webrtc.ICEConnectionState) {})
	peerConnection.OnDataChannel(func(dc *webrtc.DataChannel) {})
	peerConnection.OnTrack(func(tr *webrtc.TrackRemote, r *webrtc.RTPReceiver) {})

	peerConnectionNotClosed := peerConnection.ConnectionState() != webrtc.PeerConnectionStateClosed
	dataChannelNotClosed := false

	if n.connection.GetType() == enums.ConnectionTypeData {
		dataConnection := n.connection.(*DataConnection)
		dataChannel := dataConnection.DataChannel

		if dataChannel != nil {
			dataChannelNotClosed = dataChannel.ReadyState() != webrtc.DataChannelStateClosed
		}
	}

	if peerConnectionNotClosed || dataChannelNotClosed {
		peerConnection.Close()
	}

}

func (n *Negotiator) makeOffer() error {

	peerConnection := n.connection.GetPeerConnection()
	provider := n.connection.GetProvider()

	// TODO check offer message
	offer, err := peerConnection.CreateOffer(&webrtc.OfferOptions{
		OfferAnswerOptions: webrtc.OfferAnswerOptions{
			// VoiceActivityDetection: true,
		},
		ICERestart: false,
	})
	if err != nil {
		err1 := fmt.Errorf("makeOffer: Failed to create offer: %s", err)
		n.log.Warn(err1)
		provider.EmitError(enums.PeerErrorTypeWebRTC, err1)
		return err
	}
	n.log.Debug("Created offer")

	connOpts := n.connection.GetOptions()

	if connOpts.SDPTransform != nil {
		offer.SDP = connOpts.SDPTransform(offer.SDP)
	}

	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		err1 := fmt.Errorf("makeOffer: Failed to set local description: %s", err)
		n.log.Warn(err1)
		provider.EmitError(enums.PeerErrorTypeWebRTC, err1)
		return err
	}

	n.log.Debugf("Set localDescription: %s for:%s", offer.SDP, n.connection.GetPeerID())

	payload := models.Payload{
		Type:         n.connection.GetType(),
		ConnectionID: n.connection.GetID(),
		Metadata:     n.connection.GetMetadata(),
		SDP:          &offer,
		Browser:      DefaultBrowser,
	}

	if n.connection.GetType() == enums.ConnectionTypeData {
		dataConnection := n.connection.(*DataConnection)
		payload.Label = dataConnection.Label
		payload.Reliable = dataConnection.Reliable
		payload.Serialization = dataConnection.Serialization
	}

	msg := models.Message{
		Type:    enums.ServerMessageTypeOffer,
		Dst:     n.connection.GetPeerID(),
		Payload: payload,
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		err1 := fmt.Errorf("makeOffer: Failed to marshal socket message: %s", err)
		n.log.Warn(err1)
		provider.EmitError(enums.PeerErrorTypeWebRTC, err1)
		return err
	}

	err = provider.GetSocket().Send(raw)
	if err != nil {
		err1 := fmt.Errorf("makeOffer: Failed to send message: %s", err)
		n.log.Warn(err1)
		provider.EmitError(enums.PeerErrorTypeWebRTC, err1)
		return err
	}

	return nil
}

func (n *Negotiator) makeAnswer() error {

	peerConnection := n.connection.GetPeerConnection()
	provider := n.connection.GetProvider()

	answer, err := peerConnection.CreateAnswer(&webrtc.AnswerOptions{
		OfferAnswerOptions: webrtc.OfferAnswerOptions{
			// VoiceActivityDetection: true,
		},
	})
	if err != nil {
		err1 := fmt.Errorf("makeAnswer: Failed to create answer: %s", err)
		n.log.Warn(err1)
		provider.EmitError(enums.PeerErrorTypeWebRTC, err1)
		return err
	}

	n.log.Debug("Created answer.")

	connOpts := n.connection.GetOptions()
	if connOpts.SDPTransform != nil {
		answer.SDP = connOpts.SDPTransform(answer.SDP)
	}

	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		err1 := fmt.Errorf("makeAnswer: Failed to set local description: %s", err)
		n.log.Warn(err1)
		provider.EmitError(enums.PeerErrorTypeWebRTC, err1)
		return err
	}

	n.log.Debugf(`Set localDescription: %s for %s`, answer.SDP, n.connection.GetPeerID())

	msg := models.Message{
		Type: enums.ServerMessageTypeAnswer,
		Dst:  n.connection.GetPeerID(),
		Payload: models.Payload{
			Type:         n.connection.GetType(),
			ConnectionID: n.connection.GetID(),
			SDP:          &answer,
			Browser:      DefaultBrowser,
		},
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		err1 := fmt.Errorf("makeAnswer: Failed to marshal sockt message: %s", err)
		n.log.Warn(err1)
		provider.EmitError(enums.PeerErrorTypeWebRTC, err1)
		return err
	}

	err = provider.GetSocket().Send(raw)
	if err != nil {
		err1 := fmt.Errorf("makeAnswer: Failed to send message: %s", err)
		n.log.Warn(err1)
		provider.EmitError(enums.PeerErrorTypeWebRTC, err1)
		return err
	}

	return nil
}

// Handle an SDP.
func (n *Negotiator) handleSDP(sdpType string, sdp webrtc.SessionDescription) error {

	peerConnection := n.connection.GetPeerConnection()
	provider := n.connection.GetProvider()

	n.log.Debugf("Setting remote description %v", sdp)

	err := peerConnection.SetRemoteDescription(sdp)
	if err != nil {
		provider.EmitError(enums.PeerErrorTypeWebRTC, err)
		n.log.Warnf("handleSDP: Failed to setRemoteDescription %s", err)
		return err
	}

	n.log.Debugf(`Set remoteDescription:%s for:%s`, sdpType, n.connection.GetPeerID())

	// sdpType == OFFER
	if sdpType == enums.ServerMessageTypeOffer {
		err := n.makeAnswer()
		if err != nil {
			return err
		}
	}

	return nil
}

// HandleCandidate handles a candidate
func (n *Negotiator) HandleCandidate(iceInit *webrtc.ICECandidateInit) error {

	n.log.Debugf(`HandleCandidate: %v`, iceInit)

	// candidate := ice.ToJSON().Candidate
	// sdpMLineIndex := ice.ToJSON().SDPMLineIndex
	// sdpMid := ice.ToJSON().SDPMid

	peerConnection := n.connection.GetPeerConnection()
	provider := n.connection.GetProvider()

	err := peerConnection.AddICECandidate(*iceInit)
	if err != nil {
		provider.EmitError(enums.PeerErrorTypeWebRTC, err)
		n.log.Errorf("handleCandidate: %s", err)
		return err
	}

	n.log.Debugf(`Added ICE candidate for:%s`, n.connection.GetPeerID())

	return nil
}
