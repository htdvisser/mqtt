package mqtt

import (
	"bytes"
	"errors"
)

// ConnectPacket is the Connect packet.
type ConnectPacket struct {
	ConnectHeader
	ConnectPayload
	Properties
}

func (*ConnectPacket) _isPacket() {}

// PacketType returns the packet type of the Connect packet.
func (ConnectPacket) PacketType() PacketType { return CONNECT }

func (p ConnectPacket) fixedHeader(protocol byte) (h FixedHeader) {
	h.SetPacketType(CONNECT)
	h.remainingLength = p.size(protocol)
	return
}

func (p ConnectPacket) size(protocol byte) uint32 {
	size := 2 + len(p.ConnectHeader.ProtocolName) + 1 + 1 + 2
	size += 2 + len(p.ConnectPayload.ClientIdentifier)
	if p.ConnectHeader.Will() {
		size += 2 + len(p.ConnectPayload.WillTopic)
		size += 2 + len(p.ConnectPayload.WillMessage)
		if protocol >= 5 {
			size += int(p.ConnectPayload.WillProperties.size())
		}
	}
	if p.ConnectHeader.Username() {
		size += 2 + len(p.ConnectPayload.Username)
	}
	if p.ConnectHeader.Password() {
		size += 2 + len(p.ConnectPayload.Password)
	}
	if protocol >= 5 {
		size += int(p.Properties.size())
	}
	return uint32(size)
}

// Username returns the Username from the packet or nil if it doesn't have one.
func (p ConnectPacket) Username() []byte {
	if p.ConnectHeader.Username() {
		return p.ConnectPayload.Username
	}
	return nil
}

// SetUsername sets the Username into the packet if non-nil, or unsets it if nil.
func (p *ConnectPacket) SetUsername(username []byte) {
	p.ConnectHeader.SetUsername(username != nil)
	p.ConnectPayload.Username = username
}

// Password returns the Password from the packet or nil if it doesn't have one.
func (p ConnectPacket) Password() []byte {
	if p.ConnectHeader.Password() {
		return p.ConnectPayload.Password
	}
	return nil
}

// SetPassword sets the Password into the packet if non-nil, or unsets it if nil.
func (p *ConnectPacket) SetPassword(password []byte) {
	p.ConnectHeader.SetPassword(password != nil)
	p.ConnectPayload.Password = password
}

// Will returns the Will from the packet or nil if it doesn't have one.
func (p *ConnectPacket) Will() (properties Properties, topic, message []byte) {
	if p.ConnectHeader.Will() {
		return p.ConnectPayload.WillProperties, p.ConnectPayload.WillTopic, p.ConnectPayload.WillMessage
	}
	return nil, nil, nil
}

// SetWill sets the Will into the packet if the topic is non-nil, or unsets it if nil.
func (p *ConnectPacket) SetWill(properties Properties, topic, message []byte) {
	p.ConnectHeader.SetWill(topic != nil)
	p.ConnectPayload.WillProperties = properties
	p.ConnectPayload.WillTopic = topic
	p.ConnectPayload.WillMessage = message
}

// ConnectHeader is the header of the Connect packet.
type ConnectHeader struct {
	ProtocolName    []byte
	ProtocolVersion byte
	ConnectHeaderFlags
	KeepAlive uint16
}

// ConnectHeaderFlags are the flags in the header of the Connect packet.
type ConnectHeaderFlags byte

var errInvalidConnectHeaderFlags = errors.New("mqtt: invalid connect header flags")

func (r *PacketReader) validateConnectHeaderFlags(f ConnectHeaderFlags) error {
	if f&0x18 == 0x18 {
		return errInvalidQoS
	}
	if f&0x01 == 0x01 {
		return errInvalidConnectHeaderFlags
	}
	return nil
}

// Username returns the Username bit from the connect header flags.
func (f ConnectHeaderFlags) Username() bool { return f&0x80 == 0x80 }

