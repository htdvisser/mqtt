package mqtt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePublishFlags(t *testing.T) {
	assert := assert.New(t)

	r := PacketReader{protocol: 4}

	tests := []struct {
		flags PublishFlags
		valid bool
	}{
		{0x0, true},
		{0x1, true},
		{0x2, true},
		{0x3, true},
		{0x4, true},
		{0x5, true},
		{0x6, false},
		{0x7, false},
		{0x8, true},
		{0x9, true},
		{0xa, true},
		{0xb, true},
		{0xc, true},
		{0xd, true},
		{0xe, false},
		{0xf, false},
	}

	for _, test := range tests {
		err := r.validatePublishFlags(test.flags)
		if test.valid {
			assert.NoError(err)
		} else {
			assert.Error(err)
		}
	}
}

func TestPublishFlags(t *testing.T) {
	assert := assert.New(t)
	var f PublishFlags

	assert.False(f.Dup())
	assert.Equal(QoS0, f.QoS())
	assert.False(f.Retain())

	f.SetDup(true)
	f.SetQoS(QoS1)
	f.SetRetain(true)

	assert.True(f.Dup())
	assert.Equal(QoS1, f.QoS())
	assert.True(f.Retain())

	f.SetDup(false)
	f.SetQoS(QoS2)
	f.SetRetain(false)

	assert.False(f.Dup())
	assert.Equal(QoS2, f.QoS())
	assert.False(f.Retain())
}
