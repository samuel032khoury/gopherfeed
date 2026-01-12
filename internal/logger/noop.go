package logger

// NoopLogger is a logger that does nothing.
// Use this when you don't care about logging.
type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (l *NoopLogger) Info(args ...any)                        {}
func (l *NoopLogger) Infow(msg string, keysAndValues ...any)  {}
func (l *NoopLogger) Warnw(msg string, keysAndValues ...any)  {}
func (l *NoopLogger) Errorw(msg string, keysAndValues ...any) {}
func (l *NoopLogger) Fatal(args ...any)                       { panic(args) }
