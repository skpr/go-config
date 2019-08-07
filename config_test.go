package skprconfig

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Get(t *testing.T) {
	config := &Config{
		FileSystem: createVirtualFS(),
		TrimSuffix: DefaultTrimSuffix,
		Path:       DefaultPath,
	}

	configName := "some.config"
	// Test empty string and error are returned if value is missing.
	value, err := config.Get(configName)
	assert.Error(t, err, "correctly identified key not found")
	assert.Equal(t, "", value, "empty string returned when config missing")

	// Test single, non-overridden value.
	err = afero.WriteFile(config.FileSystem, fmt.Sprintf("/etc/skpr/config/default/%s", configName), []byte("aaa"), 0644)
	assert.Nil(t, err)
	value, err = config.Get(configName)
	assert.Nil(t, err, "no error if key exists")
	assert.Equal(t, "aaa", value, "correct default value returned")

	// Test overriding config value.
	err = afero.WriteFile(config.FileSystem, fmt.Sprintf("/etc/skpr/config/override/%s", configName), []byte("bbb"), 0644)
	assert.Nil(t, err)
	value, err = config.Get(configName)
	assert.Nil(t, err, "no error if key exists")
	assert.Equal(t, "bbb", value, "correct overridden value returned")

	// Test keys in secret have precedence over config.
	secretName := "some.secret"
	err = afero.WriteFile(config.FileSystem, fmt.Sprintf("/etc/skpr/config/override/%s", secretName), []byte("ccc"), 0644)
	assert.Nil(t, err)
	err = afero.WriteFile(config.FileSystem, fmt.Sprintf("/etc/skpr/secret/override/%s", secretName), []byte("ddd"), 0644)
	assert.Nil(t, err)
	value, err = config.Get(secretName)
	assert.Nil(t, err, "no error if key exists")
	assert.Equal(t, "ddd", value, "secret correctly overrides config value")
}

// Helper function to create the skpr config directory structure.
func createVirtualFS() afero.Fs {
	var testFS = afero.NewMemMapFs()
	paths := []string{
		"/etc/skpr/config/default",
		"/etc/skpr/config/override",
		"/etc/skpr/secret/default",
		"/etc/skpr/secret/override",
	}
	for _, path := range paths {
		err := testFS.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}
	return testFS
}
