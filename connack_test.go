package mqtt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateConnackHeaderFlags(t *testing.T) {
	assert := assert.New(t)

	r := PacketReader{protocol: 4}

	tests := []struct {
		flags ConnackHeaderFlags
		valid bool
	}{
		{0x01, true},
		{0x02, false},
	}

	for _, test := range tests {
		err := r.validateConnackHeaderFlags(test.flags)
		if test.valid {
			assert.NoError(err)
		} else {
			assert.Error(err)
		}
	}
}

func TestConnackHeaderFlags(t *testing.T) {
	assert := assert.New(t)

	var f ConnackHeaderFlags

	assert.False(f.SessionPresent())

	f.SetSessionPresent(true)

	assert.True(f.SessionPresent())

	f.SetSessionPresent(false)

	assert.False(f.SessionPresent())
}
