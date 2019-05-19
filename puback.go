package mqtt

// PubackPacket is the Puback packet.
type PubackPacket struct {
	PubackHeader
	Properties
}

func (*PubackPacket) _isPacket() {}

// PacketType returns the packet type of the Puback packet.
func (PubackPacket) PacketType() PacketType { return PUBACK }

func (p PubackPacket) fixedHeader(protocol byte) (h FixedHeader) {
	h.SetPacketType(PUBACK)
	h.remainingLength = p.size(protocol)
	return
}

func (p PubackPacket) size(protocol byte) uint32 {
	size := 2
	if protocol >= 5 {
		size++
		size += int(p.Properties.size())
	}
	return uint32(size)
}

// PubackHeader is the header of the Puback packet.
type PubackHeader struct {
	PacketIdentifier uint16
	ReasonCode
}

func (r *PacketReader) readPubackHeader() {
	packet := r.packet.(*PubackPacket)
	if packet.PubackHeader.PacketIdentifier, r.err = r.readUint16(); r.err != nil {
		return
	}
	if r.protocol >= 5 {
		var f byte
		if f, r.err = r.readByte(); r.err != nil {
			return
		}
		packet.PubackHeader.ReasonCode = ReasonCode(f)
	}
}

func (w *PacketWriter) writePubackHeader() {
	packet := w.packet.(*PubackPacket)
	if w.err = w.writeUint16(packet.PubackHeader.PacketIdentifier); w.err != nil {
		return
	}
	if w.protocol >= 5 {
		w.err = w.writeByte(byte(packet.PubackHeader.ReasonCode))
	}
}
