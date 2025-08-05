package pdalog

// Level represents logging level
type Level int8

const (
	// DebugLevel defines debug log level
	DebugLevel Level = iota
	// InfoLevel defines info log level
	InfoLevel
	// WarnLevel defines warn log level
	WarnLevel
	// ErrorLevel defines error log level
	ErrorLevel
	// FatalLevel defines fatal log level
	FatalLevel
)

var levelNames = map[Level]string{
	DebugLevel: "debug",
	InfoLevel:  "info",
	WarnLevel:  "warn",
	ErrorLevel: "error",
	FatalLevel: "fatal",
}

// String returns the string representation of the log level
func (l Level) String() string {
	if name, ok := levelNames[l]; ok {
		return name
	}
	return "unknown"
}

// ParseLevel parses a level string into a Level value
func ParseLevel(levelStr string) Level {
	switch levelStr {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel // Default is info
	}
}
