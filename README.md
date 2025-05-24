# GoEnv

A lightweight, intelligent Go library for loading environment variables from `.env` files with advanced features and robust error handling.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/raiesbo/goenv)](https://goreportcard.com/report/github.com/raiesbo/goenv)

## Features

- **Smart Discovery**: Recursively searches directories for `.env` files
- **Flexible Configuration**: Support for multiple file patterns and loading strategies
- **Type Safety**: Built-in converters for strings, integers, floats, and booleans
- **Quoted Values**: Handles single and double-quoted strings with spaces
- **Cycle Detection**: Prevents infinite loops with symbolic links
- **Depth Limiting**: Configurable search depth for performance
- **Error Context**: Detailed error messages with file and line information
- **Zero Dependencies**: Pure Go implementation with no external dependencies

## Quick Start

### Installation

```bash
go get github.com/raiesbo/goenv
```

### Basic Usage

Create a `.env` file in your project:

```env
# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_SSL=true

# API settings
API_KEY="your-secret-api-key"
API_TIMEOUT=30
DEBUG_MODE=false

# Complex values with quotes
MESSAGE='Hello, World!'
DESCRIPTION="This is a quoted string with spaces"
```

Load and use environment variables in your Go application:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/raiesbo/goenv"
)

func main() {
    // Load .env file automatically
    if err := goenv.Load(); err != nil {
        log.Fatal("Error loading .env file:", err)
    }

    // Get values with type conversion and fallbacks
    dbHost := goenv.GetString("DB_HOST", "localhost")
    dbPort := goenv.GetInt("DB_PORT", 3306)
    sslEnabled := goenv.GetBool("DB_SSL", false)
    timeout := goenv.GetFloat("API_TIMEOUT", 10.0)

    fmt.Printf("Connecting to %s:%d (SSL: %v)\n", dbHost, dbPort, sslEnabled)
    fmt.Printf("API timeout: %.1fs\n", timeout)

    // Get required variables (panics if not found)
    apiKey := goenv.MustGetString("API_KEY")
    fmt.Printf("API Key loaded: %s...\n", apiKey[:8])
}
```

## Advanced Configuration

### Custom Configuration

```go
package main

import (
    "log"
    "github.com/raiesbo/goenv"
)

func main() {
    config := &goenv.Config{
        EnvFiles:    []string{".env", ".env.local", ".env.production"},
        MaxDepth:    5,
        StopOnFirst: false, // Load all matching files
    }

    if err := goenv.LoadWithConfig(config); err != nil {
        log.Fatal("Failed to load environment:", err)
    }
}
```

### Loading Specific Files

```go
// Load a specific .env file
if err := goenv.LoadFile(".env.production"); err != nil {
    log.Fatal("Failed to load production config:", err)
}
```

### Environment-Specific Loading

```go
package main

import (
    "os"
    "github.com/raiesbo/goenv"
)

func main() {
    env := os.Getenv("GO_ENV")
    if env == "" {
        env = "development"
    }

    config := &goenv.Config{
        EnvFiles:    []string{".env", ".env." + env},
        StopOnFirst: false, // Load both base and environment-specific
    }

    goenv.LoadWithConfig(config)
}
```

## Supported File Formats

GoEnv supports various `.env` file formats:

```env
# Comments are supported
# Empty lines are ignored

# Basic key-value pairs
DATABASE_URL=postgres://localhost/mydb
PORT=8080

# Quoted strings (single or double quotes)
SECRET_KEY="my-secret-key-with-spaces"
MESSAGE='Hello, "World"!'

# Boolean values
DEBUG=true
PRODUCTION=false

# Numeric values
MAX_CONNECTIONS=100
TIMEOUT=30.5

