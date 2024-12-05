package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParse(t *testing.T) {
	expected := &Config{
		Server: Server{
			Port: 8080,
		},
		Logging: Logging{
			Level:   "debug",
			File:    "log.log",
			MaxAge:  7,
			MaxSize: 1024,
		},
	}
	_, err := os.Stat("testdata")
	if errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir("testdata", 0644); err != nil {
			t.Fatal(err)
		}
	} else if err != nil {
		t.Fatal(err)
	}

	filename := filepath.Join("testdata", "config.yaml")
	data, err := yaml.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		t.Fatal(err)
	}

	config, err := Parse(filename)
	if err != nil {
		t.Fatal(err)
	}
	if config.Server.Port != expected.Server.Port {
		t.Errorf("expected server port [%v] actual [%v]", expected.Server.Port, config.Server.Port)
	}
	if config.Logging.Level != expected.Logging.Level {
		t.Errorf("expected log level [%v] actual [%v]", expected.Logging.Level, config.Logging.Level)
	}
	if config.Logging.File != expected.Logging.File {
		t.Errorf("expected log file [%v] actual [%v]", expected.Logging.File, config.Logging.File)
	}
	if config.Logging.MaxAge != expected.Logging.MaxAge {
		t.Errorf("expected log file age [%v] actual [%v]", expected.Logging.MaxAge, config.Logging.MaxAge)
	}
	if config.Logging.MaxSize != expected.Logging.MaxSize {
		t.Errorf("expected log file size [%v] actual [%v]", expected.Logging.MaxSize, config.Logging.MaxSize)
	}
}
