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
	Init(logging)
	defer Close()
	Debug("this is a debug message")
	Info("this is an info message")
	Warn("this is a warn message")
	Error("this is an error message")
}
