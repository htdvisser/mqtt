package mqtt

// Connack returns an ConnackPacket as response to the ConnectPacket.
func (p *ConnectPacket) Connack() *ConnackPacket {
	var connack ConnackPacket
	return &connack
}

// Reply returns the appropriate reply to this packet.
func (p *ConnectPacket) Reply() Packet { return p.Connack() }

// Puback returns an PubackPacket as response to the PublishPacket.
func (p *PublishPacket) Puback() *PubackPacket {
	var puback PubackPacket
	puback.PacketIdentifier = p.PacketIdentifier
	return &puback
}

// Pubrec returns an PubrecPacket as response to the PublishPacket.
func (p *PublishPacket) Pubrec() *PubrecPacket {
	var pubrec PubrecPacket
	pubrec.PacketIdentifier = p.PacketIdentifier
	return &pubrec
}

// Reply returns the appropriate reply to this packet.
// If no reply is needed, nil is returned.
func (p *PublishPacket) Reply() Packet {
	switch p.QoS() {
	case QoS0:
		return nil
	case QoS1:
		return p.Puback()
	case QoS2:
		return p.Pubrec()
	}
	panic(errInvalidQoS)
}

// Pubrel returns an PubrelPacket as response to the PubrecPacket.
func (p *PubrecPacket) Pubrel() *PubrelPacket {
	var pubrel PubrelPacket
	pubrel.PacketIdentifier = p.PacketIdentifier
	return &pubrel
}

// Reply returns the appropriate reply to this packet.
func (p *PubrecPacket) Reply() Packet { return p.Pubrel() }

// Pubcomp returns an PubcompPacket as response to the PubrelPacket.
func (p *PubrelPacket) Pubcomp() *PubcompPacket {
	var pubcomp PubcompPacket
	pubcomp.PacketIdentifier = p.PacketIdentifier
	return &pubcomp
}

// Reply returns the appropriate reply to this packet.
func (p *PubrelPacket) Reply() Packet { return p.Pubcomp() }

// Suback returns an SubackPacket as response to the SubscribePacket.
func (p *SubscribePacket) Suback() *SubackPacket {
	var suback SubackPacket
	suback.PacketIdentifier = p.PacketIdentifier
	suback.SubackPayload = make([]ReasonCode, len(p.SubscribePayload))
	return &suback
}

// Reply returns the appropriate reply to this packet.
func (p *SubscribePacket) Reply() Packet { return p.Suback() }

// Unsuback returns an UnsubackPacket as response to the UnsubscribePacket.
func (p *UnsubscribePacket) Unsuback() *UnsubackPacket {
	var unsuback UnsubackPacket
	unsuback.PacketIdentifier = p.PacketIdentifier
	unsuback.UnsubackPayload = make([]ReasonCode, len(p.UnsubscribePayload))
	return &unsuback
}

// Reply returns the appropriate reply to this packet.
func (p *UnsubscribePacket) Reply() Packet { return p.Unsuback() }

// Pingresp returns an PingrespPacket as response to the PingreqPacket.
func (p *PingreqPacket) Pingresp() *PingrespPacket {
	var pingresp PingrespPacket
	return &pingresp
}

// Reply returns the appropriate reply to this packet.
func (p *PingreqPacket) Reply() Packet { return p.Pingresp() }

// Auth returns an AuthPacket as response to the AuthPacket.
func (p *AuthPacket) Auth() *AuthPacket {
	var auth AuthPacket
	return &auth
}

// Reply returns the appropriate reply to this packet.
func (p *AuthPacket) Reply() Packet { return p.Auth() }
