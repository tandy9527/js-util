package logger

import (
	"os"

	"gopkg.in/yaml.v3"
)

type LoggerConf struct {
	LogFilePath string `yaml:"logFilePath"`
	MaxSizeMB   int    `yaml:"maxSizeMB"`
	MaxBackups  int    `yaml:"maxBackups"`
	MaxAgeDays  int    `yaml:"maxAgeDays"`
	Compress    bool   `yaml:"compress"`
}

func LoadLoggerConfig(path string) (*LoggerConf, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg LoggerConf
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
