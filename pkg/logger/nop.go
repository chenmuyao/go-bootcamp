package logger

type NopLogger struct{}

func NewNopLogger() *NopLogger {
	return &NopLogger{}
}

func (s *NopLogger) Debug(msg string, args ...Field) {
}

func (s *NopLogger) Info(msg string, args ...Field) {
}

func (s *NopLogger) Warn(msg string, args ...Field) {
}

func (s *NopLogger) Error(msg string, args ...Field) {
}
