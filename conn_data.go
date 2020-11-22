package peer

import (
	"bytes"
	"errors"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/rs/xid"
)

const (
	//IDPrefix used as prefix for random ID
	IDPrefix = "dc_"
	//MaxBufferedAmount max amount to buffer
	MaxBufferedAmount = 8 * 1024 * 1024
)

// NewDataConnection create new DataConnection
func NewDataConnection(id string, peer *Peer, opts ConnectionOptions) (*DataConnection, error) {

	d := &DataConnection{
		BaseConnection: newBaseConnection(ConnectionTypeData, peer, opts),
		buffer:         bytes.NewBuffer([]byte{}),
	}

	d.id = opts.ConnectionID
	if d.id == "" {
		d.id = IDPrefix + xid.New().String()
	}

	d.Label = opts.Label
	if d.Label == "" {
		d.Label = d.id
	}

	d.Serialization = opts.Serialization
	if d.Serialization == "" {
		d.Serialization = SerializationTypeBinary
	}

	d.Reliable = opts.Reliable

	d.encodingQueue.On("done", d.onQueueDone)
	d.encodingQueue.On("error", d.onQueueErr)

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
	negotiator    *Negotiator
	buffer        *bytes.Buffer
	bufferSize    int
	buffering     bool
	chunkedData   map[int]chunkedData
	DataChannel   *webrtc.DataChannel
	encodingQueue *EncodingQueue
}

//   parse: (data: string) => any = JSON.parse;

func (d *DataConnection) onQueueDone(data interface{}) {
	buf := data.([]byte)
	d.bufferedSend(buf)
}

func (d *DataConnection) onQueueErr(data interface{}) {
	err := data.(error)
	d.log.Errorf(`DC#%s: Error occured in encoding from blob to arraybuffer, close DC: %s`, d.GetID(), err)
	d.Close()
}

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
		d.Emit(ConnectionEventTypeOpen, nil)
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
	panic("TODO")
	// const datatype = data.constructor;
	// isBinarySerialization := d.Serialization == SerializationTypeBinary ||
	// 	d.Serialization == SerializationTypeBinaryUTF8
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
	// Check if we've chunked--if so, piece things back together.
	// We're guaranteed that this isn't 0.
	// if deserializedData.__peerData {
	// 	d.handleChunk(deserializedData)
	// 	return
	// }
	// d.Emit(ConnectionEventTypeData, deserializedData)
}

func (d *DataConnection) handleChunk(raw []byte) {
	// const id = data.__peerData;
	// const chunkInfo = d._chunkedData[id] || {
	//   data: [],
	//   count: 0,
	//   total: data.total
	// };

	// chunkInfo.data[data.n] = data.data;
	// chunkInfo.count++;
	// d._chunkedData[id] = chunkInfo;

	// if (chunkInfo.total === chunkInfo.count) {
	//   // Clean up before making the recursive call to `_handleDataMessage`.
	//   delete d._chunkedData[id];

	//   // We've received all the chunks--time to construct the complete data.
	//   const data = new Blob(chunkInfo.data);
	//   d._handleDataMessage({ data });
	// }
}

/**
 * Exposed functionality for users.
 */

//Close allows user to close connection
func (d *DataConnection) Close() error {

	d.buffer = nil
	d.bufferSize = 0
	d.chunkedData = map[int]chunkedData{}

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

	if d.encodingQueue != nil {
		d.encodingQueue.Destroy()
		d.encodingQueue.RemoveAllListeners()
		d.encodingQueue = nil
	}

	if !d.Open {
		return nil
	}

	d.Open = false

	d.Emit(ConnectionEventTypeClose, nil)
	return nil
}

// Send allows user to send data
func (d *DataConnection) Send(data []byte, chunked bool) {
	if !d.Open {
		d.Emit(
			ConnectionEventTypeError,
			errors.New("Connection is not open. You should listen for the `open` event before sending messages"),
		)
		return
	}

	panic("TODO")

	//TODO
	// if (d.serialization == SerializationTypeJSON) {
	//   d.bufferedSend(b.toJSON(data));
	// } else if (
	//   d.Serialization == SerializationTypeBinary ||
	//   d.Serialization == SerializationTypeBinaryUTF8
	// ) {
	//   const blob = util.pack(data);

	//   if (!chunked && blob.size > util.chunkedMTU) {
	//     d._sendChunks(blob);
	//     return;
	//   }

	//   if (!util.supports.binaryBlob) {
	//     // We only do this if we really need to (e.g. blobs are not supported),
	//     // because this conversion is costly.
	//     d._encodingQueue.enque(blob);
	//   } else {
	//     d._bufferedSend(blob);
	//   }
	// } else {
	//   d._bufferedSend(data);
	// }
}

func (d *DataConnection) bufferedSend(msg []byte) {
	if d.buffering || !d.trySend(msg) {
		d.buffer.Write(msg)
		d.bufferSize = d.buffer.Len()
	}
}

// Returns true if the send succeeds.
func (d *DataConnection) trySend(msg []byte) bool {
	if !d.Open {
		return false
	}

	if d.DataChannel.BufferedAmount() > MaxBufferedAmount {
		d.buffering = true
		<-time.After(time.Millisecond * 50)
		d.buffering = false
		d.tryBuffer()
		return false
	}

	err := d.DataChannel.Send(msg)
	if err != nil {
		d.log.Errorf(`DC#%s Error sending %s`, d.GetID(), err)
		d.buffering = true
		d.Close()
		return false
	}

	return true
}

// Try to send the first message in the buffer.
func (d *DataConnection) tryBuffer() {
	if !d.Open {
		return
	}

	if d.buffer.Len() == 0 {
		return
	}

	// TODO here buffer is a slice not a continuous array
	// check or reimplement this part!
	msg := d.buffer.Bytes()
	if d.trySend(msg) {
		d.buffer.Reset()
		d.bufferSize = d.buffer.Len()
		d.tryBuffer()
	}
}

func (d *DataConnection) sendChunks(raw []byte) {
	// TODO
	panic("TODO")
	// blobs := util.chunk(raw)
	// d.log.Debugf(`DC#%s Try to send %d chunks...`, d.GetID(), len(blobs))

	// for blob := range blobs {
	// 	d.Send(blob, true)
	// }
}

func (d *DataConnection) handleMessage(message ExchangeMessage) {
	payload := message.Payload

	switch message.Type {
	case ServerMessageTypeAnswer:
		d.negotiator.handleSDP(message.Type, payload.SDP)
		break
	case ServerMessageTypeCandidate:
		panic("TODO")
		// d.negotiator.handleCandidate(payload.Candidate)
		break
	default:
		d.log.Warnf(
			"Unrecognized message type: %s from peer: %s",
			message.Type,
			d.peerID,
		)
		break
	}
}
