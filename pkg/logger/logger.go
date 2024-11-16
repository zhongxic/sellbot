package logger

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/zhongxic/sellbot/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	initOnce  sync.Once
	closeOnce sync.Once
	logWriter io.WriteCloser
)

func Init(logging config.Logging) error {
	var err error
	initOnce.Do(func() {
		err = initLogger(logging)
	})
	return err
}

func initLogger(logging config.Logging) error {
	logWriter = &lumberjack.Logger{
		Filename:  logging.File,
		MaxAge:    logging.MaxAge,
		MaxSize:   logging.MaxSize,
		Compress:  true,
		LocalTime: true,
	}
	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(logging.Level)); err != nil {
		return err
	}
	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logWriter), &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(logLevel)
	return nil
}

func Close() error {
	var err error
	closeOnce.Do(func() {
		if logWriter != nil {
			err = logWriter.Close()
		}
	})
	return err
}
