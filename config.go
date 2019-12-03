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

// NewConfig creates a Config struct.
func NewConfig(options ...func(config *Config)) (*Config, error) {
	config := &Config{
		path: DefaultPath,
	}
	for _, option := range options {
		option(config)
	}
	return config.load()
}

// load the json file to a Config map.
func (c *Config) load() (*Config, error) {

	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		return c, errors.Wrap(err, "config file does not exist")
	}

	data, err := ioutil.ReadFile(c.path)
	if err != nil {
		return c, errors.Wrap(err, "Failed to read config file")
	}

	var configData map[string]interface{}

	err = json.Unmarshal(data, &configData)
	if err != nil {
		return c, errors.Wrap(err, "Failed to unmarshal config")
	}
	c.data = configData

	return c, nil
}

// Get returns a value for the key.
func (c *Config) Get(key string) interface{} {
	return c.GetWithFallback(key, "")
}

// GetWithFallback returns the configured value of a given key, and the fallback
// value if no key does not exist.
func (c *Config) GetWithFallback(key, fallback string) interface{} {
	if _, ok := c.data[key]; !ok {
		return fallback
	}
	return c.data[key]
}
