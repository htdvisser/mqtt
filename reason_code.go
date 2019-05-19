package mqtt

// ReasonCode indicates the result of an operation
type ReasonCode byte

// ReasonCode Values.
const (
	Success                             ReasonCode = 0x00 // Success
	NormalDisconnection                 ReasonCode = 0x00 // Normal disconnection
	GrantedQoS0                         ReasonCode = 0x00 // Granted QoS 0
	GrantedQoS1                         ReasonCode = 0x01 // Granted QoS 1
	GrantedQoS2                         ReasonCode = 0x02 // Granted QoS 2
	DisconnectWithWillMessage           ReasonCode = 0x04 // Disconnect with Will Message
	NoMatchingSubscribers               ReasonCode = 0x10 // No matching subscribers
	NoSubscriptionExisted               ReasonCode = 0x11 // No subscription existed
	ContinueAuthentication              ReasonCode = 0x18 // Continue authentication
	ReAuthenticate                      ReasonCode = 0x19 // Re-authenticate
	UnspecifiedError                    ReasonCode = 0x80 // Unspecified error
	MalformedPacket                     ReasonCode = 0x81 // Malformed Packet
	ProtocolError                       ReasonCode = 0x82 // Protocol Error
	ImplementationSpecificError         ReasonCode = 0x83 // Implementation specific error
	UnsupportedProtocolVersion          ReasonCode = 0x84 // Unsupported Protocol Version
	ClientIdentifierNotValid            ReasonCode = 0x85 // Client Identifier not valid
	BadUsernameOrPassword               ReasonCode = 0x86 // Bad User Name or Password
	NotAuthorized                       ReasonCode = 0x87 // Not authorized
	ServerUnavailable                   ReasonCode = 0x88 // Server unavailable
	ServerBusy                          ReasonCode = 0x89 // Server busy
	Banned                              ReasonCode = 0x8A // Banned
	ServerShuttingDown                  ReasonCode = 0x8B // Server shutting down
	BadAuthenticationMethod             ReasonCode = 0x8C // Bad authentication method
	KeepAliveTimeout                    ReasonCode = 0x8D // Keep Alive timeout
	SessionTakenOver                    ReasonCode = 0x8E // Session taken over
	TopicFilterInvalid                  ReasonCode = 0x8F // Topic Filter invalid
	TopicNameInvalid                    ReasonCode = 0x90 // Topic Name invalid
	PacketIdentifierInUse               ReasonCode = 0x91 // Packet Identifier in use
	PacketIdentifierNotFound            ReasonCode = 0x92 // Packet Identifier not found
	ReceiveMaximumExceeded              ReasonCode = 0x93 // Receive Maximum exceeded
	TopicAliasInvalid                   ReasonCode = 0x94 // Topic Alias invalid
	PacketTooLarge                      ReasonCode = 0x95 // Packet too large
	MessageRateTooHigh                  ReasonCode = 0x96 // Message rate too high
	QuotaExceeded                       ReasonCode = 0x97 // Quota exceeded
	AdministrativeAction                ReasonCode = 0x98 // Administrative action
	PayloadFormatInvalid                ReasonCode = 0x99 // Payload format invalid
	RetainNotSupported                  ReasonCode = 0x9A // Retain not supported
	QoSNotSupported                     ReasonCode = 0x9B // QoS not supported
	UseAnotherServer                    ReasonCode = 0x9C // Use another server
	ServerMoved                         ReasonCode = 0x9D // Server moved
	SharedSubscriptionsNotSupported     ReasonCode = 0x9E // Shared Subscriptions not supported
	ConnectionRateExceeded              ReasonCode = 0x9F // Connection rate exceeded
	MaximumConnectTime                  ReasonCode = 0xA0 // Maximum connect time
	SubscriptionIdentifiersNotSupported ReasonCode = 0xA1 // Subscription Identifiers not supported
	WildcardSubscriptionsNotSupported   ReasonCode = 0xA2 // Wildcard Subscriptions not supported
)

// IsError returns whether the reason code is an error code.
func (c ReasonCode) IsError() bool { return c >= 0x80 }
