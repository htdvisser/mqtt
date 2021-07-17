// +build gofuzzbeta

package mqtt_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"htdvisser.dev/mqtt"
)

func generatePackets(f *testing.F) {
	pktCh := make(chan mqtt.Packet)
	go func() {
		generateConnects(pktCh)
		generateConnacks(pktCh)
		generatePublishes(pktCh)
		generatePubacks(pktCh)
		generatePubrecs(pktCh)
		generatePubrels(pktCh)
		generatePubcomps(pktCh)
		generateSubscribes(pktCh)
		generateSubacks(pktCh)
		generateUnsubscribes(pktCh)
		generateUnsubacks(pktCh)
		generatePingreqs(pktCh)
		generatePingresps(pktCh)
		generateDisconnects(pktCh)
		generateAuths(pktCh)
		close(pktCh)
	}()

	for pkt := range pktCh {
	nextVersion:
		for _, protocolVersion := range []byte{3, 4, 5} {
			var buf bytes.Buffer

			w := mqtt.NewWriter(&buf)
			w.SetProtocol(protocolVersion)

			switch pkt := pkt.(type) {
			case *mqtt.ConnectPacket:
				pkt.ProtocolVersion = protocolVersion
			case *mqtt.SubackPacket:
				if protocolVersion < 4 {
					for _, code := range pkt.SubackPayload {
						if code.IsError() {
							if err := w.WritePacket(pkt); err == nil {
								f.Error("WritePacket of SUBACK with error did not return error")
							}
							continue nextVersion
						}
					}
				}
			case *mqtt.UnsubackPacket:
				if protocolVersion < 5 {
					for _, code := range pkt.UnsubackPayload {
						if code.IsError() {
							if err := w.WritePacket(pkt); err == nil {
								f.Error("WritePacket of UNSUBACK with error did not return error")
							}
							continue nextVersion
						}
					}
				}
			case *mqtt.AuthPacket:
				if protocolVersion < 5 {
					continue nextVersion
				}
			}

			if err := w.WritePacket(pkt); err != nil {
				f.Errorf("WritePacket(pkt %+v): %v", pkt, err)
			}

			f.Add(buf.Bytes())
		}
	}
}

func generateConnects(pktCh chan<- mqtt.Packet) {
	mods := []func(*mqtt.ConnectPacket){
		func(pkt *mqtt.ConnectPacket) { pkt.SetCleanSession(false); pkt.ClientIdentifier = []byte("client-id") },
		func(pkt *mqtt.ConnectPacket) { pkt.SetUsername([]byte("username")) },
		func(pkt *mqtt.ConnectPacket) { pkt.SetPassword([]byte("password")) },
		func(pkt *mqtt.ConnectPacket) { pkt.SetWill(nil, []byte("will/topic"), []byte("will")) },
		func(pkt *mqtt.ConnectPacket) { pkt.SetWillQoS(mqtt.QoS0) },
		func(pkt *mqtt.ConnectPacket) { pkt.SetWillRetain(true) },
	}
	for i := 0; i < (1 << len(mods)); i++ {
		pkt := new(mqtt.ConnectPacket)
		pkt.SetCleanSession(true)
		for j := 0; j < len(mods); j++ {
			if i&(1<<j) == (1 << j) {
				mods[j](pkt)
			}
		}
		pktCh <- pkt
	}
}

func generateConnacks(pktCh chan<- mqtt.Packet) {
	{
		pkt := new(mqtt.ConnackPacket)
		pkt.SetSessionPresent(true)
		pktCh <- pkt
	}

	{
		pkt := new(mqtt.ConnackPacket)
		pkt.SetSessionPresent(false)
		pktCh <- pkt
	}

	for _, reasonCode := range []mqtt.ReasonCode{
		mqtt.UnspecifiedError,
		mqtt.MalformedPacket,
		mqtt.ProtocolError,
		mqtt.ImplementationSpecificError,
		mqtt.UnsupportedProtocolVersion,
		mqtt.ClientIdentifierNotValid,
		mqtt.BadUsernameOrPassword,
		mqtt.NotAuthorized,
		mqtt.ServerUnavailable,
		mqtt.ServerBusy,
		mqtt.Banned,
		mqtt.BadAuthenticationMethod,
		mqtt.TopicNameInvalid,
		mqtt.PacketTooLarge,
		mqtt.QuotaExceeded,
		mqtt.PayloadFormatInvalid,
		mqtt.RetainNotSupported,
		mqtt.QoSNotSupported,
		mqtt.UseAnotherServer,
		mqtt.ServerMoved,
		mqtt.ConnectionRateExceeded,
	} {
		pkt := new(mqtt.ConnackPacket)
		pkt.ReasonCode = reasonCode
		pktCh <- pkt
	}
}

