package skprconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/afero"
)

const (
	// DefaultPath is the default directory used to mount config and secrets into containers.
	DefaultPath string = "/etc/skpr"
	// DefaultTrimSuffix is the default pattern to trim from values.
	DefaultTrimSuffix string = "\n"
	// PathPartConfig is the path component where config values are mounted.
	PathPartConfig string = "config"
	// PathPartSecret is the path component where secret values are mounted.
	PathPartSecret string = "secret"
	// PathPartDefault is the path component where system-defined values are mounted.
	PathPartDefault string = "default"
	// PathPartOverride is the path component where user-defined values are mounted.
	PathPartOverride string = "override"
)

// Get returns the configured value of a given key, and the fallback value if no
// key does not exist.
//
// This is a convenience if the default skpr config mount point is being used.
func Get(key, fallback string) string {
	c := NewConfig(DefaultPath)
	return c.GetWithFallback(key, fallback)
}

// Config holds parameters for config.
type Config struct {
	Path       string
	TrimSuffix string
	FileSystem afero.Fs
}

// NewConfig returns a new Config struct with a given set of parameters.
func NewConfig(path string) *Config {
	fs := afero.NewOsFs()
	return &Config{
		Path:       path,
		TrimSuffix: DefaultTrimSuffix,
		FileSystem: fs,
	}
}

// GetWithFallback returns the configured value of a given key, and the fallback
// value if no key does not exist.
func (c *Config) GetWithFallback(key, fallback string) string {
	value, err := c.Get(key)
	if err != nil {
		return fallback
	}

	return value
}

// Get returns the configured value of a given key.
func (c *Config) Get(key string) (string, error) {
	var value string

	// These directories are weighted in order of precedence. Values specified
	// in multiple paths will be overridden by directories further the list.
	//
	// @see https://github.com/skpr/docs/blob/master/docs/config.md#deep-dive
	paths := []string{
		c.filePath(PathPartConfig, PathPartDefault, key),
		c.filePath(PathPartSecret, PathPartDefault, key),
		c.filePath(PathPartConfig, PathPartOverride, key),
		c.filePath(PathPartSecret, PathPartOverride, key),
	}

	configNoExist := true
	for _, file := range paths {
		if _, err := c.FileSystem.Stat(file); os.IsNotExist(err) {
			continue
		}
		contents, err := c.readFile(file)
		if err != nil {
			// This could indicate permissions issues with the mount, so exit.
			return "", err
		}
		configNoExist = false
		value = string(contents)
	}

	if configNoExist {
		err := fmt.Errorf("key not found")
		return "", err
	}

	return strings.TrimSuffix(value, DefaultTrimSuffix), nil
}

// dirPath returns the directory path of a specific config type and
// system/user provided value.
func (c *Config) dirPath(configType, defaultOrOverride string) string {
	return fmt.Sprintf("%s/%s/%s", c.Path, configType, defaultOrOverride)
}

// filePath returns the file path of a specific config type, system/user
// provided value, and key.
func (c *Config) filePath(configType, defaultOrOverride, key string) string {
	return fmt.Sprintf("%s/%s", c.dirPath(configType, defaultOrOverride), key)
}

// readFile adapted from io/ioutil.ReadFile() to use injected filesystem.
func (c *Config) readFile(filename string) ([]byte, error) {
	f, err := c.FileSystem.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
