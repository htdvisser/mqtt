package mqtt

// PingrespPacket is the Pingresp packet.
type PingrespPacket struct{}

func (*PingrespPacket) _isPacket() {}

// PacketType returns the packet type of the Pingresp packet.
func (PingrespPacket) PacketType() PacketType { return PINGRESP }

func (p PingrespPacket) fixedHeader(_ byte) (h FixedHeader) {
	h.SetPacketType(PINGRESP)
	return
}
