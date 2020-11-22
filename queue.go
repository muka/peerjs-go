package peer

// NewEncodingQueue initializes an EncodingQueue
func NewEncodingQueue() EncodingQueue {
	return EncodingQueue{
		Emitter: NewEmitter(),
	}
}

//EncodingQueue encoding queue
type EncodingQueue struct {
	Emitter
}

// Destroy destroys the queue instance
func (e *EncodingQueue) Destroy() {
	panic("TODO")
}

// RemoveAllListeners remove all event listener
func (e *EncodingQueue) RemoveAllListeners() {
	panic("TODO")
}
