package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultLogLevel                = "debug"
	defaultLogFile                 = "log.log"
	defaultLogFileAge              = 7
	defaultLogFileSize             = 1024
	defaultProcessTestDirectory    = "/opt/deployments/process/test/"
	defaultProcessReleaseDirectory = "/opt/deployments/process/release/"
)

type Config struct {
	Server  Server  `yaml:"server"`
	Logging Logging `yaml:"logging"`
	Process Process `yaml:"process"`
}

type Server struct {
	Port int `yaml:"port"`
}

type Logging struct {
	Level   string `yaml:"level"`
	File    string `yaml:"file"`
	MaxAge  int    `yaml:"max-age"`
	MaxSize int    `yaml:"max-size"`
}

type Process struct {
	Directory ProcessDirectory `yaml:"directory"`
}

type ProcessDirectory struct {
	Test    string `yaml:"test"`
	Release string `yaml:"release"`
}

func Parse(file string) (*Config, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(f, config)
	if err != nil {
		return nil, err
	}
	applyDefault(config)
	return config, nil
}

func applyDefault(config *Config) {
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Logging.Level == "" {
		config.Logging.Level = defaultLogLevel
	}
	if config.Logging.File == "" {
		config.Logging.File = defaultLogFile
	}
	if config.Logging.MaxAge == 0 {
		config.Logging.MaxAge = defaultLogFileAge
	}
	if config.Logging.MaxSize == 0 {
		config.Logging.MaxSize = defaultLogFileSize
	}
	if config.Process.Directory.Test == "" {
		config.Process.Directory.Test = defaultProcessTestDirectory
	}
	if config.Process.Directory.Release == "" {
		config.Process.Directory.Release = defaultProcessReleaseDirectory
	}
}
