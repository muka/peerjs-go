package peer

import (
	"bytes"
	"errors"

	"github.com/muka/peerjs-go/enums"
	"github.com/muka/peerjs-go/models"
	"github.com/muka/peerjs-go/util"
	"github.com/pion/webrtc/v3"
)

const (
	//DataChannelIDPrefix used as prefix for random ID
	DataChannelIDPrefix = "dc_"
	//MaxBufferedAmount max amount to buffer
	MaxBufferedAmount = 8 * 1024 * 1024
	// ChunkedMTU payload size for a single message
	ChunkedMTU = 16300
)

// NewDataConnection create new DataConnection
func NewDataConnection(peerID string, peer *Peer, opts ConnectionOptions) (*DataConnection, error) {

	d := &DataConnection{
		BaseConnection: newBaseConnection(enums.ConnectionTypeData, peer, opts),
		buffer:         bytes.NewBuffer([]byte{}),
		// encodingQueue:  NewEncodingQueue(),
	}

	d.peerID = peerID

	d.id = opts.ConnectionID
	if d.id == "" {
		d.id = DataChannelIDPrefix + util.RandomToken()
	}

	d.Label = opts.Label
	if d.Label == "" {
		d.Label = d.id
	}

	d.Serialization = opts.Serialization
	if d.Serialization == "" {
		d.Serialization = enums.SerializationTypeRaw
	}

	d.Reliable = opts.Reliable

	// d.encodingQueue.On("done", d.onQueueDone)
	// d.encodingQueue.On("error", d.onQueueErr)

	d.negotiator = NewNegotiator(d, opts)
	err := d.negotiator.StartConnection(opts)

	return d, err
}

type chunkedData struct {
	Data  []byte
	Count int
	Total int
}

// DataConnection track a connection with a remote Peer
type DataConnection struct {
	BaseConnection
	buffer     *bytes.Buffer
	bufferSize int
	buffering  bool
	// chunkedData map[int]chunkedData
	// encodingQueue *EncodingQueue
}

//   parse: (data: string) => any = JSON.parse;

// func (d *DataConnection) onQueueDone(data interface{}) {
// 	buf := data.([]byte)
// 	d.bufferedSend(buf)
// }

// func (d *DataConnection) onQueueErr(data interface{}) {
// 	err := data.(error)
// 	d.log.Errorf(`DC#%s: Error occured in encoding from blob to arraybuffer, close DC: %s`, d.GetID(), err)
// 	d.Close()
// }

//Initialize called by the Negotiator when the DataChannel is ready
func (d *DataConnection) Initialize(dc *webrtc.DataChannel) {
	d.DataChannel = dc
	d.configureDataChannel()
}

func (d *DataConnection) configureDataChannel() {
	// TODO
	// d.DataChannel.binaryType = "arraybuffer";

	d.DataChannel.OnOpen(func() {
		//TODO
		d.log.Debugf(`DC#%s dc connection success`, d.GetID())
		d.Open = true
		d.Emit(enums.ConnectionEventTypeOpen, nil)
	})

	d.DataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		d.log.Debugf(`DC#%s dc onmessage: %v`, d.GetID(), msg.Data)
		d.handleDataMessage(msg)
	})

	d.DataChannel.OnClose(func() {
		d.log.Debugf(`DC#%s dc closed for %s`, d.GetID(), d.peerID)
		d.Close()
	})

}

// Handles a DataChannel message.
func (d *DataConnection) handleDataMessage(msg webrtc.DataChannelMessage) {

	// isBinarySerialization := d.Serialization == SerializationTypeBinary ||
	// 	d.Serialization == SerializationTypeBinaryUTF8

	if msg.IsString {
		d.Emit(enums.ConnectionEventTypeData, string(msg.Data))
	} else {
		d.Emit(enums.ConnectionEventTypeData, msg.Data)
	}

	// if (isBinarySerialization) {
	// 	if (datatype == Blob) {
	//     // Datatype should never be blob
	//     util.blobToArrayBuffer(data as Blob, (ab) => {
	//       const unpackedData = util.unpack(ab);
	//       d.emit(ConnectionEventType.Data, unpackedData);
	//     });
	//     return;
	//   } else if (datatype === ArrayBuffer) {
	//     deserializedData = util.unpack(data as ArrayBuffer);
	//   } else if (datatype === String) {
	//     // String fallback for binary data for browsers that don't support binary yet
	//     const ab = util.binaryStringToArrayBuffer(data as string);
	//     deserializedData = util.unpack(ab);
	//   }
	// } else if (d.serialization === SerializationType.JSON) {
	//   deserializedData = d.parse(data as string);
	// }
	// // Check if we've chunked--if so, piece things back together.
	// // We're guaranteed that this isn't 0.
	// if deserializedData.__peerData {
	// 	d.handleChunk(deserializedData)
	// 	return
	// }
	// d.Emit(ConnectionEventTypeData, deserializedData)
}

// func (d *DataConnection) handleChunk(raw []byte) {
// 	// const id = data.__peerData;
// 	// const chunkInfo = d._chunkedData[id] || {
// 	//   data: [],
// 	//   count: 0,
// 	//   total: data.total
// 	// };

// 	// chunkInfo.data[data.n] = data.data;
// 	// chunkInfo.count++;
// 	// d._chunkedData[id] = chunkInfo;

