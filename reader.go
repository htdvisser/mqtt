package mqtt

import (
	"bufio"
	"encoding/binary"
	"io"
	"sync"
	"unicode/utf8"
)

// ReaderOption is an option for the PacketReader.
type ReaderOption interface {
	apply(*PacketReader)
}

type readerOptionFunc func(*PacketReader)

func (f readerOptionFunc) apply(r *PacketReader) {
	f(r)
}

// WithMaxPacketLength returns a ReaderOption that configures a maximum packet
// length on the Reader.
func WithMaxPacketLength(bytes uint32) ReaderOption {
	return readerOptionFunc(func(r *PacketReader) {
		r.maxPacketLength = bytes
	})
}

type reader interface {
	io.Reader
	io.ByteReader
}

// PacketReader reads MQTT packets.
type PacketReader struct {
	maxPacketLength uint32
	r               reader
	protocol        byte
	mu              sync.Mutex
	nRead           uint32
	header          FixedHeader
	packet          Packet
	err             error
}

// SetProtocol sets the MQTT protocol version.
func (r *PacketReader) SetProtocol(protocol byte) {
	r.mu.Lock()
	r.protocol = protocol
	r.mu.Unlock()
}

// NewReader returns a new Reader on top of the given io.Reader.
func NewReader(rd io.Reader, opts ...ReaderOption) *PacketReader {
	pr := &PacketReader{
		protocol: DefaultProtocolVersion,
	}
	if r, ok := rd.(reader); ok {
		pr.r = r
	} else {
		pr.r = bufio.NewReader(rd)
	}
	for _, opt := range opts {
		opt.apply(pr)
	}
	return pr
}

var errUnknownPacket = NewReasonCodeError(ProtocolError, "mqtt: unknown packet")

func (r *PacketReader) readVariableHeader() {
	switch r.packet.PacketType() {
	case CONNECT:
		r.readConnectHeader()
		r.protocol = r.packet.(*ConnectPacket).ConnectHeader.ProtocolVersion
	case CONNACK:
		r.readConnackHeader()
	case PUBLISH:
		r.readPublishHeader()
	case PUBACK:
		r.readPubackHeader()
	case PUBREC:
		r.readPubrecHeader()
	case PUBREL:
		r.readPubrelHeader()
	case PUBCOMP:
		r.readPubcompHeader()
	case SUBSCRIBE:
		r.readSubscribeHeader()
	case SUBACK:
		r.readSubackHeader()
	case UNSUBSCRIBE:
		r.readUnsubscribeHeader()
	case UNSUBACK:
		r.readUnsubackHeader()
	case DISCONNECT:
		r.readDisconnectHeader()
	case AUTH:
		if r.protocol < 5 {
			r.err = errUnknownPacket
			return
		}
		r.readAuthHeader()
	}
}

var errRemainingData = NewReasonCodeError(ProtocolError, "mqtt: unexpected remaining data after reading packet")

func (r *PacketReader) readPayload() {
	switch r.packet.PacketType() {
	case CONNECT:
		r.readConnectPayload()
	case PUBLISH:
		r.readPublishPayload()
	case SUBSCRIBE:
		r.readSubscribePayload()
	case SUBACK:
		r.readSubackPayload()
	case UNSUBSCRIBE:
		r.readUnsubscribePayload()
	case UNSUBACK:
		r.readUnsubackPayload()
	default:
		if r.remaining() > 0 {
			r.err = errRemainingData
		}
	}
}

// ReadPacket reads the next packet.
func (r *PacketReader) ReadPacket() (Packet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nRead = 0
	r.readFixedHeader()
	if r.err != nil {
		return nil, r.err
	}
	switch r.header.PacketType() {
	case CONNECT:
		r.packet = new(ConnectPacket)
	case CONNACK:
		r.packet = new(ConnackPacket)
	case PUBLISH:
		packet := new(PublishPacket)
		packet.PublishFlags = PublishFlags(r.header.typeAndFlags) & 0xf
		r.packet = packet
	case PUBACK:
		r.packet = new(PubackPacket)
	case PUBREC:
		r.packet = new(PubrecPacket)
	case PUBREL:
		r.packet = new(PubrelPacket)
	case PUBCOMP:
		r.packet = new(PubcompPacket)
	case SUBSCRIBE:
		r.packet = new(SubscribePacket)
	case SUBACK:
		r.packet = new(SubackPacket)
	case UNSUBSCRIBE:
		r.packet = new(UnsubscribePacket)
	case UNSUBACK:
		r.packet = new(UnsubackPacket)
	case PINGREQ:
		r.packet = new(PingreqPacket)
	case PINGRESP:
		r.packet = new(PingrespPacket)
	case DISCONNECT:
		r.packet = new(DisconnectPacket)
	case AUTH:
		r.packet = new(AuthPacket)
	}
	r.nRead = 0
	r.readVariableHeader()
	if r.err != nil {
		return nil, r.err
	}
	r.readPacketProperties()
	if r.err != nil {
		return nil, r.err
	}
	r.readPayload()
	if r.err != nil {
		return nil, r.err
	}
	return r.packet, nil
}

var errInsufficientRemainingBytes = NewReasonCodeError(MalformedPacket, "insufficient remaining bytes")

func (r *PacketReader) read(b []byte) error {
	if r.remaining() < uint32(len(b)) {
		return errInsufficientRemainingBytes
	}
	n, err := io.ReadFull(r.r, b)
	if err != nil {
		return err
	}
	r.nRead += uint32(n)
	return nil
}

func (r *PacketReader) readByte() (b byte, err error) {
	if r.remaining() < 1 {
		return 0, errInsufficientRemainingBytes
	}
	b, err = r.r.ReadByte()
	if err != nil {
		return 0, err
	}
	r.nRead++
	return
}

func (r *PacketReader) readUint16() (uint16, error) {
	var b [2]byte
	err := r.read(b[:])
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(b[:]), nil
}

func (r *PacketReader) readUint32() (uint32, error) {
	var b [4]byte
	err := r.read(b[:])
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b[:]), nil
}

type countingByteReader struct {
	*PacketReader
}

func (r *PacketReader) readUvarint() (i uint64, err error) {
	return binary.ReadUvarint(countingByteReader{r})
}

func (r countingByteReader) ReadByte() (byte, error) {
	return r.PacketReader.readByte()
}

func (r *PacketReader) readBytes() ([]byte, error) {
	var length uint16
	length, err := r.readUint16()
	if err != nil {
		return nil, err
	}
	if length == 0 {
		return nil, nil
	}
	b := make([]byte, length)
	err = r.read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

var errInvalidUTF8 = NewReasonCodeError(MalformedPacket, "mqtt: invalid utf-8 string")

func (r *PacketReader) readString() ([]byte, error) {
	b, err := r.readBytes()
	if err != nil {
		return nil, err
	}
	if r.protocol >= 4 && !utf8.Valid(b) {
		return nil, errInvalidUTF8
	}
	return b, nil
}

func (r *PacketReader) readStringPair() (k, v []byte, err error) {
	k, err = r.readString()
	if err != nil {
		return nil, nil, err
	}
	v, err = r.readString()
	if err != nil {
		return nil, nil, err
	}
	return k, v, nil
}

func (r *PacketReader) remaining() uint32 {
	return r.header.remainingLength - r.nRead
}

func (r *PacketReader) readRemaining() ([]byte, error) {
	b := make([]byte, int(r.remaining()))
	err := r.read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
