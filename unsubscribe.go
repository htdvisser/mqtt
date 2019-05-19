package mqtt

// UnsubscribePacket is the Unsubscribe packet.
type UnsubscribePacket struct {
	UnsubscribeHeader
	Properties
	UnsubscribePayload []TopicFilter
}

func (*UnsubscribePacket) _isPacket() {}

// PacketType returns the packet type of the Unsubscribe packet.
func (UnsubscribePacket) PacketType() PacketType { return UNSUBSCRIBE }

func (p UnsubscribePacket) fixedHeader(protocol byte) (h FixedHeader) {
	h.SetPacketType(UNSUBSCRIBE)
	h.typeAndFlags |= 0x02
	h.remainingLength = p.size(protocol)
	return
}

func (p UnsubscribePacket) size(protocol byte) uint32 {
	size := 2
	if protocol >= 5 {
		size += int(p.Properties.size())
	}
	for _, topicFilter := range p.UnsubscribePayload {
		size += 2 + len(topicFilter)
	}
	return uint32(size)
}

// UnsubscribeHeader is the header of the Unsubscribe packet.
type UnsubscribeHeader struct {
	PacketIdentifier uint16
}

func (r *PacketReader) readUnsubscribeHeader() {
	packet := r.packet.(*UnsubscribePacket)
	if packet.UnsubscribeHeader.PacketIdentifier, r.err = r.readUint16(); r.err != nil {
		return
	}
}

func (w *PacketWriter) writeUnsubscribeHeader() {
	packet := w.packet.(*UnsubscribePacket)
	w.err = w.writeUint16(packet.UnsubscribeHeader.PacketIdentifier)
}

func (r *PacketReader) readUnsubscribePayload() {
	packet := r.packet.(*UnsubscribePacket)
	for r.remaining() > 0 {
		var topicFilter TopicFilter
		if topicFilter, r.err = r.readBytes(); r.err != nil {
			return
		}
		packet.UnsubscribePayload = append(packet.UnsubscribePayload, topicFilter)
	}
}

func (w *PacketWriter) writeUnsubscribePayload() {
	packet := w.packet.(*UnsubscribePacket)
	for _, topicFilter := range packet.UnsubscribePayload {
		if w.err = w.writeBytes(topicFilter); w.err != nil {
			return
		}
	}
}
