package config

import (
	"encoding/json"
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
	Server  Server  `yaml:"server" json:"server"`
	Logging Logging `yaml:"logging" json:"logging"`
}

type Server struct {
	Port int `yaml:"port" json:"port"`
}

type Logging struct {
	Level   string `yaml:"level" json:"level"`
	File    string `yaml:"file" json:"file"`
	MaxAge  int    `yaml:"max-age" json:"maxAge"`
	MaxSize int    `yaml:"max-size" json:"maxSize"`
}

func (c *Config) String() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(bytes)
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
	if config.Server.Port <= 0 {
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
