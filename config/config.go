package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultLogLevel    = "debug"
	defaultLogFile     = "log.log"
	defaultLogFileAge  = 7
	defaultLogFileSize = 1024
)

type Config struct {
	Server    Server    `yaml:"server"`
	Logging   Logging   `yaml:"logging"`
	Process   Process   `yaml:"process"`
	Tokenizer Tokenizer `yaml:"tokenizer"`
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
	Cache     Cache     `yaml:"cache"`
	Directory Directory `yaml:"directory"`
}

type Directory struct {
	Test    string `yaml:"test"`
	Release string `yaml:"release"`
}

type Cache struct {
	Expiration      int `yaml:"expiration"`
	CleanupInterval int `yaml:"cleanup-interval"`
}

type Tokenizer struct {
	DictFile string `yaml:"dict-file"`
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
}
