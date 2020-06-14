package mqtt

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWriteProperties(t *testing.T) {
	for id := range allowedPropertyIdentifiers {
		t.Run(fmt.Sprint(id), func(t *testing.T) {
			assert := assert.New(t)

			buf := &bytes.Buffer{}
			w := NewWriter(buf)

			subject := Property{Identifier: id}

			size := subject.size()

			w.writeProperty(subject)
			assert.NoError(w.err)

			assert.Equal(int(size), buf.Len())

			buf = bytes.NewBuffer(buf.Bytes())

			r := NewReader(buf)
			r.header.remainingLength = w.nWritten

			p := r.readProperty()
			assert.NoError(r.err)
			assert.Equal(subject, p)
		})
	}

	var properties Properties
	for id := range allowedPropertyIdentifiers {
		properties = append(properties, Property{Identifier: id})
	}

	assert := assert.New(t)

	buf := &bytes.Buffer{}
	w := NewWriter(buf)

	w.writeProperties(properties)
	assert.NoError(w.err)

	buf = bytes.NewBuffer(buf.Bytes())

	r := NewReader(buf)
	r.header.remainingLength = w.nWritten

	p := r.readProperties()
	assert.NoError(r.err)
	assert.Equal(properties, p)
}
