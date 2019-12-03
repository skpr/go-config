package skprconfig

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

const (
	// DefaultPath is the default file path used to mount config.
	DefaultPath string = "/etc/skpr/data/config.json"
)

// Config represents the config.
type Config struct {
	path string
	data map[string]interface{}
}

// Load loads Config from file.
func Load(options ...func(config *Config)) (*Config, error) {
	config := &Config{
		path: DefaultPath,
	}
	for _, option := range options {
		option(config)
	}

	if _, err := os.Stat(config.path); os.IsNotExist(err) {
		return config, errors.Wrap(err, "config file does not exist")
	}

	data, err := ioutil.ReadFile(config.path)
	if err != nil {
		return config, errors.Wrap(err, "Failed to read config file")
	}

	var configData map[string]interface{}

	err = json.Unmarshal(data, &configData)
	if err != nil {
		return config, errors.Wrap(err, "Failed to unmarshal config")
	}
	config.data = configData

	return config, nil
}

// Get returns a string value for the key.
func (c *Config) Get(key string) (string, bool) {
	value, ok := c.getValue(key)
	if value == nil {
		value = ""
	}
	return value.(string), ok
}

// GetBool returns a boolean value for the key.
func (c *Config) GetBool(key string) (bool, bool) {
	value, ok := c.getValue(key)
	if value == nil {
		value = false
	}
	return value.(bool), ok
}

// GetInt returns a int value for the key.
func (c *Config) GetInt(key string) (int, bool) {
	value, ok := c.getValue(key)
	if value == nil {
		value = 0
	}
	return int(value.(float64)), ok
}

// GetFloat returns a float value for the key.
func (c *Config) GetFloat(key string) (float64, bool) {
	value, ok := c.getValue(key)
	if value == nil {
		value = 0.0
	}
	return value.(float64), ok
}

func (c *Config) getValue(key string) (interface{}, bool) {
	if _, ok := c.data[key]; !ok {
		return nil, false
	}
	return c.data[key], true
}

// GetWithFallback returns the configured value of a given key, and the fallback
// value if no key does not exist.
func (c *Config) GetWithFallback(key, fallback string) string {
	if _, ok := c.getValue(key); !ok {
		return fallback
	}
	value, _ := c.getValue(key)
	return value.(string)
}
