package mqtt

import "errors"

// PublishPacket is the Publish packet.
type PublishPacket struct {
	PublishFlags
	PublishHeader
	Properties
	PublishPayload []byte
}

func (*PublishPacket) _isPacket() {}

// PacketType returns the packet type of the Publish packet.
func (PublishPacket) PacketType() PacketType { return PUBLISH }

func (p PublishPacket) fixedHeader(protocol byte) (h FixedHeader) {
	h.SetPacketType(PUBLISH)
	h.typeAndFlags |= byte(p.PublishFlags)
	h.remainingLength = p.size(protocol)
	return
}

func (p PublishPacket) size(protocol byte) uint32 {
	size := 2 + len(p.PublishHeader.TopicName)
	if p.PublishFlags.QoS() > 0 {
		size += 2
	}
	if protocol >= 5 {
		size += int(p.Properties.size())
	}
	size += len(p.PublishPayload)
	return uint32(size)
}

func (p PublishPacket) publishFlags() PublishFlags { return p.PublishFlags }

// QoS is the MQTT quality of service of a Publish packet.
type QoS byte

// QoS values.
const (
	QoS0 QoS = 0 // At Most Once
	QoS1 QoS = 1 // At Least Once
	QoS2 QoS = 2 // Exactly Once
)

func (r *PacketReader) validateQoS(qos QoS) error {
	if qos > 2 {
		return errInvalidQoS
	}
	return nil
}

// PublishFlags are the fixed header flags for a Publish packet.
type PublishFlags byte

var errInvalidQoS = errors.New("mqtt: invalid QoS")

func (r *PacketReader) validatePublishFlags(f PublishFlags) error {
	return r.validateQoS(f.QoS())
}

// Dup returns the Dup bit from the publish flags.
func (f PublishFlags) Dup() bool { return f&0x8 == 0x8 }

// SetDup sets the Dup bit into the publish flags.
func (f *PublishFlags) SetDup(dup bool) {
	*f &^= 0x8
	if dup {
		*f |= 0x8
	}
}

// QoS returns the QoS from the publish flags.
func (f PublishFlags) QoS() QoS { return QoS(f >> 1 & 0x03) }

// SetQoS sets the QoS into the publish flags.
func (f *PublishFlags) SetQoS(qos QoS) {
	*f &^= 0x6
	switch qos {
	case 1:
		*f |= 0x2
	case 2:
		*f |= 0x4
	}
}

// Retain returns the Retain bit from the publish flags.
func (f PublishFlags) Retain() bool { return f&0x1 == 0x1 }

// SetRetain sets the Retain bit into the publish flags.
func (f *PublishFlags) SetRetain(retain bool) {
	*f &^= 0x1
	if retain {
		*f |= 0x1
	}
}

// PublishHeader is the header of the Publish packet.
type PublishHeader struct {
	TopicName        []byte
	PacketIdentifier uint16
}

func (r *PacketReader) readPublishHeader() {
	packet := r.packet.(*PublishPacket)
	if packet.PublishHeader.TopicName, r.err = r.readBytes(); r.err != nil {
		return
	}
	if packet.PublishFlags.QoS() > 0 {
		if packet.PublishHeader.PacketIdentifier, r.err = r.readUint16(); r.err != nil {
			return
		}
	}
}

func (w *PacketWriter) writePublishHeader() {
	packet := w.packet.(*PublishPacket)
	if w.err = w.writeBytes(packet.PublishHeader.TopicName); w.err != nil {
		return
	}
	if PublishFlags(packet.PublishFlags).QoS() > 0 {
		if w.err = w.writeUint16(packet.PublishHeader.PacketIdentifier); w.err != nil {
			return
		}
	}
}

func (r *PacketReader) readPublishPayload() {
	packet := r.packet.(*PublishPacket)
	packet.PublishPayload, r.err = r.readRemaining()
}

func (w *PacketWriter) writePublishPayload() {
	packet := w.packet.(*PublishPacket)
	w.err = w.write(packet.PublishPayload)
}
