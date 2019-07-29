package mqtt

// PubrelPacket is the Pubrel packet.
type PubrelPacket struct {
	PubrelHeader
	Properties
}

func (*PubrelPacket) _isPacket() {}

// PacketType returns the packet type of the Pubrel packet.
func (PubrelPacket) PacketType() PacketType { return PUBREL }

func (p PubrelPacket) fixedHeader(protocol byte) (h FixedHeader) {
	h.SetPacketType(PUBREL)
	h.typeAndFlags |= 0x02
	h.remainingLength = p.size(protocol)
	return
}

func (p PubrelPacket) size(protocol byte) uint32 {
	size := 2
	if protocol >= 5 {
		size++
		size += int(p.Properties.size())
	}
	return uint32(size)
}

// PubrelHeader is the header of the Pubrel packet.
type PubrelHeader struct {
	PacketIdentifier uint16
	ReasonCode
}

func (r *PacketReader) readPubrelHeader() {
	packet := r.packet.(*PubrelPacket)
	if packet.PubrelHeader.PacketIdentifier, r.err = r.readUint16(); r.err != nil {
		return
	}
	if r.protocol >= 5 && r.remaining() > 0 {
		var f byte
		if f, r.err = r.readByte(); r.err != nil {
			return
		}
		packet.PubrelHeader.ReasonCode = ReasonCode(f)
	}
}

func (w *PacketWriter) writePubrelHeader() {
	packet := w.packet.(*PubrelPacket)
	if w.err = w.writeUint16(packet.PubrelHeader.PacketIdentifier); w.err != nil {
		return
	}
	if w.protocol >= 5 {
		w.err = w.writeByte(byte(packet.PubrelHeader.ReasonCode))
	}
}
