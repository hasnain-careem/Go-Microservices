package logger

import (
	"io"
	"log/slog"
	"os"
	"sync/atomic"
)

var (
	DBErrorCount     atomic.Int64
	NetworkErrorCount atomic.Int64
	ValidationErrorCount atomic.Int64
)

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(serviceName string) *Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(jsonHandler).With(
		"service", serviceName,
	)

	return &Logger{
		logger: logger,
	}
}

func NewLoggerWithWriter(serviceName string, w io.Writer) *Logger {
	jsonHandler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(jsonHandler).With(
		"service", serviceName,
	)

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}


func (l *Logger) WithValues(keyValues ...any) *Logger {
	return &Logger{
		logger: l.logger.With(keyValues...),
	}
}

func (l *Logger) LogRequest(method string, req any) {
	l.Info("received request", "method", method, "payload", req)
}

func (l *Logger) LogResponse(method string, res any) {
	l.Info("sending response", "method", method, "payload", res)
}

func IncrementDBErrorCount() {
	DBErrorCount.Add(1)
}

func IncrementNetworkErrorCount() {
	NetworkErrorCount.Add(1)
}

func IncrementValidationErrorCount() {
	ValidationErrorCount.Add(1)
}

func GetErrorCounts() map[string]int64 {
	return map[string]int64{
		"db_errors":        DBErrorCount.Load(),
		"network_errors":   NetworkErrorCount.Load(),
		"validation_errors": ValidationErrorCount.Load(),
	}
}