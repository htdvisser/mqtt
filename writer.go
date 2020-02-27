package mqtt

import (
	"encoding/binary"
	"io"
	"sync"
)

// WriterOption is an option for the PacketWriter.
type WriterOption interface {
	apply(*PacketWriter)
}

type writerOptionFunc func(*PacketWriter)

func (f writerOptionFunc) apply(w *PacketWriter) {
	f(w)
}

// PacketWriter writes MQTT packets.
type PacketWriter struct {
	w        io.Writer
	protocol byte
	mu       sync.Mutex
	nWritten uint32
	packet   Packet
	err      error
}

// SetProtocol sets the MQTT protocol version.
func (w *PacketWriter) SetProtocol(protocol byte) {
	w.mu.Lock()
	w.protocol = protocol
	w.mu.Unlock()
}

// NewWriter returns a new Writer on top of the given io.Writer.
func NewWriter(wr io.Writer, opts ...WriterOption) *PacketWriter {
	pw := &PacketWriter{
		w:        wr,
		protocol: DefaultProtocolVersion,
	}
	for _, opt := range opts {
		opt.apply(pw)
	}
	return pw
}

func (w *PacketWriter) writeVariableHeader() {
	switch w.packet.PacketType() {
	case CONNECT:
		w.writeConnectHeader()
	case CONNACK:
		w.writeConnackHeader()
	case PUBLISH:
		w.writePublishHeader()
	case PUBACK:
		w.writePubackHeader()
	case PUBREC:
		w.writePubrecHeader()
	case PUBREL:
		w.writePubrelHeader()
	case PUBCOMP:
		w.writePubcompHeader()
	case SUBSCRIBE:
		w.writeSubscribeHeader()
	case SUBACK:
		w.writeSubackHeader()
	case UNSUBSCRIBE:
		w.writeUnsubscribeHeader()
	case UNSUBACK:
		w.writeUnsubackHeader()
	case DISCONNECT:
		w.writeDisconnectHeader()
	case AUTH:
		w.writeAuthHeader()
	}
}

func (w *PacketWriter) writePayload() {
	switch w.packet.PacketType() {
	case CONNECT:
		w.writeConnectPayload()
	case PUBLISH:
		w.writePublishPayload()
	case SUBSCRIBE:
		w.writeSubscribePayload()
	case SUBACK:
		w.writeSubackPayload()
	case UNSUBSCRIBE:
		w.writeUnsubscribePayload()
	case UNSUBACK:
		w.writeUnsubackPayload()
	}
}

// WritePacket writes the given packet.
func (w *PacketWriter) WritePacket(packet Packet) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.packet = packet
	w.writeFixedHeader()
	if w.err != nil {
		return w.err
	}
	w.writeVariableHeader()
	if w.err != nil {
		return w.err
	}
	w.writePacketProperties()
	if w.err != nil {
		return w.err
	}
	w.writePayload()
	if w.err != nil {
		return w.err
	}
	return nil
}

func (w *PacketWriter) write(buf []byte) error {
	n, err := w.w.Write(buf)
	if err != nil {
		return err
	}
	w.nWritten += uint32(n)
	return nil
}

func (w *PacketWriter) writeByte(b byte) error {
	return w.write([]byte{b})
}

func (w *PacketWriter) writeUint16(i uint16) error {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], i)
	return w.write(b[:])
}

func (w *PacketWriter) writeUint32(i uint32) error {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], i)
	return w.write(b[:])
}

func (w *PacketWriter) writeUvarint(i uint32) error {
	var b [4]byte
	n := binary.PutUvarint(b[:], uint64(i))
	return w.write(b[:n])
}

var errInvalidBytesLength = NewReasonCodeError(ProtocolError, "mqtt: invalid bytes length")

func (w *PacketWriter) writeBytes(b []byte) error {
	if len(b) > 65535 {
		return errInvalidBytesLength
	}
	err := w.writeUint16(uint16(len(b)))
	if err != nil {
		return err
	}
	return w.write(b)
}

func (w *PacketWriter) writeBytesPair(k, v []byte) error {
	if err := w.writeBytes(k); err != nil {
		return err
	}
	return w.writeBytes(v)
}
