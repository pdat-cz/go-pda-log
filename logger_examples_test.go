package pdalog_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pdat-cz/go-pda-log"
	"strings"
	"time"
)

// setupLogger creates a new logger with a buffer for testing
func setupLogger() (*pdalog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	opts := pdalog.Options{
		Writer: &buf,
		Level:  pdalog.InfoLevel,
	}
	log := pdalog.New(opts)
	return log, &buf
}

// parseLogEntry parses the JSON log entry from the buffer
func parseLogEntry(buf *bytes.Buffer) map[string]interface{} {
	var entry map[string]interface{}
	json.Unmarshal(buf.Bytes(), &entry)
	return entry
}

// Example demonstrates basic usage of the logger
func Example() {
	log, buf := setupLogger()

	// Log messages at different levels
	log.Info().Msg("This is an info message")
	log.Warn().Msg("This is a warning message")
	log.Error().Msg("This is an error message")

	// Verify that messages were logged
	output := buf.String()
	fmt.Println("Info message logged:", strings.Contains(output, "info"))
	fmt.Println("Warn message logged:", strings.Contains(output, "warn"))
	fmt.Println("Error message logged:", strings.Contains(output, "error"))

	// Output:
	// Info message logged: true
	// Warn message logged: true
	// Error message logged: true
}

// ExampleNew demonstrates creating a logger with custom options
func ExampleNew() {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create custom options
	opts := pdalog.Options{
		Writer: &buf,
		Level:  pdalog.InfoLevel,
	}

	// Create a new logger with custom options
	log := pdalog.New(opts)

	// Log messages
	log.Debug().Msg("This debug message won't be logged due to level setting")
	log.Info().Msg("This info message will be logged")

	// Verify the debug message was filtered out and info was logged
	output := buf.String()
	fmt.Println("Debug message logged:", strings.Contains(output, "debug"))
	fmt.Println("Info message logged:", strings.Contains(output, "info"))

	// Output:
	// Debug message logged: false
	// Info message logged: true
}

// ExampleNewConsoleLogger demonstrates creating a logger with default console output
func ExampleNewConsoleLogger() {
	// This example is simplified since we can't test console output directly
	log := pdalog.NewConsoleLogger()

	// Just verify the logger was created with the correct level
	fmt.Println("Default level:", log.GetLevel() == pdalog.InfoLevel)

	// Output:
	// Default level: true
}

// ExampleLogger_With demonstrates using contextual logging
func ExampleLogger_With() {
	log, buf := setupLogger()

	// Create a logger with context fields
	requestLogger := log.With("request_id", "req-123456").With("user_id", "user-789")

	// Log a message with the context fields
	requestLogger.Info().Msg("Processing request")

	// Parse the JSON to verify fields
	entry := parseLogEntry(buf)

	// Print the relevant fields
	fmt.Println("Message:", entry["message"])
	fmt.Println("Request ID:", entry["request_id"])
	fmt.Println("User ID:", entry["user_id"])

	// Output:
	// Message: Processing request
	// Request ID: req-123456
	// User ID: user-789
}

// ExampleEvent_Str demonstrates adding string fields to log events
func ExampleEvent_Str() {
	log, buf := setupLogger()

	// Add string fields to the log event
	log.Info().
		Str("service", "user-service").
		Str("method", "GetUser").
		Msg("User retrieved successfully")

	// Parse the JSON to verify fields
	entry := parseLogEntry(buf)

	// Print the relevant fields
	fmt.Println("Message:", entry["message"])
	fmt.Println("Service:", entry["service"])
	fmt.Println("Method:", entry["method"])

	// Output:
	// Message: User retrieved successfully
	// Service: user-service
	// Method: GetUser
}

// ExampleEvent_Int demonstrates adding integer fields to log events
func ExampleEvent_Int() {
	log, buf := setupLogger()

	// Add integer fields to the log event
	log.Info().
		Int("status_code", 200).
		Int("response_time_ms", 45).
		Msg("Request completed")

	// Parse the JSON to verify fields
	entry := parseLogEntry(buf)

	// Print the relevant fields
	fmt.Println("Message:", entry["message"])
	fmt.Println("Status code:", int(entry["status_code"].(float64)))
	fmt.Println("Response time (ms):", int(entry["response_time_ms"].(float64)))

	// Output:
	// Message: Request completed
	// Status code: 200
	// Response time (ms): 45
}

// ExampleEvent_Bool demonstrates adding boolean fields to log events
func ExampleEvent_Bool() {
	log, buf := setupLogger()

	// Add boolean fields to the log event
	log.Info().
		Bool("cache_hit", true).
		Bool("authenticated", false).
		Msg("Cache status")

	// Parse the JSON to verify fields
	entry := parseLogEntry(buf)

	// Print the relevant fields
	fmt.Println("Message:", entry["message"])
	fmt.Println("Cache hit:", entry["cache_hit"])
	fmt.Println("Authenticated:", entry["authenticated"])

	// Output:
	// Message: Cache status
	// Cache hit: true
	// Authenticated: false
}

