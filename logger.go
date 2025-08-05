package pdalog

import (
	"io"
	"os"
	"sync"
	"time"
)

// Logger represents the core logger structure
type Logger struct {
	writer        io.Writer
	level         Level
	mu            sync.Mutex
	timeFormat    string
	contextFields map[string]interface{}
	hooks         []Hook
}

// Options for configuring a new logger
type Options struct {
	Writer     io.Writer
	Level      Level
	TimeFormat string
}

// DefaultOptions returns the default logger options
func DefaultOptions() Options {
	return Options{
		Writer:     os.Stdout,
		Level:      InfoLevel,
		TimeFormat: time.RFC3339,
	}
}

// New creates a new logger with the given options
func New(opts Options) *Logger {
	if opts.Writer == nil {
		opts.Writer = os.Stdout
	}
	if opts.TimeFormat == "" {
		opts.TimeFormat = time.RFC3339
	}

	return &Logger{
		writer:        opts.Writer,
		level:         opts.Level,
		timeFormat:    opts.TimeFormat,
		contextFields: make(map[string]interface{}),
	}
}

// NewConsoleLogger creates a new logger with console output
func NewConsoleLogger() *Logger {
	opts := DefaultOptions()
	return New(opts)
}

// SetLevel sets the logger's minimum level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel returns the current logger level
func (l *Logger) GetLevel() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// With returns a new logger with the given field added to its context
func (l *Logger) With(key string, value interface{}) *Logger {
	newLogger := &Logger{
		writer:        l.writer,
		level:         l.level,
		timeFormat:    l.timeFormat,
		contextFields: make(map[string]interface{}),
	}

	// Copy existing context fields
	for k, v := range l.contextFields {
		newLogger.contextFields[k] = v
	}

	// Add new field
	newLogger.contextFields[key] = value

	return newLogger
}

// Debug returns a debug level event logger
func (l *Logger) Debug() *Event {
	return l.newEvent(DebugLevel)
}

// Info returns an info level event logger
func (l *Logger) Info() *Event {
	return l.newEvent(InfoLevel)
}

// Warn returns a warn level event logger
func (l *Logger) Warn() *Event {
	return l.newEvent(WarnLevel)
}

// Error returns an error level event logger
func (l *Logger) Error() *Event {
	return l.newEvent(ErrorLevel)
}

// Fatal returns a fatal level event logger
func (l *Logger) Fatal() *Event {
	return l.newEvent(FatalLevel)
}

// newEvent creates a new Event with the given level
func (l *Logger) newEvent(level Level) *Event {
	if level < l.level {
		return nil
	}

	e := &Event{
		logger: l,
		level:  level,
		fields: make(map[string]interface{}),
		time:   time.Now(),
	}

	// Add context fields
	for k, v := range l.contextFields {
		e.fields[k] = v
	}

	return e
}

// AddHook adds a hook to the logger
func (l *Logger) AddHook(hook Hook) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hooks = append(l.hooks, hook)
	return l
}

// RemoveHook removes a hook from the logger
func (l *Logger) RemoveHook(hook Hook) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i, h := range l.hooks {
		if h == hook {
			l.hooks = append(l.hooks[:i], l.hooks[i+1:]...)
			break
		}
	}
	return l
}