# Empty values
OPTIONAL_CONFIG=
```

## API Reference

### Loading Functions

| Function | Description |
|----------|-------------|
| `Load()` | Load `.env` file using default configuration |
| `LoadWithConfig(config *Config)` | Load with custom configuration |
| `LoadFile(path string)` | Load a specific file |

### Type Conversion Functions

| Function | Description | Example |
|----------|-------------|---------|
| `GetString(key, fallback string)` | Get string value | `GetString("API_URL", "http://localhost")` |
| `GetInt(key string, fallback int)` | Get integer value | `GetInt("PORT", 8080)` |
| `GetBool(key string, fallback bool)` | Get boolean value | `GetBool("DEBUG", false)` |
| `GetFloat(key string, fallback float64)` | Get float value | `GetFloat("RATE", 1.5)` |
| `MustGetString(key string)` | Get required string (panics if missing) | `MustGetString("DATABASE_URL")` |

### Configuration Options

```go
type Config struct {
    EnvFiles    []string // File patterns to search for
    MaxDepth    int      // Maximum directory depth to search
    StopOnFirst bool     // Stop after finding first file
}
```

## Use Cases

### Web Applications

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/raiesbo/goenv"
)

func main() {
    goenv.Load()
    
    port := goenv.GetString("PORT", "8080")
    dbURL := goenv.MustGetString("DATABASE_URL")
    
    fmt.Printf("Starting server on port %s\n", port)
    fmt.Printf("Database: %s\n", dbURL)
    
    http.ListenAndServe(":" + port, nil)
}
```

### Microservices Configuration

```go
type Config struct {
    ServiceName string
    Port        int
    DatabaseURL string
    RedisURL    string
    LogLevel    string
    Debug       bool
}

func LoadConfig() *Config {
    goenv.Load()
    
    return &Config{
        ServiceName: goenv.GetString("SERVICE_NAME", "unknown"),
        Port:        goenv.GetInt("PORT", 8080),
        DatabaseURL: goenv.MustGetString("DATABASE_URL"),
        RedisURL:    goenv.GetString("REDIS_URL", "redis://localhost:6379"),
        LogLevel:    goenv.GetString("LOG_LEVEL", "info"),
        Debug:       goenv.GetBool("DEBUG", false),
    }
}
```

### Docker Integration

```dockerfile
# Dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o myapp

# Use .env file in development
CMD ["./myapp"]
```

```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    build: .
    env_file:
      - .env
      - .env.local
    ports:
      - "${PORT:-8080}:8080"
```

## Testing

GoEnv includes comprehensive tests covering all functionality:

```bash
# Run tests
go test -v

# Run tests with coverage
go test -cover

# Run benchmarks
go test -bench=.
```

## Error Handling

GoEnv provides detailed error messages for debugging:

```go
if err := goenv.Load(); err != nil {
    // Errors include file path and line numbers
    fmt.Printf("Failed to load .env: %v\n", err)
    // Example: "invalid format in /path/.env at line 5: INVALID_LINE"
}
```

## Comparison with Other Libraries

| Feature | GoEnv | godotenv | viper |
|---------|-------|----------|-------|
| File Discovery | ✅ Recursive | ❌ Manual path | ✅ Multiple sources |
| Type Conversion | ✅ Built-in | ❌ Manual | ✅ Comprehensive |
| Quoted Strings | ✅ Yes | ✅ Yes | ✅ Yes |
| Error Context | ✅ Detailed | ❌ Basic | ✅ Good |
| Zero Dependencies | ✅ Yes | ✅ Yes | ❌ Many deps |
| Configuration | ✅ Flexible | ❌ Limited | ✅ Extensive |
| Performance | ✅ Fast | ✅ Fast | ⚠️ Moderate |

## Best Practices

1. **Environment-Specific Files**: Use `.env.development`, `.env.production`, etc.
2. **Never Commit Secrets**: Add `.env*` to your `.gitignore`
3. **Provide Fallbacks**: Always use fallback values for non-critical settings
4. **Validate Required Variables**: Use `MustGetString()` for essential configuration
5. **Document Variables**: Create a `.env.example` file with all required variables

### Example `.env.example`

```env
# Database Configuration
DATABASE_URL=postgres://username:password@localhost:5432/dbname

# API Keys
API_KEY=your-api-key-here
SECRET_KEY=your-secret-key-here

# Application Settings
PORT=8080
DEBUG=false
LOG_LEVEL=info

# External Services
REDIS_URL=redis://localhost:6379
SMTP_HOST=smtp.example.com
SMTP_PORT=587
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the Node.js [dotenv](https://github.com/motdotla/dotenv) library
- Thanks to the Go community for feedback and contributions

---

**GoEnv** - Making environment configuration simple and reliable for Go applications.