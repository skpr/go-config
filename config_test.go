package skprconfig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigGet(t *testing.T) {
	options := func(c *Config) {
		c.path = "fixtures/config.json"
	}
	config, err := Load(options)

	assert.Nil(t, err)

	// Assert string.
	value, ok := config.Get("chip.shop")
	assert.True(t, ok)
	assert.Equal(t, "snax", value)

	// Assert int.
	intValue, ok := config.GetInt("foo.bar")
	assert.True(t, ok)
	assert.Equal(t, 123, intValue)

	// Assert nil value.
	value, ok = config.Get("somewhat.secret")
	assert.True(t, ok)
	assert.Empty(t, value)

	// Assert missing key.
	_, ok = config.Get("fish.shop")
	assert.False(t, ok)

	// Assert missing key with fallback.
	assert.Equal(t, "yums", config.GetWithFallback("fish.shop", "yums"))

}
