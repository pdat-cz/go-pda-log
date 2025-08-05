package pdalog

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestDurationField(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := Options{
		Writer: buf,
		Level:  DebugLevel,
	}
	log := New(opts)

	// Test with a valid duration
	duration := 5 * time.Second
	log.Info().
		Duration("elapsed", duration).
		Msg("test duration")

	// Parse the JSON output
	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check duration field (JSON marshals time.Duration as a string)
	if entry["elapsed"] != float64(5000000000) { // 5 seconds in nanoseconds as float64
		t.Errorf("Expected elapsed to be 5000000000, got %v (type: %T)", entry["elapsed"], entry["elapsed"])
	}

	// Test with nil receiver
	var nilEvent *Event
	if nilEvent.Duration("test", duration) != nil {
		t.Error("Expected nil.Duration to return nil")
	}
}

func TestTimeField(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := Options{
		Writer: buf,
		Level:  DebugLevel,
	}
	log := New(opts)

	// Test with a valid time
	testTime := time.Date(2025, 8, 4, 21, 2, 0, 0, time.UTC)
	log.Info().
		Time("timestamp", testTime).
		Msg("test time")

	// Parse the JSON output
	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check time field (JSON marshals time.Time in RFC3339 format)
	expectedTimeStr := "2025-08-04T21:02:00Z"
	if entry["timestamp"] != expectedTimeStr {
		t.Errorf("Expected timestamp to be %s, got %v (type: %T)", expectedTimeStr, entry["timestamp"], entry["timestamp"])
	}

	// Test with nil receiver
	var nilEvent *Event
	if nilEvent.Time("test", testTime) != nil {
		t.Error("Expected nil.Time to return nil")
	}
}

func TestHexField(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := Options{
		Writer: buf,
		Level:  DebugLevel,
	}
	log := New(opts)

	// Test with a valid byte slice
	data := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	log.Info().
		Hex("data", data).
		Msg("test hex")

	// Parse the JSON output
	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check hex field
	expectedHex := "deadbeef"
	if entry["data"] != expectedHex {
		t.Errorf("Expected data to be %s, got %v (type: %T)", expectedHex, entry["data"], entry["data"])
	}

	// Test with nil receiver
	var nilEvent *Event
	if nilEvent.Hex("test", data) != nil {
		t.Error("Expected nil.Hex to return nil")
	}
}
