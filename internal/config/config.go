package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Watch         WatchConfig         `yaml:"watch"`
	Report        ReportConfig        `yaml:"report"`
	Notifications NotificationsConfig `yaml:"notifications"`
	Server        ServerConfig        `yaml:"server"`
}

type WatchConfig struct {
	Paths  []string `yaml:"paths"`
	Ignore []string `yaml:"ignore"`
}

type ReportConfig struct {
	Formats   []string `yaml:"formats"`
	OutputDir string   `yaml:"output_dir"`
}

type NotificationsConfig struct {
	Enable   bool     `yaml:"enable"`
	Channels []string `yaml:"channels"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