func generatePublishes(pktCh chan<- mqtt.Packet) {
	mods := []func(*mqtt.PublishPacket){
		func(pkt *mqtt.PublishPacket) { pkt.SetDup(true) },
		func(pkt *mqtt.PublishPacket) { pkt.SetQoS(1); pkt.PacketIdentifier = 1234 },
		func(pkt *mqtt.PublishPacket) { pkt.SetQoS(2); pkt.PacketIdentifier = 1234 },
		func(pkt *mqtt.PublishPacket) { pkt.SetRetain(true) },
	}
	for i := 0; i < (1 << len(mods)); i++ {
		pkt := new(mqtt.PublishPacket)
		pkt.TopicName = []byte("publish/topic")
		pkt.PublishPayload = []byte("publish payload")
		for j := 0; j < len(mods); j++ {
			if i&(1<<j) == (1 << j) {
				mods[j](pkt)
			}
		}
		pktCh <- pkt
	}
}

func generatePubacks(pktCh chan<- mqtt.Packet) {
	for _, reasonCode := range []mqtt.ReasonCode{
		mqtt.Success,
		mqtt.NoMatchingSubscribers,
		mqtt.UnspecifiedError,
		mqtt.ImplementationSpecificError,
		mqtt.NotAuthorized,
		mqtt.TopicNameInvalid,
		mqtt.PacketIdentifierInUse,
		mqtt.QuotaExceeded,
		mqtt.PayloadFormatInvalid,
	} {
		pkt := new(mqtt.PubackPacket)
		pkt.PacketIdentifier = 1234
		pkt.ReasonCode = reasonCode
		pktCh <- pkt
	}
}

func generatePubrecs(pktCh chan<- mqtt.Packet) {
	for _, reasonCode := range []mqtt.ReasonCode{
		mqtt.Success,
		mqtt.NoMatchingSubscribers,
		mqtt.UnspecifiedError,
		mqtt.ImplementationSpecificError,
		mqtt.NotAuthorized,
		mqtt.TopicNameInvalid,
		mqtt.PacketIdentifierInUse,
		mqtt.QuotaExceeded,
		mqtt.PayloadFormatInvalid,
	} {
		pkt := new(mqtt.PubrecPacket)
		pkt.PacketIdentifier = 1234
		pkt.ReasonCode = reasonCode
		pktCh <- pkt
	}
}

func generatePubrels(pktCh chan<- mqtt.Packet) {
	for _, reasonCode := range []mqtt.ReasonCode{
		mqtt.Success,
		mqtt.PacketIdentifierNotFound,
	} {
		pkt := new(mqtt.PubrelPacket)
		pkt.PacketIdentifier = 1234
		pkt.ReasonCode = reasonCode
		pktCh <- pkt
	}
}

func generatePubcomps(pktCh chan<- mqtt.Packet) {
	for _, reasonCode := range []mqtt.ReasonCode{
		mqtt.Success,
		mqtt.PacketIdentifierNotFound,
	} {
		pkt := new(mqtt.PubcompPacket)
		pkt.PacketIdentifier = 1234
		pkt.ReasonCode = reasonCode
		pktCh <- pkt
	}
}

func generateSubscribes(pktCh chan<- mqtt.Packet) {
	mods := []func(*mqtt.SubscribePacket){
		func(pkt *mqtt.SubscribePacket) {
			pkt.SubscribePayload = append(pkt.SubscribePayload, mqtt.Subscription{
				TopicFilter: []byte("topic/qos0"), QoS: mqtt.QoS0,
			})
		},
		func(pkt *mqtt.SubscribePacket) {
			pkt.SubscribePayload = append(pkt.SubscribePayload, mqtt.Subscription{
				TopicFilter: []byte("topic/qos1"), QoS: mqtt.QoS1,
			})
		},
		func(pkt *mqtt.SubscribePacket) {
			pkt.SubscribePayload = append(pkt.SubscribePayload, mqtt.Subscription{
				TopicFilter: []byte("topic/qos2"), QoS: mqtt.QoS2,
			})
		},
	}
	for i := 0; i < (1 << len(mods)); i++ {
		pkt := new(mqtt.SubscribePacket)
		pkt.PacketIdentifier = 1234
		for j := 0; j < len(mods); j++ {
			if i&(1<<j) == (1 << j) {
				mods[j](pkt)
			}
		}
		pktCh <- pkt
	}
}

func generateSubacks(pktCh chan<- mqtt.Packet) {
	for _, reasonCode := range []mqtt.ReasonCode{
		mqtt.GrantedQoS0,
		mqtt.GrantedQoS1,
		mqtt.GrantedQoS2,
		mqtt.UnspecifiedError,
		mqtt.ImplementationSpecificError,
		mqtt.NotAuthorized,
		mqtt.TopicFilterInvalid,
		mqtt.PacketIdentifierInUse,
		mqtt.QuotaExceeded,
		mqtt.SharedSubscriptionsNotSupported,
		mqtt.SubscriptionIdentifiersNotSupported,
		mqtt.WildcardSubscriptionsNotSupported,
	} {
		pkt := new(mqtt.SubackPacket)
		pkt.PacketIdentifier = 1234
		pkt.SubackPayload = []mqtt.ReasonCode{reasonCode}
		pktCh <- pkt
	}
}

