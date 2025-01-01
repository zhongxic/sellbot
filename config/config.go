package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultServerPort                  = 8080
	defaultLogLevel                    = "debug"
	defaultLogFile                     = "log.log"
	defaultLogFileAge                  = 7
	defaultLogFileSize                 = 1024
	defaultProcessCacheExpiration      = 1800
	defaultProcessCacheCleanupInterval = 900
	defaultTestProcessDirectory        = "data/process/test/"
	defaultReleaseProcessTestDirectory = "data/process/release/"
	defaultSessionRepository           = "memory"
	defaultSessionExpiration           = 1800
)

type Config struct {
	Server    Server    `yaml:"server"`
	Logging   Logging   `yaml:"logging"`
	Process   Process   `yaml:"process"`
	Tokenizer Tokenizer `yaml:"tokenizer"`
	Session   Session   `yaml:"session"`
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

type Cache struct {
	Expiration      int `yaml:"expiration"`
	CleanupInterval int `yaml:"cleanup-interval"`
}

type Directory struct {
	Test    string `yaml:"test"`
	Release string `yaml:"release"`
}

type Tokenizer struct {
	ExtraDict     string `yaml:"extra-dict"`
	StopWordsDict string `yaml:"stop-words-dict"`
}

type Session struct {
	Repository string `yaml:"repository"`
	Expiration int    `yaml:"expiration"`
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
		config.Server.Port = defaultServerPort
	}
	if config.Logging.Level == "" {
		config.Logging.Level = defaultLogLevel
	}
	if config.Logging.File == "" {
		config.Logging.File = defaultLogFile
	}
	if config.Logging.MaxAge <= 0 {
		config.Logging.MaxAge = defaultLogFileAge
	}
	if config.Logging.MaxSize <= 0 {
		config.Logging.MaxSize = defaultLogFileSize
	}
	if config.Process.Cache.Expiration <= 0 {
		config.Process.Cache.Expiration = defaultProcessCacheExpiration
	}
	if config.Process.Cache.CleanupInterval <= 0 {
		config.Process.Cache.CleanupInterval = defaultProcessCacheCleanupInterval
	}
	if config.Process.Directory.Test == "" {
		config.Process.Directory.Test = defaultTestProcessDirectory
	}
	if config.Process.Directory.Release == "" {
		config.Process.Directory.Release = defaultReleaseProcessTestDirectory
	}
	if config.Session.Repository == "" {
		config.Session.Repository = defaultSessionRepository
	}
	if config.Session.Expiration <= 0 {
		config.Session.Expiration = defaultSessionExpiration
	}
}
