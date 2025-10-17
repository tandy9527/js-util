package tools

import (
	"os"

	"github.com/tandy9527/js-util/logger"
	"gopkg.in/yaml.v3"
)

// Loadyaml 加载.yaml
func Loadyaml[T any](filePath string) *T {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic("failed to read config file: " + filePath)
	}
	var config T
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic("failed to unmarshal config file: " + filePath)
	}

	logger.Infof("[Loadyaml] load successful: %s", filePath)
	return &config
}