func generateUnsubscribes(pktCh chan<- mqtt.Packet) {
	mods := []func(*mqtt.UnsubscribePacket){
		func(pkt *mqtt.UnsubscribePacket) {
			pkt.UnsubscribePayload = append(pkt.UnsubscribePayload, mqtt.TopicFilter("topic/qos0"))
		},
		func(pkt *mqtt.UnsubscribePacket) {
			pkt.UnsubscribePayload = append(pkt.UnsubscribePayload, mqtt.TopicFilter("topic/qos1"))
		},
		func(pkt *mqtt.UnsubscribePacket) {
			pkt.UnsubscribePayload = append(pkt.UnsubscribePayload, mqtt.TopicFilter("topic/qos2"))
		},
	}
	for i := 0; i < (1 << len(mods)); i++ {
		pkt := new(mqtt.UnsubscribePacket)
		pkt.PacketIdentifier = 1234
		for j := 0; j < len(mods); j++ {
			if i&(1<<j) == (1 << j) {
				mods[j](pkt)
			}
		}
		pktCh <- pkt
	}
}

func generateUnsubacks(pktCh chan<- mqtt.Packet) {
	for _, reasonCode := range []mqtt.ReasonCode{
		mqtt.Success,
		mqtt.NoSubscriptionExisted,
		mqtt.UnspecifiedError,
		mqtt.ImplementationSpecificError,
		mqtt.NotAuthorized,
		mqtt.TopicFilterInvalid,
		mqtt.PacketIdentifierInUse,
	} {
		pkt := new(mqtt.UnsubackPacket)
		pkt.PacketIdentifier = 1234
		pkt.UnsubackPayload = []mqtt.ReasonCode{reasonCode}
		pktCh <- pkt
	}
}

func generatePingreqs(pktCh chan<- mqtt.Packet) {
	pktCh <- new(mqtt.PingreqPacket)
}

func generatePingresps(pktCh chan<- mqtt.Packet) {
	pktCh <- new(mqtt.PingrespPacket)
}

func generateDisconnects(pktCh chan<- mqtt.Packet) {
	for _, reasonCode := range []mqtt.ReasonCode{
		mqtt.NormalDisconnection,
		mqtt.DisconnectWithWillMessage,
		mqtt.UnspecifiedError,
		mqtt.MalformedPacket,
		mqtt.ProtocolError,
		mqtt.ImplementationSpecificError,
		mqtt.NotAuthorized,
		mqtt.ServerBusy,
		mqtt.ServerShuttingDown,
		mqtt.KeepAliveTimeout,
		mqtt.SessionTakenOver,
		mqtt.TopicFilterInvalid,
		mqtt.TopicNameInvalid,
		mqtt.ReceiveMaximumExceeded,
		mqtt.TopicAliasInvalid,
		mqtt.PacketTooLarge,
		mqtt.MessageRateTooHigh,
		mqtt.QuotaExceeded,
		mqtt.AdministrativeAction,
		mqtt.PayloadFormatInvalid,
		mqtt.RetainNotSupported,
		mqtt.QoSNotSupported,
		mqtt.UseAnotherServer,
		mqtt.ServerMoved,
		mqtt.SharedSubscriptionsNotSupported,
		mqtt.ConnectionRateExceeded,
		mqtt.MaximumConnectTime,
		mqtt.SubscriptionIdentifiersNotSupported,
		mqtt.WildcardSubscriptionsNotSupported,
	} {
		pkt := new(mqtt.DisconnectPacket)
		pkt.ReasonCode = reasonCode
		pktCh <- pkt
	}
}

func generateAuths(pktCh chan<- mqtt.Packet) {
	for _, reasonCode := range []mqtt.ReasonCode{
		mqtt.Success,
		mqtt.ContinueAuthentication,
		mqtt.ReAuthenticate,
	} {
		pkt := new(mqtt.AuthPacket)
		pkt.ReasonCode = reasonCode
		pktCh <- pkt
	}
}

func FuzzReadWrite(f *testing.F) {
	generatePackets(f)

	f.Fuzz(func(t *testing.T, data []byte) {
		protocolVersion := mqtt.DefaultProtocolVersion

		r := mqtt.NewReader(bytes.NewBuffer(data), mqtt.WithMaxPacketLength(8*1024))
		r.SetProtocol(protocolVersion)

		pkt0, err := r.ReadPacket()
		if err != nil {
			if pkt0 != nil {
				t.Fatal("pkt0 != nil on error")
			}
			t.Skip()
		}

		if connectPkt, ok := pkt0.(*mqtt.ConnectPacket); ok {
			protocolVersion = connectPkt.ProtocolVersion
		}

		var buf bytes.Buffer

		w := mqtt.NewWriter(&buf)
		w.SetProtocol(protocolVersion)
		err = w.WritePacket(pkt0)
		if err != nil {
			t.Fatal(err)
		}

		r2 := mqtt.NewReader(&buf)
		pkt1, err := r2.ReadPacket()
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(pkt0, pkt1); diff != "" {
			t.Fatalf("pkt0 is not equal to pkt1: %v", diff)
		}
	})
}
