package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	Logger      LoggerConf    `yaml:"logger"`
	Storage     string        `yaml:"storage"`
	DBConnetion string        `yaml:"dbConnetion"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) ReadConfig(path string) error {
	confFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer confFile.Close()

	decoder := yaml.NewDecoder(confFile)
	if err = decoder.Decode(c); err != nil {
		return err
	}
	return c.ValidateConfig()
}

func (c *Config) ValidateConfig() error {
	if c.Port == "" {
		return fmt.Errorf("port is not specified")
	}
	if c.Host == "" {
		return fmt.Errorf("host are not specified")
	}
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}
	return nil
}
