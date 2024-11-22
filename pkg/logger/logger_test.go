package logger

import (
	"log/slog"
	"testing"

	"github.com/zhongxic/sellbot/config"
)

func TestLogging(t *testing.T) {
	logging := config.Logging{
		Level:   "debug",
		File:    "log.log",
		MaxAge:  7,
		MaxSize: 1024,
	}
	if err := Init(logging); err != nil {
		t.Fatal(err)
	}
	if logWriter == nil {
		t.Fatal("the writer has not been initialized")
	}

	defer Close()

	slog.Debug("this is a debug message")
	slog.Info("this is an info message")
	slog.Warn("this is a warn message")
	slog.Error("this is an error message")
}
