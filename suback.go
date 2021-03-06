package mqtt

var errInvalidSubscribeReasonCode = NewReasonCodeError(ProtocolError, "mqtt: invalid subscribe reason code")

func (r *PacketReader) validateSubscribeReasonCode(c ReasonCode) error {
	switch {
	case r.protocol >= 5:
		switch c {
		case GrantedQoS0,
			GrantedQoS1,
			GrantedQoS2,
			UnspecifiedError,
			ImplementationSpecificError,
			NotAuthorized,
			TopicFilterInvalid,
			PacketIdentifierInUse,
			QuotaExceeded,
			SharedSubscriptionsNotSupported,
			SubscriptionIdentifiersNotSupported,
			WildcardSubscriptionsNotSupported:
			return nil
		default:
			return errInvalidSubscribeReasonCode
		}
	default:
		switch c {
		case GrantedQoS0,
			GrantedQoS1,
			GrantedQoS2,
			UnspecifiedError:
			return nil
		default:
			return errInvalidSubscribeReasonCode
		}
	}
}

// SubackPacket is the Suback packet.
type SubackPacket struct {
	SubackHeader
	Properties
	SubackPayload []ReasonCode
}

func (*SubackPacket) _isPacket() {}

// PacketType returns the packet type of the Suback packet.
func (SubackPacket) PacketType() PacketType { return SUBACK }

func (p SubackPacket) fixedHeader(protocol byte) (h FixedHeader) {
	h.SetPacketType(SUBACK)
	h.remainingLength = p.size(protocol)
	return
}

func (p SubackPacket) size(protocol byte) uint32 {
	size := 2
	if protocol >= 5 {
		size += int(p.Properties.size())
	}
	size += len(p.SubackPayload)
	return uint32(size)
}

// SubackHeader is the header of the Suback packet.
type SubackHeader struct {
	PacketIdentifier uint16
}

func (r *PacketReader) readSubackHeader() {
	packet := r.packet.(*SubackPacket)
	if packet.SubackHeader.PacketIdentifier, r.err = r.readUint16(); r.err != nil {
		return
	}
}

func (w *PacketWriter) writeSubackHeader() {
	packet := w.packet.(*SubackPacket)
	w.err = w.writeUint16(packet.SubackHeader.PacketIdentifier)
}

func (r *PacketReader) readSubackPayload() {
	packet := r.packet.(*SubackPacket)
	for r.remaining() > 0 {
		var b byte
		if b, r.err = r.readByte(); r.err != nil {
			return
		}
		returnCode := ReasonCode(b)
		if r.err = r.validateSubscribeReasonCode(returnCode); r.err != nil {
			return
		}
		packet.SubackPayload = append(packet.SubackPayload, returnCode)
	}
}

func (w *PacketWriter) writeSubackPayload() {
	packet := w.packet.(*SubackPacket)
	switch {
	case w.protocol < 4:
		for _, returnCode := range packet.SubackPayload {
			if returnCode.IsError() {
				w.err = NewReasonCodeError(returnCode, "")
				return
			}
		}
	case w.protocol < 5:
		for i, returnCode := range packet.SubackPayload {
			if returnCode.IsError() {
				packet.SubackPayload[i] = UnspecifiedError
			}
		}
	}
	for _, returnCode := range packet.SubackPayload {
		if w.err = w.writeByte(byte(returnCode)); w.err != nil {
			return
		}
	}
}
