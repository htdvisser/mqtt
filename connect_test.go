package mqtt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateConnectHeaderFlags(t *testing.T) {
	assert := assert.New(t)

	r := PacketReader{protocol: 4}

	tests := []struct {
		flags ConnectHeaderFlags
		valid bool
	}{
		{0x01, false},
		{0x02, true},
		{0x18, false},
	}

	for _, test := range tests {
		err := r.validateConnectHeaderFlags(test.flags)
		if test.valid {
			assert.NoError(err)
		} else {
			assert.Error(err)
		}
	}
}

func TestConnectHeaderFlags(t *testing.T) {
	assert := assert.New(t)

	var f ConnectHeaderFlags

	assert.False(f.Username())
	assert.False(f.Password())
	assert.False(f.WillRetain())
	assert.Equal(QoS0, f.WillQoS())
	assert.False(f.Will())
	assert.False(f.CleanSession())

	f.SetUsername(true)
	f.SetPassword(true)
	f.SetWillRetain(true)
	f.SetWillQoS(QoS1)
	f.SetWill(true)
	f.SetCleanSession(true)

	assert.True(f.Username())
	assert.True(f.Password())
	assert.True(f.WillRetain())
	assert.Equal(QoS1, f.WillQoS())
	assert.True(f.Will())
	assert.True(f.CleanSession())

	f.SetUsername(false)
	f.SetPassword(false)
	f.SetWillRetain(false)
	f.SetWillQoS(QoS2)
	f.SetWill(false)
	f.SetCleanSession(false)

	assert.False(f.Username())
	assert.False(f.Password())
	assert.False(f.WillRetain())
	assert.Equal(QoS2, f.WillQoS())
	assert.False(f.Will())
	assert.False(f.CleanSession())
}

func TestConnectPacket(t *testing.T) {
	assert := assert.New(t)

	var p ConnectPacket

	assert.Nil(p.Username())
	assert.Nil(p.Password())
	properties, topic, message := p.Will()
	assert.Nil(properties)
	assert.Nil(topic)
	assert.Nil(message)

	p.SetUsername([]byte("username"))
	p.SetPassword([]byte("password"))
	p.SetWill(nil, []byte("will-topic"), []byte("will-message"))

	assert.True(p.ConnectHeader.Username())
	assert.True(p.ConnectHeader.Password())
	assert.True(p.ConnectHeader.Will())

	assert.Equal([]byte("username"), p.Username())
	assert.Equal([]byte("password"), p.Password())
	_, topic, message = p.Will()
	assert.Equal([]byte("will-topic"), topic)
	assert.Equal([]byte("will-message"), message)
}
