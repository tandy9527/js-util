package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/tandy9527/js-util/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once // 保证只初始化一次
)

func LoadMysql(path string) {
	cfg := LoadMySQLConf(path)
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)
		logger.Infof("LoadMysql DSN: %s", dsn)
		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Sprintf("msyql Connection failed: %v", err))
		}
		sqlDB, _ := DB.DB()
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Hour)
		logger.Infof("LoadMysql successful")

	})
}

// CloseMySQL 关闭连接
func CloseMySQL() {
	if DB != nil {
		sqlDB, _ := DB.DB()
		_ = sqlDB.Close()
		logger.Infof("close MySQL ")
	}
}

// CallProcedure 执行存储过程并将结果映射到 dest（结构体或切片）
// procName: 存储过程名
// dest: 指针类型，结构体或切片，用于接收结果集
// args: 输入参数
func CallProcedure(dest any, procName string, args ...any) error {
	// 构造 CALL 语句，例如 CALL procName(?, ?, ?)
	callStmt := fmt.Sprintf("CALL %s(%s)", procName, placeholders(len(args)))
	return DB.Raw(callStmt, args...).Scan(dest).Error
}

func placeholders(n int) string {
	if n == 0 {
		return ""
	}
	s := "?"
	for i := 1; i < n; i++ {
		s += ",?"
	}
	return s
}
