package config

import (
	"testing"
)

func TestParse(t *testing.T) {
	expectedServerPort := 8080
	expectedLogLevel := "debug"
	expectedLogFile := "log.log"
	expectedLogFileAge := 7
	expectedLogFileSize := 1024

	config, err := Parse("config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if config.Server.Port != expectedServerPort {
		t.Errorf("expected server port [%v] actual [%v]", expectedServerPort, config.Server.Port)
	}
	if config.Logging.Level != expectedLogLevel {
		t.Errorf("expected log level [%v] actual [%v]", expectedLogLevel, config.Logging.Level)
	}
	if config.Logging.File != expectedLogFile {
		t.Errorf("expected log file [%v] actual [%v]", expectedLogFile, config.Logging.File)
	}
	if config.Logging.MaxAge != expectedLogFileAge {
		t.Errorf("expected log file age [%v] actual [%v]", expectedLogFile, config.Logging.MaxAge)
	}
	if config.Logging.MaxSize != expectedLogFileSize {
		t.Errorf("expected log file size [%v] actual [%v]", expectedLogFile, config.Logging.MaxSize)
	}
}
