package mqtt

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWrite(t *testing.T) {
	tests := []struct {
		packet Packet
	}{
		{&ConnectPacket{
			ConnectHeader: ConnectHeader{
				ProtocolName:    []byte("MQTT"),
				ProtocolVersion: 4,
			},
			ConnectPayload: ConnectPayload{
				ClientIdentifier: []byte("foo"),
			},
		}},
		{&ConnackPacket{}},
		{&PublishPacket{
			PublishHeader: PublishHeader{
				TopicName: []byte("foo"),
			},
			PublishPayload: []byte("foo"),
		}},
		{&PubackPacket{}},
		{&PubrecPacket{}},
		{&PubrelPacket{}},
		{&PubcompPacket{}},
		{&SubscribePacket{
			SubscribePayload: []Subscription{
				{[]byte("foo"), QoS2},
			},
		}},
		{&SubackPacket{
			SubackPayload: []ReasonCode{
				ReasonCode(QoS2),
			},
		}},
		{&UnsubscribePacket{
			UnsubscribePayload: []TopicFilter{
				[]byte("foo"),
			},
		}},
		{&UnsubackPacket{}},
		{&PingreqPacket{}},
		{&PingrespPacket{}},
		{&DisconnectPacket{}},
		{&AuthPacket{}},
	}

	for _, protocol := range []byte{3, 4, 5} {
		for _, test := range tests {
			t.Run(fmt.Sprintf("MQTT%d_%T", protocol, test.packet), func(t *testing.T) {
				assert := assert.New(t)

				buf := &bytes.Buffer{}
				w := NewWriter(buf)
				w.SetProtocol(protocol)

				testPacket := test.packet
				if connectPacket, ok := testPacket.(*ConnectPacket); ok {
					connectPacket.ConnectHeader.ProtocolVersion = protocol
					connectPacket.SetUsername([]byte("username"))
					connectPacket.SetPassword([]byte("password"))
					var props Properties
					if protocol >= 5 {
						props = Properties{
							{Identifier: WillDelayInterval, UintValue: 10},
						}
					}
					connectPacket.SetWill(props, []byte("will-topic"), []byte("will-message"))
					testPacket = connectPacket
				}

				err := w.WritePacket(testPacket)
				assert.NoError(err)

				t.Logf("%x", buf.Bytes())

				buf = bytes.NewBuffer(buf.Bytes())

				r := NewReader(buf)
				r.SetProtocol(protocol)

				pkt, err := r.ReadPacket()

				if _, ok := testPacket.(*AuthPacket); ok && protocol < 5 {
					assert.Error(err)
					return
				}

				assert.NoError(err)

				assert.Equal(testPacket, pkt)

				assert.Equal(0, buf.Len())
			})
		}
	}
}
