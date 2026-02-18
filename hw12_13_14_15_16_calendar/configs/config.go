package config

import (
	"fmt"
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Calendar  CalendarConfig  `yaml:"calendar"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
	Sender    SenderConfig    `yaml:"sender"`
}

type CalendarConfig struct {
	Server  ServerConf  `yaml:"server"`
	Logger  LoggerConf  `yaml:"logger"`
	Storage StorageConf `yaml:"storage"`
}

type SchedulerConfig struct {
	AMQP     AMQPConf    `yaml:"amqp"`
	Logger   LoggerConf  `yaml:"logger"`
	Storage  StorageConf `yaml:"storage"`
	Interval string      `yaml:"interval"`
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

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var root map[string]any
	if err := yaml.Unmarshal(data, &root); err != nil {
		return Config{}, err
	}

	if err = yaml.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	if _, ok := root["calendar"]; ok {
		if err = config.validateCalendarFields(); err != nil {
			return Config{}, err
		}
	}

	if _, ok := root["scheduler"]; ok {
		if err = config.validateSchedulerFields(); err != nil {
			return Config{}, err
		}
	}

	if _, ok := root["sender"]; ok {
		if err = config.validateSenderFields(); err != nil {
			return Config{}, err
		}
	}

	return config, nil
}

func (c *Config) validateCalendarFields() error {
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

	return nil
}

func (c *Config) validateSchedulerFields() error {
	if c.Scheduler.AMQP.URL == "" {
		return fmt.Errorf("amqp url is not specified for scheduler")
	}

	if c.Scheduler.Interval == "" {
		c.Scheduler.Interval = "15m"
	}

	return nil
}

func (c *Config) validateSenderFields() error {
	if c.Sender.AMQP.URL == "" {
		return fmt.Errorf("amqp url is not specified for sender")
	}

	return nil
}
