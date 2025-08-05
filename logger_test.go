package pdalog

import (
	"bytes"
	"encoding/json"
	"errors"
	"sync"
	"testing"
)

func TestLoggerLevels(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := Options{
		Writer: buf,
		Level:  InfoLevel,
	}
	log := New(opts)

	// Debug should not be logged at InfoLevel
	log.Debug().Msg("debug message")
	if buf.Len() > 0 {
		t.Error("Debug message was logged when level is Info")
	}

	// Info should be logged
	log.Info().Msg("info message")
	if buf.Len() == 0 {
		t.Error("Info message was not logged when level is Info")
	}

	// Reset buffer
	buf.Reset()

	// Change level to Debug
	log.SetLevel(DebugLevel)

	// Now Debug should be logged
	log.Debug().Msg("debug message")
	if buf.Len() == 0 {
		t.Error("Debug message was not logged when level is Debug")
	}
}

func TestStructuredLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := Options{
		Writer: buf,
		Level:  DebugLevel,
	}
	log := New(opts)

	// Log with structured data
	log.Info().
		Str("string", "value").
		Int("int", 123).
		Bool("bool", true).
		Err(errors.New("test error")).
		Msg("test message")

	// Parse the JSON output
	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check fields
	if entry["level"] != "info" {
		t.Errorf("Expected level to be 'info', got %v", entry["level"])
	}
	if entry["message"] != "test message" {
		t.Errorf("Expected message to be 'test message', got %v", entry["message"])
	}
	if entry["string"] != "value" {
		t.Errorf("Expected string field to be 'value', got %v", entry["string"])
	}
	if entry["int"] != float64(123) { // JSON numbers are float64
		t.Errorf("Expected int field to be 123, got %v", entry["int"])
	}
	if entry["bool"] != true {
		t.Errorf("Expected bool field to be true, got %v", entry["bool"])
	}
	if entry["error"] != "test error" {
		t.Errorf("Expected error field to be 'test error', got %v", entry["error"])
	}
}

func TestContextFields(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := Options{
		Writer: buf,
		Level:  InfoLevel,
	}
	log := New(opts)

	// Create a logger with context
	contextLogger := log.With("requestID", "123456")

	// Log with the context logger
	contextLogger.Info().Msg("context test")

	// Parse the JSON output
	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check context field
	if entry["requestID"] != "123456" {
		t.Errorf("Expected requestID to be '123456', got %v", entry["requestID"])
	}
}

func TestParseLevelString(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", DebugLevel},
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"error", ErrorLevel},
		{"fatal", FatalLevel},
		{"invalid", InfoLevel}, // Default is InfoLevel
	}

	for _, test := range tests {
		level := ParseLevel(test.input)
		if level != test.expected {
			t.Errorf("ParseLevel(%q) = %v, want %v", test.input, level, test.expected)
		}
	}
}

// MockHook is a test implementation of the Hook interface
type MockHook struct {
	mu           sync.Mutex
	Fired        bool
	FiredEntries []map[string]interface{}
	FiredLevels  []Level
	levels       []Level
}

func NewMockHook(levels ...Level) *MockHook {
	if len(levels) == 0 {
		levels = []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel}
	}
	return &MockHook{
		levels:       levels,
		FiredEntries: make([]map[string]interface{}, 0),
		FiredLevels:  make([]Level, 0),
	}
}

func (h *MockHook) Fire(entry map[string]interface{}) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Fired = true
	h.FiredEntries = append(h.FiredEntries, entry)

	// Extract and store the level
	if levelStr, ok := entry["level"].(string); ok {
		h.FiredLevels = append(h.FiredLevels, ParseLevel(levelStr))
	}

	return nil
}

func (h *MockHook) Levels() []Level {
	return h.levels
}

func TestAddRemoveHook(t *testing.T) {
	log := NewConsoleLogger()
	hook := NewMockHook()

	// Add hook
	log = log.AddHook(hook)

	// Check if hook was added
	if len(log.hooks) != 1 {
		t.Errorf("Expected 1 hook, got %d", len(log.hooks))
	}

	// Remove hook
	log = log.RemoveHook(hook)

	// Check if hook was removed
	if len(log.hooks) != 0 {
		t.Errorf("Expected 0 hooks, got %d", len(log.hooks))
	}
}

func TestHookFiring(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := Options{
		Writer: buf,
		Level:  DebugLevel,
	}
	log := New(opts)

	// Create a hook that only fires for Info and Error levels
	hook := NewMockHook(InfoLevel, ErrorLevel)
	log.AddHook(hook)

	// Debug should not trigger the hook
	log.Debug().Msg("debug message")
	if hook.Fired {
		t.Error("Hook fired for Debug level when it should only fire for Info and Error")
	}

	// Info should trigger the hook
	log.Info().Msg("info message")
	if !hook.Fired {
		t.Error("Hook did not fire for Info level")
	}

	// Reset the hook
	hook.Fired = false

	// Error should trigger the hook
	log.Error().Msg("error message")
	if !hook.Fired {
		t.Error("Hook did not fire for Error level")
	}

	// Check the number of fired entries
	if len(hook.FiredEntries) != 2 {
		t.Errorf("Expected 2 fired entries, got %d", len(hook.FiredEntries))
	}

	// Check the fired levels
	if len(hook.FiredLevels) != 2 {
		t.Errorf("Expected 2 fired levels, got %d", len(hook.FiredLevels))
	}

	if hook.FiredLevels[0] != InfoLevel {
		t.Errorf("Expected first fired level to be InfoLevel, got %v", hook.FiredLevels[0])
	}

	if hook.FiredLevels[1] != ErrorLevel {
		t.Errorf("Expected second fired level to be ErrorLevel, got %v", hook.FiredLevels[1])
	}
}

func TestNatsHook(t *testing.T) {
	// Create a mock NATS connection
	mockConn := &MockNatsConn{
		PublishedMessages: make(map[string][]byte),
	}

	// Create a NatsHook with a template subject
	hook := NewNatsHook(mockConn, "logs.{level}.{component}", InfoLevel, ErrorLevel)

	// Create a logger with the hook
	log := NewConsoleLogger()
	log.AddHook(hook)

	// Log a message with component field
	log.Info().Str("component", "api").Msg("API request received")

	// Check if the message was published to the correct subject
	if _, ok := mockConn.PublishedMessages["logs.info.api"]; !ok {
		t.Errorf("Message not published to expected subject logs.info.api")
	}

	// Log another message with a different component
	log.Error().Str("component", "database").Msg("Database connection failed")

	// Check if the message was published to the correct subject
	if _, ok := mockConn.PublishedMessages["logs.error.database"]; !ok {
		t.Errorf("Message not published to expected subject logs.error.database")
	}

	// Check the total number of published messages
	if len(mockConn.PublishedMessages) != 2 {
		t.Errorf("Expected 2 published messages, got %d", len(mockConn.PublishedMessages))
	}
}

// MockNatsConn is a mock implementation of the NATS connection
type MockNatsConn struct {
	PublishedMessages map[string][]byte
}

func (m *MockNatsConn) Publish(subject string, data []byte) error {
	m.PublishedMessages[subject] = data
	return nil
}
