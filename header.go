package mqtt

import (
	"encoding/binary"
	"errors"
)

// PacketType is the MQTT packet type.
type PacketType byte

// PacketType values.
const (
	_           PacketType = 0  // Reserved
	CONNECT     PacketType = 1  // Client request to connect to Server
	CONNACK     PacketType = 2  // Connect acknowledgment
	PUBLISH     PacketType = 3  // Publish message
	PUBACK      PacketType = 4  // Publish acknowledgment
	PUBREC      PacketType = 5  // Publish received (assured delivery part 1)
	PUBREL      PacketType = 6  // Publish release (assured delivery part 2)
	PUBCOMP     PacketType = 7  // Publish complete (assured delivery part 3)
	SUBSCRIBE   PacketType = 8  // Subscribe request
	SUBACK      PacketType = 9  // Subscribe acknowledgment
	UNSUBSCRIBE PacketType = 10 // Unsubscribe request
	UNSUBACK    PacketType = 11 // Unsubscribe acknowledgment
	PINGREQ     PacketType = 12 // PING request
	PINGRESP    PacketType = 13 // PING response
	DISCONNECT  PacketType = 14 // Client is disconnecting
	AUTH        PacketType = 15 // Authentication exchange
)

// FixedHeader is the fixed header of an MQTT packet.
type FixedHeader struct {
	typeAndFlags    byte
	remainingLength uint32
}

// PacketType returns the packet type from the fixed header.
func (h FixedHeader) PacketType() PacketType {
	return PacketType(h.typeAndFlags >> 4)
}

// SetPacketType sets the packet type into the fixed header.
func (h *FixedHeader) SetPacketType(p PacketType) {
	h.typeAndFlags |= byte(p) << 4
}

var errReservedPacketType = errors.New("mqtt: reserved packed type")
var errInvalidHeaderFlags = errors.New("mqtt: invalid header flags")

func (r *PacketReader) validateFixedHeader(h FixedHeader) error {
	switch h.PacketType() {
	case 0:
		return errReservedPacketType
	case PUBLISH:
		return r.validatePublishFlags(PublishFlags(h.typeAndFlags))
	case PUBREL, SUBSCRIBE, UNSUBSCRIBE:
		if h.typeAndFlags&0x0F == 0x02 {
			return nil
		}
	default:
		if h.typeAndFlags&0x0F == 0x00 {
			return nil
		}
	}
	return errInvalidHeaderFlags
}

const maxRemainingLength = 268435455

var errInvalidRemainingLength = errors.New("mqtt: invalid remaining length")

func (r *PacketReader) readFixedHeader() {
	r.header.typeAndFlags, r.err = r.readByte()
	if r.err != nil {
		return
	}
	var remainingLength uint64
	remainingLength, r.err = binary.ReadUvarint(r.r)
	if r.err != nil {
		return
	}
	if remainingLength > maxRemainingLength {
		r.err = errInvalidRemainingLength
		return
	}
	r.header.remainingLength = uint32(remainingLength)
	r.err = r.validateFixedHeader(r.header)
}

func (w *PacketWriter) writeFixedHeader() (err error) {
	header := w.packet.fixedHeader(w.protocol)
	if header.remainingLength > maxRemainingLength {
		return errInvalidRemainingLength
	}
	var buf [5]byte
	buf[0] = header.typeAndFlags
	n := binary.PutUvarint(buf[1:], uint64(header.remainingLength))
	_, err = w.w.Write(buf[:n+1])
	return err
}
