package mqtt_test

import (
	"errors"
	"log"
	"net"

	"htdvisser.dev/mqtt"
)

func ListenAndAccept(address string, handle func(net.Conn) error) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		go func() {
			defer conn.Close()
			handle(conn)
		}()
	}
}

func checkAuth(username, password []byte) error {
	return nil
}

func Example_server() {
	ListenAndAccept("localhost:1883", func(conn net.Conn) error {
		reader, writer := mqtt.NewReader(conn), mqtt.NewWriter(conn)

		packet, err := reader.ReadPacket()
		if err != nil {
			return err
		}
		if packet.PacketType() != mqtt.CONNECT {
			return errors.New("first packet was not CONNECT")
		}
		connect := packet.(*mqtt.ConnectPacket)
		writer.SetProtocol(connect.ProtocolVersion)

		connack := connect.Connack()

		if err := checkAuth(connect.Username(), connect.Password()); err != nil {
			connack.ReasonCode = mqtt.BadUsernameOrPassword
			return writer.WritePacket(connack)
		}

		// TODO: Handle keepalive, session, will, ...

		if err = writer.WritePacket(connack); err != nil {
			return err
		}

		go func() {
			// TODO: Handle publishes to client
		}()

		for {
			packet, err := reader.ReadPacket()
			if err != nil {
				return err
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
			case mqtt.SUBSCRIBE:
				subscribe := packet.(*mqtt.SubscribePacket)
				suback := subscribe.Suback()
				// TODO: Handle subscribe, set reason codes in suback
				err = writer.WritePacket(suback)
			case mqtt.UNSUBSCRIBE:
				unsubscribe := packet.(*mqtt.UnsubscribePacket)
				unsuback := unsubscribe.Unsuback()
				// TODO: Handle unsubscribe, set reason codes in unsuback
				err = writer.WritePacket(unsuback)
			case mqtt.PINGREQ:
				pingreq := packet.(*mqtt.PingreqPacket)
				pingresp := pingreq.Pingresp()
				err = writer.WritePacket(pingresp)
			case mqtt.DISCONNECT:
				disconnect := packet.(*mqtt.DisconnectPacket)
				_ = disconnect.ReasonCode
				// TODO: Handle disconnect
				return nil
			case mqtt.AUTH:
				auth := packet.(*mqtt.AuthPacket)
				if auth.ReasonCode != mqtt.ContinueAuthentication {
					return errors.New("received auth packet with invalid reason code")
				}
				// TODO: Handle auth
			}
			if err != nil {
				return err
			}
		}
	})
}
