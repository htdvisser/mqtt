package mqtt

import (
	"encoding/binary"
	"fmt"
)

// Properties is a slice of MQTT Properties.
type Properties []Property

func (properties Properties) size() uint32 {
	size := properties.innerSize()
	size += uint32(binary.PutUvarint(make([]byte, 4), uint64(size)))
	return size
}

func (properties Properties) innerSize() (size uint32) {
	for _, property := range properties {
		size += property.size()
	}
	return size
}

// Property is a single MQTT Property
type Property struct {
	Identifier      PropertyIdentifier
	UintValue       uint64 // used for all uints
	BytesValue      []byte // used for bytes and strings
	StringPairValue StringPair
	ByteValue       byte
}

func (p Property) size() uint32 {
	switch p.Identifier {
	// UintValue:
	case SubscriptionIdentifier:
		return 1 + uint32(binary.PutUvarint(make([]byte, 4), p.UintValue))
	// Uint32Value:
	case MessageExpiryInterval,
		SessionExpiryInterval,
		WillDelayInterval,
		MaximumPacketSize:
		return 1 + 4
	// Uint16Value:
	case ServerKeepAlive,
		ReceiveMaximum,
		TopicAliasMaximum,
		TopicAlias:
		return 1 + 2
	// StringValue:
	case ContentType,
		ResponseTopic,
		AssignedClientIdentifier,
		AuthenticationMethod,
		ResponseInformation,
		ServerReference,
		ReasonString:
		return 1 + 2 + uint32(len(p.BytesValue))
	// BytesValue:
	case CorrelationData,
		AuthenticationData:
		return 1 + 2 + uint32(len(p.BytesValue))
	// StringPairValue:
	case UserProperty:
		return 1 + 2 + uint32(len(p.StringPairValue.Key)) + 2 + uint32(len(p.StringPairValue.Value))
	// ByteValue:
	case PayloadFormatIndicator,
		RequestProblemInformation,
		RequestResponseInformation,
		MaximumQoS,
		RetainAvailable,
		WildcardSubscriptionAvailable,
		SubscriptionIdentifierAvailable,
		SharedSubscriptionAvailable:
		return 1 + 1
	default:
		panic(errUnknownProperty)
	}
}

// PropertyIdentifier is the identifier for MQTT properties.
type PropertyIdentifier uint64

// PropertyIdentifier values.
const (
	_                               PropertyIdentifier = 0  // Reserved
	PayloadFormatIndicator                             = 1  // Payload Format Indicator
	MessageExpiryInterval                              = 2  // Message Expiry Interval
	ContentType                                        = 3  // Content Type
	ResponseTopic                                      = 8  // Response Topic
	CorrelationData                                    = 9  // Correlation Data
	SubscriptionIdentifier                             = 11 // Subscription Identifier
	SessionExpiryInterval                              = 17 // Session Expiry Interval
	AssignedClientIdentifier                           = 18 // Assigned Client Identifier
	ServerKeepAlive                                    = 19 // Server Keep Alive
	AuthenticationMethod                               = 21 // Authentication Method
	AuthenticationData                                 = 22 // Authentication Data
	RequestProblemInformation                          = 23 // Request Problem Information
	WillDelayInterval                                  = 24 // Will Delay Interval
	RequestResponseInformation                         = 25 // Request Response Information
	ResponseInformation                                = 26 // Response Information
	ServerReference                                    = 28 // Server Reference
	ReasonString                                       = 31 // Reason String
	ReceiveMaximum                                     = 33 // Receive Maximum
	TopicAliasMaximum                                  = 34 // Topic Alias Maximum
	TopicAlias                                         = 35 // Topic Alias
	MaximumQoS                                         = 36 // Maximum QoS
	RetainAvailable                                    = 37 // Retain Available
	UserProperty                                       = 38 // User Property
	MaximumPacketSize                                  = 39 // Maximum Packet Size
	WildcardSubscriptionAvailable                      = 40 // Wildcard Subscription Available
	SubscriptionIdentifierAvailable                    = 41 // Subscription Identifier Available
	SharedSubscriptionAvailable                        = 42 // Shared Subscription Available
)

const willProperties PacketType = 255

var allowedPropertyIdentifiers = make(map[PropertyIdentifier]map[PacketType]bool)

func allowPropertyIdentifier(id PropertyIdentifier, packetTypes ...PacketType) {
	allowed, ok := allowedPropertyIdentifiers[id]
	if !ok {
		allowed = make(map[PacketType]bool)
		allowedPropertyIdentifiers[id] = allowed
	}
	for _, packetType := range packetTypes {
		allowed[packetType] = true
	}
}

