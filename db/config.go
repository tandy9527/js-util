package db

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbName"`
	Charset         string        `yaml:"charset"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`    // 最大连接数 0 -无限制
	MaxIdleConns    int           `yaml:"maxIdleConns"`    // 最大空闲连接
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"` // 单个连接最大存活时间 m 分
	Timeout         time.Duration `yaml:"timeout"`         // 默认超时时间 s 秒
	LogSQL          bool          `yaml:"logSQL"`          // 是否打印 SQL 日志
}

type DBConfig struct {
	DB Config `yaml:"mysql"`
}

func LoadMySQLConf(path string) (*DBConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg DBConfig
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
