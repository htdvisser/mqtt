package mqtt

// ConnackPacket is the Connack packet.
type ConnackPacket struct {
	ConnackHeader
	Properties
}

func (*ConnackPacket) _isPacket() {}

// PacketType returns the packet type of the Connack packet.
func (ConnackPacket) PacketType() PacketType { return CONNACK }

func (p ConnackPacket) fixedHeader(protocol byte) (h FixedHeader) {
	h.SetPacketType(CONNACK)
	h.remainingLength = p.size(protocol)
	return
}

func (p ConnackPacket) size(protocol byte) uint32 {
	size := 2
	if protocol >= 5 {
		size += int(p.Properties.size())
	}
	return uint32(size)
}

// ConnackHeader is the header of the Connack packet.
type ConnackHeader struct {
	ConnackHeaderFlags
	ReasonCode
}

// ConnackHeaderFlags are the flags in the header of the Connack packet.
type ConnackHeaderFlags byte

var errInvalidConnackHeaderFlags = NewReasonCodeError(ProtocolError, "mqtt: invalid connack header flags")

func (r *PacketReader) validateConnackHeaderFlags(f ConnackHeaderFlags) error {
	if r.protocol < 4 && f != 0x00 {
		return errInvalidConnackHeaderFlags
	}
	if f&0xFE != 0x00 {
		return errInvalidConnackHeaderFlags
	}
	return nil
}

// SessionPresent returns Session Present bit from the connack header flags.
func (f ConnackHeaderFlags) SessionPresent() bool { return f&0x01 == 0x01 }

// SetSessionPresent sets the Session Present bit into the connack header flags.
func (f *ConnackHeaderFlags) SetSessionPresent(sessionPresent bool) {
	*f &^= 0x01
	if sessionPresent {
		*f |= 0x01
	}
}

func (r *PacketReader) readConnackHeader() {
	packet := r.packet.(*ConnackPacket)
	var f byte
	if f, r.err = r.readByte(); r.err != nil {
		return
	}
	packet.ConnackHeader.ConnackHeaderFlags = ConnackHeaderFlags(f)
	if r.err = r.validateConnackHeaderFlags(packet.ConnackHeader.ConnackHeaderFlags); r.err != nil {
		return
	}
	if f, r.err = r.readByte(); r.err != nil {
		return
	}
	packet.ConnackHeader.ReasonCode = ReasonCode(f)
}

func (w *PacketWriter) writeConnackHeader() {
	packet := w.packet.(*ConnackPacket)
	var flags byte
	if w.protocol >= 4 {
		flags = byte(packet.ConnackHeader.ConnackHeaderFlags)
	}
	if w.err = w.writeByte(flags); w.err != nil {
		return
	}
	w.err = w.writeByte(byte(packet.ConnackHeader.ReasonCode))
}
