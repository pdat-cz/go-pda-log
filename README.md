# go-pda-log

[![Go Reference](https://pkg.go.dev/badge/github.com/pdat-cz/go-pda-log.svg)](https://pkg.go.dev/github.com/pdat-cz/go-pda-log)
[![Go Report Card](https://goreportcard.com/badge/github.com/pdat-cz/go-pda-log)](https://goreportcard.com/report/github.com/pdat-cz/go-pda-log)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A simple, structured logging library for Go applications with a focus on ease of use and flexibility.

## Features

- **Leveled logging**: Debug, Info, Warn, Error, Fatal
- **Structured logging**: Add key-value pairs to your logs
- **JSON output format**: Machine-readable logs
- **Context fields**: Include context in all log messages
- **Hooks system**: Send logs to multiple destinations
- **NATS integration**: Built-in support for NATS messaging system
- **Simple and intuitive API**: Inspired by zerolog's fluent API

## Installation

```bash
go get github.com/pdat-cz/go-pda-log
```

## Basic Usage

```go
package main

import (
    "errors"
    "github.com/pdat-cz/go-pda-log"
)

func main() {
    // Create a new logger
    log := pdalog.NewConsoleLogger()
    
    // Set log level
    log.SetLevel(pdalog.DebugLevel)
    
    // Basic logging
    log.Info().Msg("Application started")
    
    // Structured logging
    log.Info().
        Str("service", "api").
        Int("port", 8080).
        Msg("Server listening")
    
    // With context
    contextLogger := log.With("requestID", "12345")
    contextLogger.Debug().Str("path", "/users").Msg("Request received")
    
    // Error logging
    err := errors.New("database connection failed")
    log.Error().
        Err(err).
        Str("component", "database").
        Msg("Failed to connect to database")
}
```

## Advanced Usage

### Custom Configuration

You can configure the logger with custom options:

```go
opts := pdalog.Options{
    Writer:     os.Stdout,
    Level:      pdalog.InfoLevel,
    TimeFormat: time.RFC3339,
}
log := pdalog.New(opts)
```

### Using Hooks

Hooks allow you to send log entries to multiple destinations.

#### Custom File Hook Example

```go
// Create a custom hook that writes to a file
type FileHook struct {
    file   *os.File
    levels []pdalog.Level
}

func NewFileHook(filename string, levels ...pdalog.Level) (*FileHook, error) {
    if len(levels) == 0 {
        levels = []pdalog.Level{pdalog.InfoLevel, pdalog.WarnLevel, pdalog.ErrorLevel, pdalog.FatalLevel}
    }

    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return nil, err
    }

    return &FileHook{
        file:   file,
        levels: levels,
    }, nil
}

func (h *FileHook) Fire(entry map[string]interface{}) error {
    line := fmt.Sprintf("[%s] %s: %s\n",
        entry["time"],
        entry["level"],
        entry["message"],
    )

    _, err := h.file.WriteString(line)
    return err
}

func (h *FileHook) Levels() []golog.Level {
    return h.levels
}

// Usage
log := golog.NewConsoleLogger()
fileHook, err := NewFileHook("logs.txt", golog.InfoLevel, golog.ErrorLevel)
if err != nil {
    log.Fatal().Err(err).Msg("Failed to create file hook")
}
log.AddHook(fileHook)
```

#### NATS Hook

The library includes a built-in hook for sending logs to a NATS server. Here's a comprehensive example of how to use golog with NATS:

##### Basic NATS Integration

Here's a basic example of integrating golog with NATS:

```go
package main

import (
    "github.com/nats-io/nats.go"
    "github.com/pdat-cz/go-pda-log"
    "os"
    "time"
)

func main() {
    // Create a logger
    log := pdalog.NewConsoleLogger()
    
    // Connect to NATS
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to connect to NATS")
    }
    defer nc.Close()
    
    // Create a NATS hook with template-based subject
    // The subject will be constructed as "logs.{level}.{component}"
    natsHook := pdalog.NewNatsHook(nc, "logs.{level}.{component}")
    log.AddHook(natsHook)
    
    // Log messages with different components
    log.Info().Str("component", "api").Msg("API server started")
    log.Error().Str("component", "database").Msg("Database connection failed")
    
    // Give time for messages to be sent
    time.Sleep(100 * time.Millisecond)
}
```

##### Advanced NATS Integration with Subscribers

This example shows how to set up both a publisher (logger) and subscribers to listen for log messages:

```go
package main

import (
    "fmt"
    "github.com/nats-io/nats.go"
    "github.com/pdat-cz/go-pda-log"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Create a logger
    log := pdalog.NewConsoleLogger()
    
    // Connect to NATS
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to connect to NATS")
    }
    defer nc.Close()
    
    // Set up subscribers for different log types
    
    // Subscribe to all logs
    nc.Subscribe("logs.>", func(msg *nats.Msg) {
        fmt.Printf("Received log: %s\n", string(msg.Data))
    })
    
    // Subscribe only to error logs
    nc.Subscribe("logs.error.>", func(msg *nats.Msg) {
        fmt.Printf("ERROR LOG: %s\n", string(msg.Data))
    })
    
    // Create a NATS hook with template-based subject
    natsHook := pdalog.NewNatsHook(nc, "logs.{level}.{component}")
    log.AddHook(natsHook)
    
    // Start logging in a separate goroutine
    go func() {
        for {
            log.Info().Str("component", "api").Msg("API request processed")
            log.Error().Str("component", "database").Msg("Database timeout")
            time.Sleep(2 * time.Second)
        }
    }()
    
    // Wait for interrupt signal
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh
    
    fmt.Println("Shutting down...")
}
```

##### Subject Patterns and Variable Substitution

The NATS hook supports variable substitution in the subject. Any field in your log entry can be used in the subject pattern:

```go
// Basic subject
natsHook := pdalog.NewNatsHook(nc, "logs")

// Level-based subject
natsHook := pdalog.NewNatsHook(nc, "logs.{level}")

// Service and level-based subject
natsHook := pdalog.NewNatsHook(nc, "logs.{service}.{level}")

// Using with structured logging
log.Info().
    Str("service", "auth").
    Str("user_id", "12345").
    Msg("User authenticated")  // Published to "logs.auth.info"
```

##### Filtering Log Levels

You can specify which log levels should trigger the NATS hook:

```go
// Only send error and fatal logs to NATS
natsHook := pdalog.NewNatsHook(nc, "logs.{level}", pdalog.ErrorLevel, pdalog.FatalLevel)
log.AddHook(natsHook)
```

## Log Levels

The following log levels are available, in order of increasing severity:

- `DebugLevel`: Debug information for developers
- `InfoLevel`: General information about application progress
- `WarnLevel`: Warning events that might cause issues
- `ErrorLevel`: Error events that might still allow the application to continue
- `FatalLevel`: Fatal events that cause the application to exit

## Contributing

Contributions are welcome! Here's how you can contribute:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature/my-feature`
5. Submit a pull request

Please make sure your code follows the Go coding standards and includes appropriate tests.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

This library is inspired by [zerolog](https://github.com/rs/zerolog) and aims to provide a simple, yet powerful logging solution for Go applications.