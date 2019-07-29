package mqtt_test

import (
	"log"
	"net"
	"time"

	"htdvisser.dev/mqtt"
)

func Example_client() {
	conn, err := net.Dial("tcp", "localhost:1883")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader, writer := mqtt.NewReader(conn), mqtt.NewWriter(conn)

	connect := new(mqtt.ConnectPacket)
	connect.SetCleanSession(true)
	connect.SetUsername([]byte("username"))
	connect.SetPassword([]byte("password"))

	if err = writer.WritePacket(connect); err != nil {
		log.Fatal(err)
	}

	packet, err := reader.ReadPacket()
	if err != nil {
		log.Fatal(err)
	}
	if packet.PacketType() != mqtt.CONNACK {
		log.Fatal("first packet was not CONNACK")
	}
	connack := packet.(*mqtt.ConnackPacket)

	if connack.ReasonCode != mqtt.Success {
		// TODO: Check actual reason code
		log.Fatal("connect failed")
	}

	go func() {
		subscribe := new(mqtt.SubscribePacket)
		subscribe.PacketIdentifier = 1 // TODO: Keep track of these.
		subscribe.SubscribePayload = append(subscribe.SubscribePayload, mqtt.Subscription{
			TopicFilter: []byte("time/+"),
			QoS:         1,
		})
		if err := writer.WritePacket(subscribe); err != nil {
			log.Fatal(err)
		}

		for t := range time.Tick(time.Second) {
			publish := new(mqtt.PublishPacket)
			publish.SetRetain(true)
			publish.TopicName = []byte("time/now")
			publish.PublishPayload = []byte(t.Format(time.RFC3339))
			if err = writer.WritePacket(publish); err != nil {
				log.Fatal(err)
			}
		}
	}()

	for {
		packet, err := reader.ReadPacket()
		if err != nil {
			log.Fatal(err)
		}

		switch packet.PacketType() {
		case mqtt.PUBLISH:
			publish := packet.(*mqtt.PublishPacket)
			switch publish.QoS() {
			case mqtt.QoS1:
				puback := publish.Puback()
				// TODO: Handle QoS 1 publish
				err = writer.WritePacket(puback)
			case mqtt.QoS2:
				pubrec := publish.Pubrec()
				// TODO: Handle QoS 2 publish
				err = writer.WritePacket(pubrec)
			}
		case mqtt.PUBACK:
			puback := packet.(*mqtt.PubackPacket)
			// TODO: Handle QoS 1 publish
			_ = puback.PacketIdentifier
		case mqtt.PUBREC:
			pubrec := packet.(*mqtt.PubrecPacket)
			pubrel := pubrec.Pubrel()
			// TODO: Handle QoS 2 publish
			err = writer.WritePacket(pubrel)
		case mqtt.PUBREL:
			pubrel := packet.(*mqtt.PubrelPacket)
			pubcomp := pubrel.Pubcomp()
			// TODO: Handle QoS 2 publish
			err = writer.WritePacket(pubcomp)
		case mqtt.PUBCOMP:
			pubcomp := packet.(*mqtt.PubcompPacket)
			// TODO: Handle QoS 2 publish
			_ = pubcomp.PacketIdentifier
		case mqtt.SUBACK:
			suback := packet.(*mqtt.SubackPacket)
			// TODO: Find subscribe by ID
			_ = suback.PacketIdentifier
			// TODO: Handle reason codes for subscribes
			_ = suback.SubackPayload
		case mqtt.UNSUBACK:
			unsuback := packet.(*mqtt.UnsubackPacket)
			// TODO: Find unsubscribe by ID
			_ = unsuback.PacketIdentifier
			// TODO: Handle reason codes for unsubscribes
			_ = unsuback.UnsubackPayload
		case mqtt.PINGRESP:
			// TODO: Handle pingresp
		case mqtt.DISCONNECT:
			disconnect := packet.(*mqtt.DisconnectPacket)
			_ = disconnect.ReasonCode
			// TODO: Handle disconnect
			return
		case mqtt.AUTH:
			auth := packet.(*mqtt.AuthPacket)
			if auth.ReasonCode != mqtt.ContinueAuthentication {
				log.Fatal("received auth packet with invalid reason code")
			}
			// TODO: Handle auth
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}
