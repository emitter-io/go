package emitter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	assert.Equal(t, "me=0", WithoutEcho().String())
	assert.Equal(t, "ttl=5", WithTTL(5).String())
	assert.Equal(t, "last=5", WithLast(5).String())
}
