package logger

// Logger defines the logging interface used by internal packages.
// This allows internal packages to remain decoupled from specific logging implementations.
type Logger interface {
	// Info logs an informational message
	Info(args ...any)

	// Infow logs an informational message with structured key-value pairs
	Infow(msg string, keysAndValues ...any)

	// Warn logs a warning message with structured key-value pairs
	Warnw(msg string, keysAndValues ...any)

	// Error logs an error message with structured key-value pairs
	Errorw(msg string, keysAndValues ...any)
	// Fatal logs a fatal message and exits the program
	Fatal(args ...any)
}
