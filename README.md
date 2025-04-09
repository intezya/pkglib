I'll help you improve your README to highlight the modular structure of the library and the examples package. Here's a revised version:

# pkglib

**pkglib** is a modular, and extensible Go library with reusable packages for building robust backend services. It provides a collection of utility packages designed to be composable and independent, allowing you to import only what you need.

## 📦 Architecture

The key feature of **pkglib** is its truly modular design. Each package is a standalone Go module with its own `go.mod` file and specific dependencies, meaning:

- You can import just the packages you need
- Each package pulls in only the dependencies it requires
- No dependency bloat in your projects

For example:

```
module github.com/intezya/pkglib/logger
require go.uber.org/zap v1.27.0
require go.uber.org/multierr v1.10.0 // indirect
```

```
module github.com/intezya/pkglib/sliceutils
go 1.20
```

```
module github.com/intezya/pkglib/jwt
require github.com/golang-jwt/jwt/v5 v5.2.2
```

## 📚 Available Packages

### `configloader`
Load environment variables with fallback values and type conversion.
```go
configloader.LoadEnv()               // Loads .env file (supports --env-file flag)
port := configloader.GetEnvInt("PORT", 8080)
debug := configloader.GetEnvBool("DEBUG", false)
```

### `sliceutils`
Generic helper functions for slice operations.
```go
ids := []int{1, 2, 2, 3}
unique := sliceutils.ToSet(ids)      // map[int]bool{1: true, 2: true, 3: true}
```

### `jwt`
Easy-to-use JWT generation and validation with support for custom claims.
```go
service := jwt.New(myConfig)
token, _ := service.Generate(jwt.AuthenticationData{"user_id": "123"})
data, _ := service.Validate(token)   // jwt.AuthenticationData{"user_id": "123"}
```

### `crypto`
Cryptographic utilities for secure applications.
```go
salt, _ := crypto.Salt(16)
hash := crypto.HashArgon2("password123", salt, nil)
valid := crypto.VerifyArgon2("password123", hash, nil)  // true
encoded := crypto.Base64Encode([]byte("hello world"))
```

### `generate`
Generate random strings and cryptographic hashes.
```go
password := generate.RandomString(12, generate.Charset.AllCharset)
sha := generate.RandomSHA512()
uuid := generate.RandomUUID()
```

### `itertools`
Generic functions for working with collections.
```go
squares := itertools.Map(func(n int) int { return n * n }, []int{1, 2, 3})    // []int{1, 4, 9}
evens := itertools.Filter(func(n int) bool { return n%2 == 0 }, []int{1, 2, 3, 4})  // []int{2, 4}
```

### `logger`
Structured logging built on top of zap with optional Loki integration.
```go
log, _ := logger.New(
    logger.WithDebug(true),
    logger.WithEnvironment(logger.EnvDevelopment),
    logger.WithTimeZone("UTC"),
)
log.Info("Server started", "port", 8080, "mode", "development")
```

## 🔍 Examples

The library includes a dedicated `examples` package with practical use cases for each module. These examples demonstrate how to use each package effectively:

```
examples/
├── configloader/
│   └── dotenv_example.go
├── crypto/
│   ├── encode_example.go
│   ├── generate_salt_example.go
│   ├── hash_example.go
│   └── password_example.go
├── generate/
│   ├── charset_example.go
│   └── generative_example.go
├── itertools/
│   └── itertools_example.go
├── jwt/
│   └── service_example.go
├── logger/
│   ├── logger_example.go
│   └── logger_example_scenarios.go
└── sliceutils/
    └── to_set_example.go
```

The examples package also has its own `go.mod` file to demonstrate real-world usage.

## 🚀 Getting Started

### Installation

You can import specific packages directly:

```bash
go get github.com/intezya/pkglib/logger
go get github.com/intezya/pkglib/jwt
```

Then import only what you need:

```go
import (
    "github.com/intezya/pkglib/configloader"
    "github.com/intezya/pkglib/logger"
    "github.com/intezya/pkglib/jwt"
)
```

### Example Usage

```go
package main

import (
    "github.com/intezya/pkglib/configloader"
    "github.com/intezya/pkglib/logger"
    "github.com/intezya/pkglib/generate"
)

func main() {
    // Load environment variables
    configloader.LoadEnv()
    
    // Initialize logger
    log, _ := logger.New(
        logger.WithDebug(configloader.GetEnvBool("DEBUG", false)),
        logger.WithEnvironment(configloader.GetEnv("ENV", "dev")),
    )
    
    // Generate a random API key
    apiKey := generate.RandomString(32, generate.Charset.AlphaNumericCharset)
    
    log.Info("Application initialized", 
        "api_key", apiKey[:5]+"...",  // Log only first 5 chars for security
        "debug", configloader.GetEnvBool("DEBUG", false),
    )
}
```

## 🔧 Advanced Features

### Logging to Loki

The `logger` package supports sending logs to Grafana Loki for centralized log storage:

```go
lokiConfig := logger.NewLokiConfig(
    "http://loki:3100/loki/api/v1/push",
    map[string]string{
        "app": "my-service",
        "env": "production",
    },
)

log, _ := logger.New(
    logger.WithEnvironment(logger.EnvProduction),
    logger.WithLoki(lokiConfig),
)

log.Info("This log will be sent to both stdout and Loki")
```

### JWT with Custom Claims

```go
type CustomClaims struct {
    UserID    string `json:"user_id"`
    UserRole  string `json:"role"`
}

config := jwt.Config{
    Secret: "my-secret-key",
    Issuer: "my-service",
    ExpiryMinutes: 60,
}

service := jwt.New(config)

claims := CustomClaims{
    UserID: "123",
    UserRole: "admin",
}

token, _ := service.Generate(claims)
```

## 🛠️ Requirements

- Go 1.20 or higher
- Each package pulls in only the dependencies it requires