// ExampleEvent_Err demonstrates adding error fields to log events
func ExampleEvent_Err() {
	log, buf := setupLogger()

	// Create an error
	err := errors.New("connection refused")

	// Add the error to the log event
	log.Error().
		Err(err).
		Str("server", "db-01").
		Msg("Database connection failed")

	// Parse the JSON to verify fields
	entry := parseLogEntry(buf)

	// Print the relevant fields
	fmt.Println("Message:", entry["message"])
	fmt.Println("Error:", entry["error"])
	fmt.Println("Server:", entry["server"])

	// Output:
	// Message: Database connection failed
	// Error: connection refused
	// Server: db-01
}

// ExampleEvent_Duration demonstrates adding duration fields to log events
func ExampleEvent_Duration() {
	log, buf := setupLogger()

	// Create a duration
	duration := 150 * time.Millisecond

	// Add the duration to the log event
	log.Info().
		Duration("response_time", duration).
		Msg("Request processed")

	// Parse the JSON to verify fields
	entry := parseLogEntry(buf)

	// Print the relevant fields
	fmt.Println("Message:", entry["message"])
	fmt.Println("Has response time:", entry["response_time"] != nil)

	// Output:
	// Message: Request processed
	// Has response time: true
}

// ExampleEvent_Time demonstrates adding time fields to log events
func ExampleEvent_Time() {
	log, buf := setupLogger()

	// Create a time
	timestamp := time.Date(2025, 8, 5, 9, 58, 0, 0, time.UTC)

	// Add the time to the log event
	log.Info().
		Time("created_at", timestamp).
		Msg("Resource created")

	// Parse the JSON to verify fields
	entry := parseLogEntry(buf)

	// Print the relevant fields
	fmt.Println("Message:", entry["message"])
	fmt.Println("Has created_at:", entry["created_at"] != nil)

	// Output:
	// Message: Resource created
	// Has created_at: true
}

// ExampleEvent_Hex demonstrates adding hex-encoded byte fields to log events
func ExampleEvent_Hex() {
	log, buf := setupLogger()

	// Create a byte slice
	data := []byte{0xDE, 0xAD, 0xBE, 0xEF}

	// Add the hex-encoded data to the log event
	log.Info().
		Hex("signature", data).
		Msg("Data signed")

	// Parse the JSON to verify fields
	entry := parseLogEntry(buf)

	// Print the relevant fields
	fmt.Println("Message:", entry["message"])
	fmt.Println("Signature:", entry["signature"])

	// Output:
	// Message: Data signed
	// Signature: deadbeef
}

// ExampleEvent_Any demonstrates adding arbitrary fields to log events
func ExampleEvent_Any() {
	log, buf := setupLogger()

	// Create a complex data structure
	data := map[string]interface{}{
		"id":   123,
		"name": "example",
	}

	// Add the data to the log event
	log.Info().
		Any("metadata", data).
		Msg("Data processed")

	// Parse the JSON to verify fields
	entry := parseLogEntry(buf)

	// Print the relevant fields
	fmt.Println("Message:", entry["message"])
	fmt.Println("Has metadata:", entry["metadata"] != nil)

	// Output:
	// Message: Data processed
	// Has metadata: true
}

// PrintHook is a simple hook that captures log entries
type PrintHook struct {
	levels    []pdalog.Level
	lastEntry map[string]interface{}
}

// Fire is called when a log event occurs
func (h *PrintHook) Fire(entry map[string]interface{}) error {
	h.lastEntry = entry
	return nil
}

// Levels returns the log levels this hook should be triggered for
func (h *PrintHook) Levels() []pdalog.Level {
	return h.levels
}

// ExampleLogger_AddHook demonstrates adding a hook to the logger
func ExampleLogger_AddHook() {
	log, _ := setupLogger()

	// Create and add a hook that only fires for error and fatal levels
	hook := &PrintHook{
		levels: []pdalog.Level{pdalog.ErrorLevel, pdalog.FatalLevel},
	}
	log.AddHook(hook)

	// This won't trigger the hook
	log.Info().Msg("This is an info message")

	// This will trigger the hook
	log.Error().Msg("This is an error message")

	// Verify the hook was triggered for error but not info
	fmt.Println("Hook triggered for error level:", hook.lastEntry != nil)
	if hook.lastEntry != nil {
		fmt.Println("Hook received message:", hook.lastEntry["message"])
	}

	// Output:
	// Hook triggered for error level: true
	// Hook received message: This is an error message
}

// ExampleParseLevel demonstrates parsing a level string into a Level value
func ExampleParseLevel() {
	// Parse level strings
	debugLevel := pdalog.ParseLevel("debug")
	infoLevel := pdalog.ParseLevel("info")
	warnLevel := pdalog.ParseLevel("warn")
	errorLevel := pdalog.ParseLevel("error")
	fatalLevel := pdalog.ParseLevel("fatal")
	unknownLevel := pdalog.ParseLevel("unknown")

	// Print the levels
	fmt.Println("Debug level:", debugLevel)
	fmt.Println("Info level:", infoLevel)
	fmt.Println("Warn level:", warnLevel)
	fmt.Println("Error level:", errorLevel)
	fmt.Println("Fatal level:", fatalLevel)
	fmt.Println("Unknown level (defaults to info):", unknownLevel)

	// Output:
	// Debug level: debug
	// Info level: info
	// Warn level: warn
	// Error level: error
	// Fatal level: fatal
	// Unknown level (defaults to info): info
}
