package db

import "time"

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string        // 连接库名
	MaxOpen  int           // 最大连接数 0 -无限制
	MaxIdle  int           // 最大空闲连接
	Timeout  time.Duration // 默认超时时间 单位：秒
	LogSQL   bool          // 是否打印 SQL 日志
}
