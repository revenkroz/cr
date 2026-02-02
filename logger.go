package cr

import "log"

type Logger interface {
	Logf(format string, args ...interface{})
}

// NewStdLogger creates a new logger that writes to the standard logger.
func NewStdLogger() Logger {
	return stdLogger{}
}

type stdLogger struct{}

func (n stdLogger) Logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// NewVoidLogger creates a new logger that does nothing.
func NewVoidLogger() Logger {
	return voidLogger{}
}

type voidLogger struct{}

func (n voidLogger) Logf(_ string, _ ...interface{}) {
}
