package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/zhongxic/sellbot/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	addition     = "addition"
	loggerClosed = "the file logger has been closed"
)

var (
	initOnce  sync.Once
	closeOnce sync.Once
	logWriter io.WriteCloser
	logger    *slog.Logger
	stdLogger *slog.Logger
	closed    bool
)

func init() {
	stdLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

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
	logger = slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logWriter), &slog.HandlerOptions{Level: logLevel}))
	return nil
}

func Debug(msg string, args ...any) {
	if isLoggerInvalid() {
		logToStdOut(msg, args...)
		return
	}
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	if isLoggerInvalid() {
		logToStdOut(msg, args...)
		return
	}
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	if isLoggerInvalid() {
		logToStdOut(msg, args...)
		return
	}
	logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	if isLoggerInvalid() {
		logToStdOut(msg, args...)
		return
	}
	logger.Error(msg, args...)
}

func isLoggerInvalid() bool {
	return logger == nil || closed
}

func logToStdOut(msg string, args ...any) {
	record := slog.NewRecord(time.Now(), slog.LevelInfo, msg, 0)
	record.Add(args...)
	record.Add(slog.String(addition, loggerClosed))
	stdLogger.Handler().Handle(context.Background(), record)
}

func Close() error {
	var err error
	closeOnce.Do(func() {
		closed = true
		if logWriter != nil {
			err = logWriter.Close()
		}
	})
	return err
}
