package mqtt

// PubrecPacket is the Pubrec packet.
type PubrecPacket struct {
	PubrecHeader
	Properties
}

func (*PubrecPacket) _isPacket() {}

// PacketType returns the packet type of the Pubrec packet.
func (PubrecPacket) PacketType() PacketType { return PUBREC }

func (p PubrecPacket) fixedHeader(protocol byte) (h FixedHeader) {
	h.SetPacketType(PUBREC)
	h.remainingLength = p.size(protocol)
	return
}

func (p PubrecPacket) size(protocol byte) uint32 {
	size := 2
	if protocol >= 5 {
		size++
		size += int(p.Properties.size())
	}
	return uint32(size)
}

// PubrecHeader is the header of the Pubrec packet.
type PubrecHeader struct {
	PacketIdentifier uint16
	ReasonCode
}

func (r *PacketReader) readPubrecHeader() {
	packet := r.packet.(*PubrecPacket)
	if packet.PubrecHeader.PacketIdentifier, r.err = r.readUint16(); r.err != nil {
		return
	}
	if r.protocol >= 5 {
		var f byte
		if f, r.err = r.readByte(); r.err != nil {
			return
		}
		packet.PubrecHeader.ReasonCode = ReasonCode(f)
	}
}

func (w *PacketWriter) writePubrecHeader() {
	packet := w.packet.(*PubrecPacket)
	if w.err = w.writeUint16(packet.PubrecHeader.PacketIdentifier); w.err != nil {
		return
	}
	if w.protocol >= 5 {
		w.err = w.writeByte(byte(packet.PubrecHeader.ReasonCode))
	}
}