func init() {
	allowPropertyIdentifier(PayloadFormatIndicator, PUBLISH, willProperties)
	allowPropertyIdentifier(MessageExpiryInterval, PUBLISH, willProperties)
	allowPropertyIdentifier(ContentType, PUBLISH, willProperties)
	allowPropertyIdentifier(ResponseTopic, PUBLISH, willProperties)
	allowPropertyIdentifier(CorrelationData, PUBLISH, willProperties)
	allowPropertyIdentifier(SubscriptionIdentifier, PUBLISH, SUBSCRIBE)
	allowPropertyIdentifier(SessionExpiryInterval, CONNECT, CONNACK, DISCONNECT)
	allowPropertyIdentifier(AssignedClientIdentifier, CONNACK)
	allowPropertyIdentifier(ServerKeepAlive, CONNACK)
	allowPropertyIdentifier(AuthenticationMethod, CONNECT, CONNACK, AUTH)
	allowPropertyIdentifier(AuthenticationData, CONNECT, CONNACK, AUTH)
	allowPropertyIdentifier(RequestProblemInformation, CONNECT)
	allowPropertyIdentifier(WillDelayInterval, willProperties)
	allowPropertyIdentifier(RequestResponseInformation, CONNECT)
	allowPropertyIdentifier(ResponseInformation, CONNACK)
	allowPropertyIdentifier(ServerReference, CONNACK, DISCONNECT)
	allowPropertyIdentifier(ReasonString, CONNACK, PUBACK, PUBREC, PUBREL, PUBCOMP, SUBACK, UNSUBACK, DISCONNECT, AUTH)
	allowPropertyIdentifier(ReceiveMaximum, CONNECT, CONNACK)
	allowPropertyIdentifier(TopicAliasMaximum, CONNECT, CONNACK)
	allowPropertyIdentifier(TopicAlias, PUBLISH)
	allowPropertyIdentifier(MaximumQoS, CONNACK)
	allowPropertyIdentifier(RetainAvailable, CONNACK)
	allowPropertyIdentifier(UserProperty, CONNECT, CONNACK, PUBLISH, willProperties, PUBACK, PUBREC, PUBREL, PUBCOMP, SUBSCRIBE, SUBACK, UNSUBSCRIBE, UNSUBACK, DISCONNECT, AUTH)
	allowPropertyIdentifier(MaximumPacketSize, CONNECT, CONNACK)
	allowPropertyIdentifier(WildcardSubscriptionAvailable, CONNACK)
	allowPropertyIdentifier(SubscriptionIdentifierAvailable, CONNACK)
	allowPropertyIdentifier(SharedSubscriptionAvailable, CONNACK)
}

func (r *PacketReader) validateProperties(properties Properties, packetType PacketType) error {
	for _, property := range properties {
		if !allowedPropertyIdentifiers[property.Identifier][packetType] {
			return NewReasonCodeError(MalformedPacket, fmt.Sprintf("mqtt: invalid property %d for %s", property.Identifier, packetType))
		}
	}
	return nil
}

// StringPair is a key-value pair of MQTT strings.
type StringPair struct {
	Key   []byte
	Value []byte
}

// Strings returns the key-value pair as strings.
func (p StringPair) Strings() (key, value string) {
	return string(p.Key), string(p.Value)
}

func (r *PacketReader) readPacketProperties() {
	if r.protocol < 5 || r.remaining() == 0 {
		return
	}
	var properties Properties
	switch pkt := r.packet.(type) {
	case *ConnectPacket:
		properties = r.readProperties()
		pkt.Properties = properties
	case *ConnackPacket:
		properties = r.readProperties()
		pkt.Properties = properties
	case *PublishPacket:
		properties = r.readProperties()
		pkt.Properties = properties
	case *PubackPacket:
		properties = r.readProperties()
		pkt.Properties = properties
	case *PubrecPacket:
		properties = r.readProperties()
		pkt.Properties = properties
	case *PubrelPacket:
		properties = r.readProperties()
		pkt.Properties = properties
	case *PubcompPacket:
		properties = r.readProperties()
		pkt.Properties = properties
	case *SubscribePacket:
		properties = r.readProperties()
		pkt.Properties = properties
	case *SubackPacket:
		properties = r.readProperties()
		pkt.Properties = properties
	case *UnsubscribePacket:
		properties = r.readProperties()
		pkt.Properties = properties
	case *UnsubackPacket:
		properties = r.readProperties()
		pkt.Properties = properties
	case *DisconnectPacket:
		properties = r.readProperties()
		pkt.Properties = properties
	case *AuthPacket:
		properties = r.readProperties()
		pkt.Properties = properties
	default:
		return
	}
	if r.err != nil {
		return
	}
	r.err = r.validateProperties(properties, r.packet.PacketType())
}

func (r *PacketReader) readProperties() Properties {
	var properties Properties
	var propertyLength uint64
	if propertyLength, r.err = r.readUvarint(); r.err != nil {
		return nil
	}
	nReadBefore := r.nRead
	for uint64(r.nRead-nReadBefore) < propertyLength {
		property := r.readProperty()
		if r.err != nil {
			return nil
		}
		properties = append(properties, property)
	}
	return properties
}