// SetUsername sets the Username bit into the connect header flags.
func (f *ConnectHeaderFlags) SetUsername(username bool) {
	*f &^= 0x80
	if username {
		*f |= 0x80
	}
}

// Password returns the Password bit from the connect header flags.
func (f ConnectHeaderFlags) Password() bool { return f&0x40 == 0x40 }

// SetPassword sets the Password bit into the connect header flags.
func (f *ConnectHeaderFlags) SetPassword(password bool) {
	*f &^= 0x40
	if password {
		*f |= 0x40
	}
}

// WillRetain returns the WillRetain bit from the connect header flags.
func (f ConnectHeaderFlags) WillRetain() bool { return f&0x20 == 0x20 }

// SetWillRetain sets the WillRetain bit into the connect header flags.
func (f *ConnectHeaderFlags) SetWillRetain(willRetain bool) {
	*f &^= 0x20
	if willRetain {
		*f |= 0x20
	}
}

// Will returns the Will bit from the connect header flags.
func (f ConnectHeaderFlags) Will() bool { return f&0x04 == 0x04 }

// SetWill sets the Will bit into the connect header flags.
func (f *ConnectHeaderFlags) SetWill(will bool) {
	*f &^= 0x04
	if will {
		*f |= 0x04
	}
}

// CleanStart returns the CleanStart bit from the connect header flags.
func (f ConnectHeaderFlags) CleanStart() bool { return f&0x02 == 0x02 }

// CleanSession is an alias for CleanStart.
func (f ConnectHeaderFlags) CleanSession() bool {
	return f.CleanStart()
}

// SetCleanStart sets the CleanStart bit into the connect header flags.
func (f *ConnectHeaderFlags) SetCleanStart(cleanSession bool) {
	*f &^= 0x02
	if cleanSession {
		*f |= 0x02
	}
}

// SetCleanSession is an alias for SetCleanStart.
func (f *ConnectHeaderFlags) SetCleanSession(cleanSession bool) {
	f.SetCleanStart(cleanSession)
}

// WillQoS returns the WillQoS from the connect header flags.
func (f ConnectHeaderFlags) WillQoS() QoS { return QoS(f >> 3 & 0x03) }

// SetWillQoS sets the WillQoS into the connect header flags.
func (f *ConnectHeaderFlags) SetWillQoS(qos QoS) {
	*f &^= 0x18
	switch qos {
	case 1:
		*f |= 0x08
	case 2:
		*f |= 0x10
	}
}

var (
	protocolMQIsdp = []byte("MQIsdp")
	protocolMQTT   = []byte("MQTT")
)

var (
	errUnknownProtocolName        = errors.New("mqtt: unknown protocol name")
	errUnsupportedProtocolVersion = errors.New("mqtt: unsupported protocol version")
)

func (r *PacketReader) readConnectHeader() {
	packet := r.packet.(*ConnectPacket)
	packet.ConnectHeader.ProtocolName, r.err = r.readBytes()
	if r.err != nil {
		return
	}
	switch {
	case bytes.Equal(packet.ConnectHeader.ProtocolName, protocolMQIsdp):
		packet.ConnectHeader.ProtocolName = protocolMQIsdp
	case bytes.Equal(packet.ConnectHeader.ProtocolName, protocolMQTT):
		packet.ConnectHeader.ProtocolName = protocolMQTT
	default:
		r.err = errUnknownProtocolName
		return
	}
	packet.ConnectHeader.ProtocolVersion, r.err = r.readByte()
	if r.err != nil {
		return
	}
	switch packet.ConnectHeader.ProtocolVersion {
	case 3, 4, 5:
	default:
		r.err = errUnsupportedProtocolVersion
		return
	}
	var f byte
	f, r.err = r.readByte()
	if r.err != nil {
		return
	}
	packet.ConnectHeader.ConnectHeaderFlags = ConnectHeaderFlags(f)
	if r.err = r.validateConnectHeaderFlags(packet.ConnectHeader.ConnectHeaderFlags); r.err != nil {
		return
	}
	packet.ConnectHeader.KeepAlive, r.err = r.readUint16()
	if r.err != nil {
		return
	}
}

