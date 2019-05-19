package mqtt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacketType(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		b    byte
		want PacketType
	}{
		{0x10, CONNECT},
		{0x20, CONNACK},
		{0x30, PUBLISH},
		{0x40, PUBACK},
		{0x50, PUBREC},
		{0x60, PUBREL},
		{0x70, PUBCOMP},
		{0x80, SUBSCRIBE},
		{0x90, SUBACK},
		{0xa0, UNSUBSCRIBE},
		{0xb0, UNSUBACK},
		{0xc0, PINGREQ},
		{0xd0, PINGRESP},
		{0xe0, DISCONNECT},
		{0xf0, AUTH},
	}

	for _, test := range tests {
		assert.Equal(test.want, (FixedHeader{typeAndFlags: test.b}).PacketType())
	}
}

type headerTestPacket struct {
	FixedHeader
}

func (p headerTestPacket) _isPacket()                     {}
func (p headerTestPacket) PacketType() PacketType         { return p.FixedHeader.PacketType() }
func (p headerTestPacket) fixedHeader(_ byte) FixedHeader { return p.FixedHeader }

func TestReadWriteFixedHeader(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		bin    []byte
		header FixedHeader
	}{
		{[]byte{0x30, 0x00}, FixedHeader{typeAndFlags: 0x30, remainingLength: 0}},
		{[]byte{0x30, 0x7F}, FixedHeader{typeAndFlags: 0x30, remainingLength: 127}},
		{[]byte{0x30, 0x80, 0x01}, FixedHeader{typeAndFlags: 0x30, remainingLength: 128}},
		{[]byte{0x30, 0xFF, 0x7F}, FixedHeader{typeAndFlags: 0x30, remainingLength: 16383}},
		{[]byte{0x30, 0x80, 0x80, 0x01}, FixedHeader{typeAndFlags: 0x30, remainingLength: 16384}},
		{[]byte{0x30, 0xFF, 0xFF, 0x7F}, FixedHeader{typeAndFlags: 0x30, remainingLength: 2097151}},
		{[]byte{0x30, 0x80, 0x80, 0x80, 0x01}, FixedHeader{typeAndFlags: 0x30, remainingLength: 2097152}},
		{[]byte{0x30, 0xFF, 0xFF, 0xFF, 0x7F}, FixedHeader{typeAndFlags: 0x30, remainingLength: 268435455}},
	}

	for _, test := range tests {
		r := NewReader(bytes.NewBuffer(test.bin))
		r.readFixedHeader()
		assert.NoError(r.err)
		assert.Equal(test.header, r.header)

		buf := &bytes.Buffer{}
		w := NewWriter(buf)
		w.packet = headerTestPacket{test.header}
		w.writeFixedHeader()
		assert.NoError(r.err)
		assert.Equal(test.bin, buf.Bytes())
	}
}

func TestValidateFixedHeader(t *testing.T) {
	assert := assert.New(t)

	r := PacketReader{protocol: 4}

	tests := []struct {
		flags byte
		valid bool
	}{
		{0x10, true},
		{0x11, false},
		{0x3D, true},
		{0x82, true},
		{0x80, false},
	}

	for _, test := range tests {
		err := r.validateFixedHeader(FixedHeader{typeAndFlags: test.flags})
		if test.valid {
			assert.NoError(err)
		} else {
			assert.Error(err)
		}
	}
}
