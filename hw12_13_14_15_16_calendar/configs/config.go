package config

import (
	"fmt"
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Calendar  CalendarConfig
	Scheduler SchedulerConfig
	Sender    SenderConfig
}

type CalendarConfig struct {
	Server  ServerConf  `yaml:"server"`
	Logger  LoggerConf  `yaml:"logger"`
	Storage StorageConf `yaml:"storage"`
}

type SchedulerConfig struct {
	AMQP    AMQPConf    `yaml:"amqp"`
	Logger  LoggerConf  `yaml:"logger"`
	Storage StorageConf `yaml:"storage"`
}

type SenderConfig struct {
	AMQP    AMQPConf    `yaml:"amqp"`
	Logger  LoggerConf  `yaml:"logger"`
	Storage StorageConf `yaml:"storage"`
}

type ServerConf struct {
	Type    string        `yaml:"type"`
	Host    string        `yaml:"host"`
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type AMQPConf struct {
	URL   string `yaml:"url"`
	Topic string `yaml:"topic"`
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
	if c.Calendar.Server.Port == "" {
		return fmt.Errorf("port is not specified")
	}
	if c.Calendar.Server.Host == "" {
		return fmt.Errorf("host is not specified")
	}
	if c.Calendar.Storage.Type == "" {
		return fmt.Errorf("storage is not specified")
	}
	if c.Calendar.Server.Timeout == 0 {
		c.Calendar.Server.Timeout = 30 * time.Second
	}

	if c.Scheduler.AMQP.URL == "" {
		return fmt.Errorf("amqp url is not specified for scheduler")
	}

	if c.Sender.AMQP.URL == "" {
		return fmt.Errorf("amqp url is not specified for sender")
	}

	return nil
}
