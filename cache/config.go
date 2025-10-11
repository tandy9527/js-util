package cache

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RedisConf struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
	PoolTimeout  int    `yaml:"pool_timeout"`
}
type RedisMap struct {
	Redis map[string]RedisConf `yaml:"redis"`
}

func LoadRedisConf(path string) (*RedisMap, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg RedisMap
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
