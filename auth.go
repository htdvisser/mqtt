package mqtt

// AuthPacket is the Auth packet.
type AuthPacket struct {
	AuthHeader
	Properties
}

func (*AuthPacket) _isPacket() {}

// PacketType returns the packet type of the Auth packet.
func (AuthPacket) PacketType() PacketType { return AUTH }

func (p AuthPacket) fixedHeader(protocol byte) (h FixedHeader) {
	h.SetPacketType(AUTH)
	h.remainingLength = p.size(protocol)
	return
}

func (p AuthPacket) size(protocol byte) uint32 {
	size := 0
	if protocol >= 5 {
		size++
		size += int(p.Properties.size())
	}
	return uint32(size)
}

// AuthHeader is the header of the Auth packet.
type AuthHeader struct {
	ReasonCode
}

func (r *PacketReader) readAuthHeader() {
	packet := r.packet.(*AuthPacket)
	if r.protocol >= 5 {
		var f byte
		if f, r.err = r.readByte(); r.err != nil {
			return
		}
		packet.AuthHeader.ReasonCode = ReasonCode(f)
	}
}

func (w *PacketWriter) writeAuthHeader() {
	packet := w.packet.(*AuthPacket)
	if w.protocol >= 5 {
		w.err = w.writeByte(byte(packet.AuthHeader.ReasonCode))
	}
}
