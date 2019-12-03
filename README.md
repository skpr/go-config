# skprconfig

[![CircleCI](https://circleci.com/gh/skpr/go-config.svg?style=svg)](https://circleci.com/gh/skpr/go-config)

This is a go package providing an interface to read config values on the skpr.io platform.

## Usage

```go
import "github.com/skpr/skprconfig"

config, err := skprconfig.NewConfig().Load()
if err != nil {
  panic("failed to load config")
}
bar := config.Get("foo")

// Reload config.
config, err := config.Load()

// Get the configured value for "port", with a default fallback if missing.
listenPort := config.GetWithFallback("port", "8888")

// Get the configured value for "token", and return an error if missing.
token, err := config.Get("token")
if err != nil {
  panic("auth token not configured")
}
```