// 	// if (chunkInfo.total === chunkInfo.count) {
// 	//   // Clean up before making the recursive call to `_handleDataMessage`.
// 	//   delete d._chunkedData[id];

// 	//   // We've received all the chunks--time to construct the complete data.
// 	//   const data = new Blob(chunkInfo.data);
// 	//   d._handleDataMessage({ data });
// 	// }
// }

/**
 * Exposed functionality for users.
 */

//Close allows user to close connection
func (d *DataConnection) Close() error {

	d.buffer = nil
	d.bufferSize = 0
	// d.chunkedData = map[int]chunkedData{}

	if d.negotiator != nil {
		d.negotiator.Cleanup()
		d.negotiator = nil
	}

	if d.Provider != nil {
		d.Provider.RemoveConnection(d)
		d.Provider = nil
	}

	if d.DataChannel != nil {
		d.DataChannel.OnOpen(func() {})
		d.DataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {})
		d.DataChannel.OnClose(func() {})
		d.DataChannel = nil
	}

	// if d.encodingQueue != nil {
	// 	d.encodingQueue.Destroy()
	// 	d.encodingQueue = nil
	// }

	if !d.Open {
		return nil
	}

	d.Open = false

	d.Emit(enums.ConnectionEventTypeClose, nil)
	return nil
}

// Send allows user to send data.
func (d *DataConnection) Send(data []byte, chunked bool) error {
	if !d.Open {
		err := errors.New("Connection is not open. You should listen for the `open` event before sending messages")
		d.Emit(
			enums.ConnectionEventTypeError,
			err,
		)
		return err
	}

	err := d.DataChannel.Send(data)
	if err != nil {
		d.log.Warnf("Send failed: %s", err)
		return err
	}

	return nil

	// if d.Serialization == SerializationTypeJSON {
	// 	// JSON data must be marshalled before send!
	// 	d.log.Debug("Send JSON")
	// 	d.bufferedSend(raw)
	// } else if d.Serialization == SerializationTypeBinary || d.Serialization == SerializationTypeBinaryUTF8 {

	// 	panic(errors.New("binarypack encoding is not supported"))

	// 	// NOTE we pack with MessagePack not with binarypack. Understant if this is good enough
	// 	// blob, err := msgpack.Marshal(data)
	// 	// if err != nil {
	// 	// 	return fmt.Errorf("Failed to pack message: %s", err)
	// 	// }

	// 	// if !chunked && len(blob) > ChunkedMTU {
	// 	// 	d.log.Debug("Chunk payload")
	// 	// 	d.sendChunks(blob)
	// 	// 	return nil
	// 	// }

	// 	// d.log.Debugf("Send encoded payload %v", raw)
	// 	// d.bufferedSend(blob)

	// } else {
	// 	d.log.Debug("Send raw payload")
	// 	d.bufferedSend(raw)
	// }

}

// func (d *DataConnection) bufferedSend(msg []byte) {
// 	if d.buffering || !d.trySend(msg) {
// 		d.buffer.Write(msg)
// 		d.bufferSize = d.buffer.Len()
// 	}
// }

// // Returns true if the send succeeds.
// func (d *DataConnection) trySend(msg []byte) bool {
// 	if !d.Open {
// 		return false
// 	}

// 	if d.DataChannel.BufferedAmount() > MaxBufferedAmount {
// 		d.buffering = true
// 		<-time.After(time.Millisecond * 50)
// 		d.buffering = false
// 		d.tryBuffer()
// 		return false
// 	}

// 	err := d.DataChannel.Send(msg)
// 	if err != nil {
// 		d.log.Errorf(`DC#%s Error sending %s`, d.GetID(), err)
// 		d.buffering = true
// 		// d.Close()
// 		return false
// 	}

// 	return true
// }

// Try to send the first message in the buffer.
// func (d *DataConnection) tryBuffer() {
// 	if !d.Open {
// 		return
// 	}

// 	if d.buffer.Len() == 0 {
// 		return
// 	}

// 	// TODO here buffer is a slice not a continuous array
// 	// check or reimplement this part!
// 	msg := d.buffer.Bytes()
// 	if d.trySend(msg) {
// 		d.buffer.Reset()
// 		d.bufferSize = d.buffer.Len()
// 		d.tryBuffer()
// 	}
// }

// func (d *DataConnection) sendChunks(raw []byte) {
// 	panic("sendChunks: binarypack not implemented, please use SerializationTypeRaw")
// 	// // this method requires a [binarypack] encoding to work
// 	// chunks := util.Chunk(raw)
// 	// d.log.Debugf(`DC#%s Try to send %d chunks...`, d.GetID(), len(chunks))
// 	// for _, chunk := range chunks {
// 	// 	d.Send(chunk, true)
// 	// }
// }

// HandleMessage handles incoming messages
func (d *DataConnection) HandleMessage(message *models.Message) error {
	payload := message.Payload

	switch message.Type {
	case enums.ServerMessageTypeAnswer:
		d.negotiator.handleSDP(message.Type, *payload.SDP)
		break
	case enums.ServerMessageTypeCandidate:
		err := d.negotiator.HandleCandidate(payload.Candidate)
		if err != nil {
			d.log.Errorf("Failed to handle candidate for peer=%s: %s", d.peerID, err)
		}
		break
	default:
		d.log.Warnf(
			"Unrecognized message type: %s from peer: %s",
			message.Type,
			d.peerID,
		)
		break
	}

	return nil
}
