package enums

const (

	//ConnectionEventTypeOpen enum for connection open
	ConnectionEventTypeOpen string = "open"
	//ConnectionEventTypeStream enum for connection stream
	ConnectionEventTypeStream = "stream"
	//ConnectionEventTypeData enum for connection data
	ConnectionEventTypeData = "data"
	//ConnectionEventTypeClose enum for connection close
	ConnectionEventTypeClose = "close"
	//ConnectionEventTypeError enum for connection error
	ConnectionEventTypeError = "error"
	//ConnectionEventTypeIceStateChanged enum for ICE state changes
	ConnectionEventTypeIceStateChanged = "iceStateChanged"
	//ConnectionTypeData enum for data connection type
	ConnectionTypeData = "data"
	//ConnectionTypeMedia enum for media connection type
	ConnectionTypeMedia = "media"

	//PeerEventTypeOpen enum for peer open
	PeerEventTypeOpen = "open"
	//PeerEventTypeClose enum for peer close
	PeerEventTypeClose = "close"
	//PeerEventTypeConnection enum for peer connection
	PeerEventTypeConnection = "connection"
	//PeerEventTypeCall enum for peer call
	PeerEventTypeCall = "call"
	//PeerEventTypeDisconnected enum for peer disconnected
	PeerEventTypeDisconnected = "disconnected"
	//PeerEventTypeError enum for peer error
	PeerEventTypeError = "error"

	//PeerErrorTypeBrowserIncompatible enum for peer error  browser-incompatible
	PeerErrorTypeBrowserIncompatible = "browser-incompatible"
	//PeerErrorTypeDisconnected enum for peer error disconnected
	PeerErrorTypeDisconnected = "disconnected"
	//PeerErrorTypeInvalidID enum for  peer error invalid-id
	PeerErrorTypeInvalidID = "invalid-id"
	//PeerErrorTypeInvalidKey enum for  peer error invalid-key
	PeerErrorTypeInvalidKey = "invalid-key"
	//PeerErrorTypeNetwork enum for  peer error network
	PeerErrorTypeNetwork = "network"
	//PeerErrorTypePeerUnavailable enum for  peer error peer-unavailable
	PeerErrorTypePeerUnavailable = "peer-unavailable"
	//PeerErrorTypeSslUnavailable enum for  peer error ssl-unavailable
	PeerErrorTypeSslUnavailable = "ssl-unavailable"
	//PeerErrorTypeServerError enum for  peer error server-error
	PeerErrorTypeServerError = "server-error"
	//PeerErrorTypeSocketError enum for  peer error socket-error
	PeerErrorTypeSocketError = "socket-error"
	//PeerErrorTypeSocketClosed enum for  peer error socket-closed
	PeerErrorTypeSocketClosed = "socket-closed"
	//PeerErrorTypeUnavailableID enum for  peer error unavailable-id
	PeerErrorTypeUnavailableID = "unavailable-id"
	//PeerErrorTypeWebRTC enum for  peer error webrtc
	PeerErrorTypeWebRTC = "webrtc"

	//SerializationTypeBinary enum for binary serialization
	SerializationTypeBinary = "binary"
	//SerializationTypeBinaryUTF8 enum for UTF8 binary serialization
	SerializationTypeBinaryUTF8 = "binary-utf8"
	//SerializationTypeJSON enum for JSON serialization
	SerializationTypeJSON = "json"
	//SerializationTypeRaw Payload is sent as-is
	SerializationTypeRaw = "raw"

	//SocketEventTypeMessage enum for socket message
	SocketEventTypeMessage = "message"
	//SocketEventTypeDisconnected enum for socket disconnected
	SocketEventTypeDisconnected = "disconnected"
	//SocketEventTypeError enum for socket error
	SocketEventTypeError = "error"
	//SocketEventTypeClose enum for socket close
	SocketEventTypeClose = "close"

	//ServerMessageTypeHeartbeat enum for server HEARTBEAT
	ServerMessageTypeHeartbeat = "HEARTBEAT"
	//ServerMessageTypeCandidate enum for server CANDIDATE
	ServerMessageTypeCandidate = "CANDIDATE"
	//ServerMessageTypeOffer enum for server OFFER
	ServerMessageTypeOffer = "OFFER"
	//ServerMessageTypeAnswer enum for server ANSWER
	ServerMessageTypeAnswer = "ANSWER"
	//ServerMessageTypeOpen enum for server OPEN
	ServerMessageTypeOpen = "OPEN"
	//ServerMessageTypeError enum for server ERROR
	ServerMessageTypeError = "ERROR"
	//ServerMessageTypeIDTaken enum for server
	ServerMessageTypeIDTaken = "ID-TAKEN" // The selected ID is taken.
	//ServerMessageTypeInvalidKey enum for INVALID-KEY
	ServerMessageTypeInvalidKey = "INVALID-KEY"
	//ServerMessageTypeLeave enum for server LEAVE
	ServerMessageTypeLeave = "LEAVE"
	//ServerMessageTypeExpire enum for server EXPIRE
	ServerMessageTypeExpire = "EXPIRE"
)
