package config

import (
	"fmt"
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConf  `yaml:"server"`
	Logger  LoggerConf  `yaml:"logger"`
	Storage StorageConf `yaml:"storage"`
}

type ServerConf struct {
	Type    string        `yaml:"type"`
	Host    string        `yaml:"host"`
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type StorageConf struct {
	Type string `yaml:"type"`
	DSN  string `yaml:"dsn"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	confFile, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer confFile.Close()

	decoder := yaml.NewDecoder(confFile)
	if err = decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	if err = config.validateFields(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c *Config) validateFields() error {
	if c.Server.Port == "" {
		return fmt.Errorf("port is not specified")
	}
	if c.Server.Host == "" {
		return fmt.Errorf("host are not specified")
	}
	if c.Server.Timeout == 0 {
		c.Server.Timeout = 30 * time.Second
	}
	return nil
}
