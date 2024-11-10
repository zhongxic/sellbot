package config

import (
	"testing"
)

func TestParse(t *testing.T) {
	expectedLogLevel := "debug"
	expectedLogFile := "log.log"
	expectedLogFileAge := 7
	expectedLogFileSize := 1024

	config, err := Parse("config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if config.Logging.Level != expectedLogLevel {
		t.Fatalf("expected log level [%v] actual [%v]", expectedLogLevel, config.Logging.Level)
	}
	if config.Logging.File != expectedLogFile {
		t.Fatalf("expected log file [%v] actual [%v]", expectedLogFile, config.Logging.File)
	}
	if config.Logging.MaxAge != expectedLogFileAge {
		t.Fatalf("expected log file age [%v] actual [%v]", expectedLogFile, config.Logging.MaxAge)
	}
	if config.Logging.MaxSize != expectedLogFileSize {
		t.Fatalf("expected log file size [%v] actual [%v]", expectedLogFile, config.Logging.MaxSize)
	}
}
