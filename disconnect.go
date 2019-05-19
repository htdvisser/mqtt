package mqtt

// DisconnectPacket is the Disconnect packet.
type DisconnectPacket struct {
	DisconnectHeader
	Properties
}

func (*DisconnectPacket) _isPacket() {}

// PacketType returns the packet type of the Disconnect packet.
func (DisconnectPacket) PacketType() PacketType { return DISCONNECT }

func (p DisconnectPacket) fixedHeader(protocol byte) (h FixedHeader) {
	h.SetPacketType(DISCONNECT)
	h.remainingLength = p.size(protocol)
	return
}

func (p DisconnectPacket) size(protocol byte) uint32 {
	size := 0
	if protocol >= 5 {
		size++
		size += int(p.Properties.size())
	}
	return uint32(size)
}

// DisconnectHeader is the header of the Disconnect packet.
type DisconnectHeader struct {
	ReasonCode
}

func (r *PacketReader) readDisconnectHeader() {
	packet := r.packet.(*DisconnectPacket)
	if r.protocol >= 5 {
		var f byte
		if f, r.err = r.readByte(); r.err != nil {
			return
		}
		packet.DisconnectHeader.ReasonCode = ReasonCode(f)
	}
}

func (w *PacketWriter) writeDisconnectHeader() {
	packet := w.packet.(*DisconnectPacket)
	if w.protocol >= 5 {
		w.err = w.writeByte(byte(packet.DisconnectHeader.ReasonCode))
	}
}
