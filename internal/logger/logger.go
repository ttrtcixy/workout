package logger

import (
	"fmt"
	"github.com/ttrtcixy/workout/internal/logger/lib"
	"log/slog"
	"os"
)

type Logger interface {
	Info(format string, a ...any)
	Error(format string, a ...any)
	ErrorOp(op string, err error)
	Fatal(format string, a ...any)
	Debug(format string, a ...any)
	Warning(format string, a ...any)
}

type slogLogger struct {
	log *slog.Logger
}

func (s slogLogger) ErrorOp(op string, err error) {
	s.log.Error(fmt.Sprintf("%s: %s", op, err.Error()))
}

func Load() Logger {
	opts := colorLog.PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := colorLog.NewPrettyHandler(os.Stderr, opts)
	logger := slog.New(handler)

	return slogLogger{log: logger}
}

func (s slogLogger) Info(format string, a ...any) {
	s.log.Info(fmt.Sprintf(format, a...))
}

func (s slogLogger) Error(format string, a ...any) {
	s.log.Error(fmt.Sprintf(format, a...))
}

func (s slogLogger) Fatal(format string, a ...any) {
	s.log.Error(fmt.Sprintf(format, a...))
	os.Exit(1)
}

func (s slogLogger) Debug(format string, a ...any) {
	s.log.Debug(fmt.Sprintf(format, a...))
}

func (s slogLogger) Warning(format string, a ...any) {
	s.log.Warn(fmt.Sprintf(format, a...))
}
