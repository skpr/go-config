package skprconfig

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigGet(t *testing.T) {
	options := func(c *Config) {
		c.path = "fixtures/config.json"
	}
	config, err := NewConfig(options)

	assert.Nil(t, err)

	// Assert string
	assert.Equal(t, "snax", config.Get("chip.shop").(string))
	// Assert int
	assert.Equal(t, 123, int(config.Get("foo.bar").(float64)))
	// Assert nil value
	assert.Nil(t, config.Get("somewhat.secret"))
	// Assert missing key
	assert.Empty(t, config.Get("fish.shop"))
	// Assert missing key with fallback.
	assert.Equal(t, "yums", config.GetWithFallback("fish.shop", "yums").(string))

}
