package logger

import (
	"log/slog"
	"os"
	"sync"
)

type Logger struct {
	Slg *slog.Logger
}

var (
	once sync.Once
	log  *Logger
)

func NewLogger() *Logger {
	handler := slog.NewTextHandler(os.Stdout, nil)
	once.Do(func() {
		log = &Logger{
			Slg: slog.New(handler),
		}
	})
	return log
}

func Log() *slog.Logger {
	if log == nil {
		log = NewLogger()
	}
	return log.Slg
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.Slg.Info(msg, keysAndValues...)
}

func (l *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.Slg.Error(msg, append(keysAndValues, slog.Any("error", err))...)
}