var errUnknownProperty = NewReasonCodeError(ProtocolError, "mqtt: unknown property")

func (r *PacketReader) readProperty() (p Property) {
	var id uint64
	if id, r.err = r.readUvarint(); r.err != nil {
		return
	}
	p.Identifier = PropertyIdentifier(id)
	switch p.Identifier {
	case PayloadFormatIndicator,
		RequestProblemInformation,
		RequestResponseInformation,
		MaximumQoS,
		RetainAvailable,
		WildcardSubscriptionAvailable,
		SubscriptionIdentifierAvailable,
		SharedSubscriptionAvailable:
		p.ByteValue, r.err = r.readByte()
	case MessageExpiryInterval,
		SessionExpiryInterval,
		WillDelayInterval,
		MaximumPacketSize:
		var v uint32
		v, r.err = r.readUint32()
		p.UintValue = uint64(v)
	case ContentType,
		ResponseTopic,
		AssignedClientIdentifier,
		AuthenticationMethod,
		ResponseInformation,
		ServerReference,
		ReasonString:
		p.BytesValue, r.err = r.readString()
	case CorrelationData,
		AuthenticationData:
		p.BytesValue, r.err = r.readBytes()
	case SubscriptionIdentifier:
		p.UintValue, r.err = r.readUvarint()
	case ServerKeepAlive,
		ReceiveMaximum,
		TopicAliasMaximum,
		TopicAlias:
		var v uint16
		v, r.err = r.readUint16()
		p.UintValue = uint64(v)
	case UserProperty:
		pair := StringPair{}
		pair.Key, pair.Value, r.err = r.readStringPair()
		p.StringPairValue = pair
	default:
		r.err = errUnknownProperty
	}
	return
}

func (w *PacketWriter) writePacketProperties() {
	if w.protocol < 5 {
		return
	}
	switch pkt := w.packet.(type) {
	case *ConnectPacket:
		w.writeProperties(pkt.Properties)
	case *ConnackPacket:
		w.writeProperties(pkt.Properties)
	case *PublishPacket:
		w.writeProperties(pkt.Properties)
	case *PubackPacket:
		w.writeProperties(pkt.Properties)
	case *PubrecPacket:
		w.writeProperties(pkt.Properties)
	case *PubrelPacket:
		w.writeProperties(pkt.Properties)
	case *PubcompPacket:
		w.writeProperties(pkt.Properties)
	case *SubscribePacket:
		w.writeProperties(pkt.Properties)
	case *SubackPacket:
		w.writeProperties(pkt.Properties)
	case *UnsubscribePacket:
		w.writeProperties(pkt.Properties)
	case *UnsubackPacket:
		w.writeProperties(pkt.Properties)
	case *DisconnectPacket:
		w.writeProperties(pkt.Properties)
	case *AuthPacket:
		w.writeProperties(pkt.Properties)
	default:
		return
	}
}

func (w *PacketWriter) writeProperties(p Properties) {
	if w.err = w.writeUvarint(p.innerSize()); w.err != nil {
		return
	}
	for _, property := range p {
		w.writeProperty(property)
		if w.err != nil {
			return
		}
	}
}

func (w *PacketWriter) writeProperty(p Property) {
	if w.err = w.writeUvarint(uint32(p.Identifier)); w.err != nil {
		return
	}
	switch p.Identifier {
	case PayloadFormatIndicator,
		RequestProblemInformation,
		RequestResponseInformation,
		MaximumQoS,
		RetainAvailable,
		WildcardSubscriptionAvailable,
		SubscriptionIdentifierAvailable,
		SharedSubscriptionAvailable:
		w.err = w.writeByte(p.ByteValue)
	case MessageExpiryInterval,
		SessionExpiryInterval,
		WillDelayInterval,
		MaximumPacketSize:
		w.err = w.writeUint32(uint32(p.UintValue))
	case ContentType,
		ResponseTopic,
		AssignedClientIdentifier,
		AuthenticationMethod,
		ResponseInformation,
		ServerReference,
		ReasonString:
		w.err = w.writeBytes(p.BytesValue)
	case CorrelationData,
		AuthenticationData:
		w.err = w.writeBytes(p.BytesValue)
	case SubscriptionIdentifier:
		w.err = w.writeUvarint(uint32(p.UintValue))
	case ServerKeepAlive,
		ReceiveMaximum,
		TopicAliasMaximum,
		TopicAlias:
		w.err = w.writeUint16(uint16(p.UintValue))
	case UserProperty:
		if w.err = w.writeBytes(p.StringPairValue.Key); w.err != nil {
			return
		}
		w.err = w.writeBytes(p.StringPairValue.Value)
	default:
		w.err = errUnknownProperty
	}
}
