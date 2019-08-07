# skprconfig

This is a go package providing an interface to read config values on the skpr.io platform.

## Usage

```go
import "github.com/skpr/go-skprconfig"

// Get the configured value for "port", with a default fallback if missing.
listenPort := skprconfig.GetWithFallback("port", "8888")

// Get the configured value for "token", and return an error if missing.
token, err := skprconfig.Get("token")
if err != nil {
  panic("auth token not configured")
}
```
