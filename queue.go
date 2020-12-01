package peer

// NewEncodingQueue initializes an EncodingQueue
// func NewEncodingQueue() *EncodingQueue {
// 	return &EncodingQueue{
// 		Emitter: NewEmitter(),
// 		Queue:   [][]byte{},
// 	}
// }

// //EncodingQueue encoding queue
// type EncodingQueue struct {
// 	Emitter
// 	Processing bool
// 	Queue      [][]byte
// }

// // Destroy destroys the queue instance
// func (e *EncodingQueue) Destroy() {
// 	e.Processing = false
// 	e.Queue = [][]byte{}
// }

// // Size return the queue size
// func (e *EncodingQueue) Size() int {
// 	return len(e.Queue)
// }

// // Enque add element to the queue
// func (e *EncodingQueue) Enque(raw []byte) {
// 	// e.queue = append(e.queue, raw)
// 	// TODO understand if conversion to ArrayBuffer is needed
// 	e.Emit("done", raw)
// }
