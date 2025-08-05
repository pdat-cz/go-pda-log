package pdalog

// Hook represents a log hook that processes log entries
type Hook interface {
	// Fire is called when a log event occurs
	Fire(entry map[string]interface{}) error
	// Levels returns the log levels this hook should be triggered for
	Levels() []Level
}
