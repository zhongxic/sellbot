package logger

import (
	"testing"

	"github.com/zhongxic/sellbot/config"
)

var logging = config.Logging{
	Level:   "debug",
	File:    "log.log",
	MaxAge:  7,
	MaxSize: 1024,
}

func TestLogging(t *testing.T) {
	if err := Init(logging); err != nil {
		t.Fatal(err)
	}
	if logWriter == nil {
		t.Fatal("the writer has not been initialized")
	}
	if logger == nil {
		t.Fatal("the logger has not been initialized")
	}
	if stdLogger == nil {
		t.Fatal("the stdout logger has not been initialized")
	}

	defer Close()

	Debug("this is a debug message")
	Info("this is an info message")
	Warn("this is a warn message")
	Error("this is an error message")
}