func (w *PacketWriter) writeConnectHeader() {
	packet := w.packet.(*ConnectPacket)
	protocolName := packet.ConnectHeader.ProtocolName
	if len(protocolName) == 0 {
		switch w.protocol {
		case 3:
			protocolName = protocolMQIsdp
		case 4, 5:
			protocolName = protocolMQTT
		default:
			w.err = errUnsupportedProtocolVersion
			return
		}
	}
	w.err = w.writeBytes(protocolName)
	if w.err != nil {
		return
	}
	protocolVersion := packet.ConnectHeader.ProtocolVersion
	if protocolVersion == 0 {
		protocolVersion = w.protocol
	}
	w.err = w.writeByte(protocolVersion)
	if w.err != nil {
		return
	}
	w.err = w.writeByte(byte(packet.ConnectHeader.ConnectHeaderFlags))
	if w.err != nil {
		return
	}
	w.err = w.writeUint16(packet.ConnectHeader.KeepAlive)
	if w.err != nil {
		return
	}
}

// ConnectPayload is the payload of the Connect packet.
type ConnectPayload struct {
	ClientIdentifier []byte
	WillProperties   Properties
	WillTopic        []byte
	WillMessage      []byte
	Username         []byte
	Password         []byte
}

var errEmptyClientIdentifier = errors.New("mqtt: empty client identifier")

func (r *PacketReader) readConnectPayload() {
	packet := r.packet.(*ConnectPacket)
	packet.ConnectPayload.ClientIdentifier, r.err = r.readBytes()
	if r.err != nil {
		return
	}
	if r.protocol < 5 && len(packet.ConnectPayload.ClientIdentifier) == 0 && !packet.ConnectHeader.CleanSession() {
		r.err = errEmptyClientIdentifier
		return
	}
	if packet.ConnectHeader.Will() {
		if r.protocol >= 5 {
			packet.ConnectPayload.WillProperties = r.readProperties()
			if r.err != nil {
				return
			}
			r.err = r.validateProperties(packet.ConnectPayload.WillProperties, willProperties)
			if r.err != nil {
				return
			}
		}
		packet.ConnectPayload.WillTopic, r.err = r.readBytes()
		if r.err != nil {
			return
		}
		packet.ConnectPayload.WillMessage, r.err = r.readBytes()
		if r.err != nil {
			return
		}
	}
	if packet.ConnectHeader.Username() {
		packet.ConnectPayload.Username, r.err = r.readBytes()
		if r.err != nil {
			return
		}
	}
	if packet.ConnectHeader.Password() {
		packet.ConnectPayload.Password, r.err = r.readBytes()
		if r.err != nil {
			return
		}
	}
}

func (w *PacketWriter) writeConnectPayload() {
	packet := w.packet.(*ConnectPacket)
	w.err = w.writeBytes(packet.ConnectPayload.ClientIdentifier)
	if w.err != nil {
		return
	}
	if len(packet.ConnectPayload.ClientIdentifier) == 0 && !packet.ConnectHeader.CleanSession() {
		w.err = errEmptyClientIdentifier
		return
	}
	if packet.ConnectHeader.Will() {
		if w.protocol >= 5 {
			w.writeProperties(packet.ConnectPayload.WillProperties)
			if w.err != nil {
				return
			}
		}
		w.err = w.writeBytes(packet.ConnectPayload.WillTopic)
		if w.err != nil {
			return
		}
		w.err = w.writeBytes(packet.ConnectPayload.WillMessage)
		if w.err != nil {
			return
		}
	}
	if packet.ConnectHeader.Username() {
		w.err = w.writeBytes(packet.ConnectPayload.Username)
		if w.err != nil {
			return
		}
	}
	if packet.ConnectHeader.Password() {
		w.err = w.writeBytes(packet.ConnectPayload.Password)
		if w.err != nil {
			return
		}
	}
}
