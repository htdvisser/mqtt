package mqtt

// PubcompPacket is the Pubcomp packet.
type PubcompPacket struct {
	PubcompHeader
	Properties
}

func (*PubcompPacket) _isPacket() {}

// PacketType returns the packet type of the Pubcomp packet.
func (PubcompPacket) PacketType() PacketType { return PUBCOMP }

func (p PubcompPacket) fixedHeader(protocol byte) (h FixedHeader) {
	h.SetPacketType(PUBCOMP)
	h.remainingLength = p.size(protocol)
	return
}

func (p PubcompPacket) size(protocol byte) uint32 {
	size := 2
	if protocol >= 5 {
		size++
		size += int(p.Properties.size())
	}
	return uint32(size)
}

// PubcompHeader is the header of the Pubcomp packet.
type PubcompHeader struct {
	PacketIdentifier uint16
	ReasonCode
}

func (r *PacketReader) readPubcompHeader() {
	packet := r.packet.(*PubcompPacket)
	if packet.PubcompHeader.PacketIdentifier, r.err = r.readUint16(); r.err != nil {
		return
	}
	if r.protocol >= 5 && r.remaining() > 0 {
		var f byte
		if f, r.err = r.readByte(); r.err != nil {
			return
		}
		packet.PubcompHeader.ReasonCode = ReasonCode(f)
	}
}

func (w *PacketWriter) writePubcompHeader() {
	packet := w.packet.(*PubcompPacket)
	if w.err = w.writeUint16(packet.PubcompHeader.PacketIdentifier); w.err != nil {
		return
	}
	if w.protocol >= 5 {
		w.err = w.writeByte(byte(packet.PubcompHeader.ReasonCode))
	}
}
