package mqtt

// SubscribePacket is the Subscribe packet.
type SubscribePacket struct {
	SubscribeHeader
	Properties
	SubscribePayload []Subscription
}

func (*SubscribePacket) _isPacket() {}

// PacketType returns the packet type of the Subscribe packet.
func (SubscribePacket) PacketType() PacketType { return SUBSCRIBE }

func (p SubscribePacket) fixedHeader(protocol byte) (h FixedHeader) {
	h.SetPacketType(SUBSCRIBE)
	h.typeAndFlags |= 0x02
	h.remainingLength = p.size(protocol)
	return
}

func (p SubscribePacket) size(protocol byte) uint32 {
	size := 2
	if protocol >= 5 {
		size += int(p.Properties.size())
	}
	for _, subscription := range p.SubscribePayload {
		size += 2 + len(subscription.TopicFilter) + 1
	}
	return uint32(size)
}

// SubscribeHeader is the header of the Subscribe packet.
type SubscribeHeader struct {
	PacketIdentifier uint16
}

func (r *PacketReader) readSubscribeHeader() {
	packet := r.packet.(*SubscribePacket)
	if packet.SubscribeHeader.PacketIdentifier, r.err = r.readUint16(); r.err != nil {
		return
	}
}

func (w *PacketWriter) writeSubscribeHeader() {
	packet := w.packet.(*SubscribePacket)
	w.err = w.writeUint16(packet.SubscribeHeader.PacketIdentifier)
}

// TopicFilter is the topic filter for MQTT subscriptions.
type TopicFilter []byte

// Subscription is an MQTT subscription.
type Subscription struct {
	TopicFilter TopicFilter
	QoS         QoS
}

func (r *PacketReader) readSubscribePayload() {
	packet := r.packet.(*SubscribePacket)
	for r.remaining() > 0 {
		var subscription Subscription
		if subscription.TopicFilter, r.err = r.readBytes(); r.err != nil {
			return
		}
		var b byte
		if b, r.err = r.readByte(); r.err != nil {
			return
		}
		subscription.QoS = QoS(b)
		if r.err = r.validateQoS(subscription.QoS); r.err != nil {
			return
		}
		packet.SubscribePayload = append(packet.SubscribePayload, subscription)
	}
}

func (w *PacketWriter) writeSubscribePayload() {
	packet := w.packet.(*SubscribePacket)
	for _, subscription := range packet.SubscribePayload {
		if w.err = w.writeBytes(subscription.TopicFilter); w.err != nil {
			return
		}
		if w.err = w.writeByte(byte(subscription.QoS) & 0x03); w.err != nil {
			return
		}
	}
}
