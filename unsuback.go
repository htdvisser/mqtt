package mqtt

// UnsubackPacket is the Unsuback packet.
type UnsubackPacket struct {
	UnsubackHeader
	Properties
}

func (*UnsubackPacket) _isPacket() {}

// PacketType returns the packet type of the Unsuback packet.
func (UnsubackPacket) PacketType() PacketType { return UNSUBACK }

func (p UnsubackPacket) fixedHeader(protocol byte) (h FixedHeader) {
	h.SetPacketType(UNSUBACK)
	h.remainingLength = p.size(protocol)
	return
}

func (p UnsubackPacket) size(protocol byte) uint32 {
	size := 2
	if protocol >= 5 {
		size += int(p.Properties.size())
	}
	return uint32(size)
}

// UnsubackHeader is the header of the Unsuback packet.
type UnsubackHeader struct {
	PacketIdentifier uint16
}

func (r *PacketReader) readUnsubackHeader() {
	packet := r.packet.(*UnsubackPacket)
	if packet.UnsubackHeader.PacketIdentifier, r.err = r.readUint16(); r.err != nil {
		return
	}
}

func (w *PacketWriter) writeUnsubackHeader() {
	packet := w.packet.(*UnsubackPacket)
	w.err = w.writeUint16(packet.UnsubackHeader.PacketIdentifier)
}
