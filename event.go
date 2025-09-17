package pdalog

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Event represents a log event
type Event struct {
	logger *Logger
	level  Level
	fields map[string]interface{}
	time   time.Time
}

// Str adds a string field to the event
func (e *Event) Str(key, val string) *Event {
	if e == nil {
		return nil
	}
	e.fields[key] = val
	return e
}

// Int adds an integer field to the event
func (e *Event) Int(key string, val int) *Event {
	if e == nil {
		return nil
	}
	e.fields[key] = val
	return e
}

// Bool adds a boolean field to the event
func (e *Event) Bool(key string, val bool) *Event {
	if e == nil {
		return nil
	}
	e.fields[key] = val
	return e
}

// Err adds an error field to the event
func (e *Event) Err(err error) *Event {
	if e == nil {
		return nil
	}
	if err == nil {
		return e
	}
	e.fields["error"] = err.Error()
	return e
}

// Any adds a field with any value to the event
func (e *Event) Any(key string, val interface{}) *Event {
	if e == nil {
		return nil
	}
	e.fields[key] = val
	return e
}

// Duration adds a duration field to the event
func (e *Event) Duration(key string, val time.Duration) *Event {
	if e == nil {
		return nil
	}
	e.fields[key] = val
	return e
}

// Time adds a time.Time field to the event
func (e *Event) Time(key string, val time.Time) *Event {
	if e == nil {
		return nil
	}
	e.fields[key] = val
	return e
}

// Hex adds a hex-encoded byte slice field to the event
func (e *Event) Hex(key string, val []byte) *Event {
	if e == nil {
		return nil
	}
	e.fields[key] = fmt.Sprintf("%x", val)
	return e
}

// Msg sends the event with the given message
func (e *Event) Msg(msg string) {
	if e == nil {
		return
	}

	// Create the log entry
	entry := map[string]interface{}{
		"level":   e.level.String(),
		"time":    e.time.Format(e.logger.timeFormat),
		"message": msg,
	}

	// Add all fields
	for k, v := range e.fields {
		entry[k] = v
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error marshaling log entry: %v\n", err)
		return
	}

	// Write to output
	e.logger.mu.Lock()
	defer e.logger.mu.Unlock()

	jsonData = append(jsonData, '\n')
	_, err = e.logger.writer.Write(jsonData)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error writing log entry: %v\n", err)
	}

	// Fire hooks
	for _, hook := range e.logger.hooks {
		// Check if this hook should be triggered for this level
		shouldFire := false
		for _, level := range hook.Levels() {
			if level == e.level {
				shouldFire = true
				break
			}
		}

		if shouldFire {
			if err := hook.Fire(entry); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error firing hook: %v\n", err)
			}
		}
	}

	// If fatal, exit the program
	if e.level == FatalLevel {
		os.Exit(1)
	}
}
