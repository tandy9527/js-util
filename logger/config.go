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

func LoadLoggerConfig(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic("failed to read logger file: " + path)
	}
	var cfg LoggerConf
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		panic("failed to unmarshal logger file: " + path)
	}
	LoggerInit(cfg.LogFilePath, cfg.MaxSizeMB, cfg.MaxBackups, cfg.MaxAgeDays, cfg.Compress)
}
