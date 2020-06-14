package mqtt_test

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"

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
		ctx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup
		defer func() {
			cancel()
			wg.Wait()
		}()

		reader, writer := mqtt.NewReader(conn), mqtt.NewWriter(conn)

		timeout := 10 * time.Second

		conn.SetReadDeadline(time.Now().Add(timeout)) // Read deadline for CONNECT.

		packet, err := reader.ReadPacket()
		if err != nil {
			return err
		}
		if packet.PacketType() != mqtt.CONNECT {
			return errors.New("first packet was not CONNECT")
		}
		connect := packet.(*mqtt.ConnectPacket)
		writer.SetProtocol(connect.ProtocolVersion)

		conn.SetWriteDeadline(time.Now().Add(timeout)) // Write deadline for CONNACK.

		connack := connect.Connack()

		if err := checkAuth(connect.Username(), connect.Password()); err != nil {
			connack.ReasonCode = mqtt.BadUsernameOrPassword
			return writer.WritePacket(connack)
		}

		if connect.KeepAlive != 0 {
			timeout = time.Duration(connect.KeepAlive) * 1500 * time.Millisecond
		}

		// TODO: Handle session, will, ...

		if err = writer.WritePacket(connack); err != nil {
			return err
		}

		conn.SetReadDeadline(time.Now().Add(timeout))
		conn.SetWriteDeadline(time.Time{}) // Clear write deadline.

		var (
			controlPackets = make(chan mqtt.Packet)
			publishPackets = make(chan *mqtt.PublishPacket)
		)

		wg.Add(1)
		go func() { // Write routine
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case pkt := <-controlPackets:
					conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
					err = writer.WritePacket(pkt)
				case pkt := <-publishPackets:
					conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
					err = writer.WritePacket(pkt)
				}
				if err != nil {
					// TODO: Handle error
				}
			}
		}()

		for { // Read routine
			conn.SetReadDeadline(time.Now().Add(timeout))
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
					controlPackets <- puback
				case mqtt.QoS2:
					pubrec := publish.Pubrec()
					// TODO: Handle QoS 2 publish
					controlPackets <- pubrec
				}
			case mqtt.PUBACK:
				puback := packet.(*mqtt.PubackPacket)
				// TODO: Handle QoS 1 publish
				_ = puback.PacketIdentifier
			case mqtt.PUBREC:
				pubrec := packet.(*mqtt.PubrecPacket)
				pubrel := pubrec.Pubrel()
				// TODO: Handle QoS 2 publish
				controlPackets <- pubrel
			case mqtt.PUBREL:
				pubrel := packet.(*mqtt.PubrelPacket)
				pubcomp := pubrel.Pubcomp()
				// TODO: Handle QoS 2 publish
				controlPackets <- pubcomp
			case mqtt.PUBCOMP:
				pubcomp := packet.(*mqtt.PubcompPacket)
				// TODO: Handle QoS 2 publish
				_ = pubcomp.PacketIdentifier
			case mqtt.SUBSCRIBE:
				subscribe := packet.(*mqtt.SubscribePacket)
				suback := subscribe.Suback()
				// TODO: Handle subscribe, set reason codes in suback
				controlPackets <- suback
			case mqtt.UNSUBSCRIBE:
				unsubscribe := packet.(*mqtt.UnsubscribePacket)
				unsuback := unsubscribe.Unsuback()
				// TODO: Handle unsubscribe, set reason codes in unsuback
				controlPackets <- unsuback
			case mqtt.PINGREQ:
				pingreq := packet.(*mqtt.PingreqPacket)
				pingresp := pingreq.Pingresp()
				controlPackets <- pingresp
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
