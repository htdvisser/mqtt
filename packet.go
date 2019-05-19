package mqtt

// Packet interface for MQTT packets from this package.
type Packet interface {
	_isPacket()
	PacketType() PacketType
	fixedHeader(protocol byte) FixedHeader
}
