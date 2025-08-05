package pdalog

import (
	"encoding/json"
	// nats is imported for users who will pass a real nats.Conn to NewNatsHook
	_ "github.com/nats-io/nats.go"
	"strings"
)

// NatsConn is an interface that defines the methods needed from a NATS connection
type NatsConn interface {
	Publish(subject string, data []byte) error
}

// NatsHook sends log entries to NATS
type NatsHook struct {
	conn    NatsConn
	subject string
	levels  []Level
}

// NewNatsHook creates a new NATS hook.
func NewNatsHook(conn NatsConn, subject string, levels ...Level) *NatsHook {
	if len(levels) == 0 {
		levels = []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel}
	}

	return &NatsHook{
		conn:    conn,
		subject: subject,
		levels:  levels,
	}
}

// Fire sends the log entry to NATS
func (h *NatsHook) Fire(entry map[string]interface{}) error {

	subject := h.subject

	// Simple variable substitution
	for key, value := range entry {
		if strValue, ok := value.(string); ok {
			placeholder := "{" + key + "}"
			subject = strings.Replace(subject, placeholder, strValue, -1)
		}
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return h.conn.Publish(subject, data)
}

// Levels returns the log levels this hook should be triggered for
func (h *NatsHook) Levels() []Level {
	return h.levels
}
