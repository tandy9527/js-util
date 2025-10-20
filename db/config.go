package db

import (
	"github.com/tandy9527/js-util/tools"
)

type MysqlConfig struct {
	// Host            string        `yaml:"host"`
	// Port            int           `yaml:"port"`
	// User            string        `yaml:"user"`
	// Password        string        `yaml:"password"`
	// DBName          string        `yaml:"dbName"`
	// Charset         string        `yaml:"charset"`
	// MaxOpenConns    int           `yaml:"maxOpenConns"`    // 最大连接数 0 -无限制
	// MaxIdleConns    int           `yaml:"maxIdleConns"`    // 最大空闲连接
	// ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"` // 单个连接最大存活时间 m 分
	// Timeout         time.Duration `yaml:"timeout"`         // 默认超时时间 s 秒
	// LogSQL          bool          `yaml:"logSQL"`          // 是否打印 SQL 日志

	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbName"`
	Charset         string `yaml:"charset"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`    // 最大打开连接数
	MaxIdleConns    int    `yaml:"maxIdleConns"`    // 最大空闲连接数
	ConnMaxLifetime string `yaml:"connMaxLifetime"` // 连接最大生命周期，使用字符串格式方便配置文件中设置
}

func LoadMySQLConf(path string) *MysqlConfig {
	return tools.Loadyaml[MysqlConfig](path)
}
