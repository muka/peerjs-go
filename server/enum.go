package server

const (
	// ErrorInvalidKey Invalid key provided
	ErrorInvalidKey = "Invalid key provided"
	// ErrorInvalidToken Invalid token provided
	ErrorInvalidToken = "Invalid token provided"
	// ErrorInvalidWSParameters No id, token, or key supplied to websocket server
	ErrorInvalidWSParameters = "No id, token, or key supplied to websocket server"
	// ErrorConnectionLimitExceeded Server has reached its concurrent user limit
	ErrorConnectionLimitExceeded = "Server has reached its concurrent user limit"
	// MessageTypeOpen OPEN
	MessageTypeOpen = "OPEN"
	// MessageTypeLeave LEAVE
	MessageTypeLeave = "LEAVE"
	// MessageTypeCandidate CANDIDATE
	MessageTypeCandidate = "CANDIDATE"
	// MessageTypeOffer OFFER
	MessageTypeOffer = "OFFER"
	// MessageTypeAnswer ANSWER
	MessageTypeAnswer = "ANSWER"
	// MessageTypeExpire EXPIRE
	MessageTypeExpire = "EXPIRE"
	// MessageTypeHeartbeat HEARTBEAT
	MessageTypeHeartbeat = "HEARTBEAT"
	// MessageTypeIDTaken ID-TAKEN
	MessageTypeIDTaken = "ID-TAKEN"
	// MessageTypeError ERROR
	MessageTypeError = "ERROR"

	// WebsocketEventMessage message
	WebsocketEventMessage = "message"
	// WebsocketEventConnection connection
	WebsocketEventConnection = "connection"
	// WebsocketEventError error
	WebsocketEventError = "error"
	// WebsocketEventClose close
	WebsocketEventClose = "close"
)
