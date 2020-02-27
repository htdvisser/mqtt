package mqtt

var errInvalidUnsubscribeReasonCode = NewReasonCodeError(ProtocolError, "mqtt: invalid unsubscribe reason code")

func (r *PacketReader) validateUnsubscribeReasonCode(c ReasonCode) error {
	switch {
	case r.protocol >= 5:
		switch c {
		case Success,
			NoSubscriptionExisted,
			UnspecifiedError,
			ImplementationSpecificError,
			NotAuthorized,
			TopicFilterInvalid,
			PacketIdentifierInUse:
			return nil
		default:
			return errInvalidSubscribeReasonCode
		}
	default:
		return nil
	}
}

// UnsubackPacket is the Unsuback packet.
type UnsubackPacket struct {
	UnsubackHeader
	Properties
	UnsubackPayload []ReasonCode
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
		size += len(p.UnsubackPayload)
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

func (r *PacketReader) readUnsubackPayload() {
	packet := r.packet.(*UnsubackPacket)
	if r.protocol >= 5 {
		for r.remaining() > 0 {
			var b byte
			if b, r.err = r.readByte(); r.err != nil {
				return
			}
			returnCode := ReasonCode(b)
			if r.err = r.validateUnsubscribeReasonCode(returnCode); r.err != nil {
				return
			}
			packet.UnsubackPayload = append(packet.UnsubackPayload, returnCode)
		}
	}
}

func (w *PacketWriter) writeUnsubackPayload() {
	packet := w.packet.(*UnsubackPacket)
	if w.protocol >= 5 {
		for _, returnCode := range packet.UnsubackPayload {
			if w.err = w.writeByte(byte(returnCode)); w.err != nil {
				return
			}
		}
	} else {
		for _, returnCode := range packet.UnsubackPayload {
			if returnCode.IsError() {
				w.err = NewReasonCodeError(returnCode, "")
				return
			}
		}
	}
}
