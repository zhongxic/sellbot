package config

import (
	"os"
	"path/filepath"
	"reflect"
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
		Process: Process{
			Directory: ProcessDirectory{
				Test:    "/opt/deployments/process/test/",
				Release: "/opt/deployments/process/release/",
			},
		},
	}

	filename := filepath.Join("testdata", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(filename), 0644); err != nil {
		t.Fatal(err)
	}
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
	if !reflect.DeepEqual(config, expected) {
		t.Errorf("expected [%v] actual[%v]", expected, config)
	}
}
