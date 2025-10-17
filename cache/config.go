package cache

import (
	"github.com/tandy9527/js-util/utils"
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

func LoadRedisConf(path string) *RedisMap {

	return utils.Loadyaml[RedisMap](path)
}
