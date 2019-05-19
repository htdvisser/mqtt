package mqtt

// PingreqPacket is the Pingreq packet.
type PingreqPacket struct{}

func (*PingreqPacket) _isPacket() {}

// PacketType returns the packet type of the Pingreq packet.
func (PingreqPacket) PacketType() PacketType { return PINGREQ }

func (p PingreqPacket) fixedHeader(_ byte) (h FixedHeader) {
	h.SetPacketType(PINGREQ)
	return
}
