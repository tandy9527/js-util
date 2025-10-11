package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/tandy9527/js-util/logger"
)

var (
	mysql *MySQLClient
	once  sync.Once
)

type MySQLClient struct {
	db      *sql.DB
	timeout time.Duration
	logSQL  bool
}

// 初始化单例
func Init(cfg Config) (*MySQLClient, error) {
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

		db, e := sql.Open("mysql", dsn)
		if e != nil {
			err = e
			return
		}
		db.SetMaxOpenConns(cfg.MaxOpenConns)
		db.SetMaxIdleConns(cfg.MaxIdleConns)
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

		mysql = &MySQLClient{
			db:      db,
			timeout: cfg.Timeout,
			logSQL:  cfg.LogSQL,
		}
		logger.Info("[MySQL] Connect to %s:%d/%s", cfg.Host, cfg.Port, cfg.DBName)
	})
	return mysql, err
}

// 获取单例
func GetInstance() *MySQLClient {
	if mysql == nil {
		logger.Error("[MySQL]  Instance not initialized, call Init() first.")
		return nil
	}
	return mysql
}

// ===================================================
// 模式 1：带 ctx
// ===================================================

func (c *MySQLClient) ExecCtx(ctx context.Context, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	res, err := c.db.ExecContext(ctx, query, args...)
	c.logQuery("ExecCtx", query, time.Since(start), err)
	return res, err
}

func (c *MySQLClient) QueryCtx(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := c.db.QueryContext(ctx, query, args...)
	c.logQuery("QueryCtx", query, time.Since(start), err)
	return rows, err
}

func (c *MySQLClient) QueryRowCtx(ctx context.Context, query string, args ...any) *sql.Row {
	start := time.Now()
	row := c.db.QueryRowContext(ctx, query, args...)
	c.logQuery("QueryRowCtx", query, time.Since(start), nil)
	return row
}

// ===================================================
// 模式 2：默认 ctx
// ===================================================

func (c *MySQLClient) defaultCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.timeout)
}

func (c *MySQLClient) Exec(query string, args ...any) (sql.Result, error) {
	ctx, cancel := c.defaultCtx()
	defer cancel()
	return c.ExecCtx(ctx, query, args...)
}

func (c *MySQLClient) Query(query string, args ...any) (*sql.Rows, error) {
	ctx, cancel := c.defaultCtx()
	defer cancel()
	return c.QueryCtx(ctx, query, args...)
}

func (c *MySQLClient) QueryRow(query string, args ...any) *sql.Row {
	ctx, cancel := c.defaultCtx()
	defer cancel()
	return c.QueryRowCtx(ctx, query, args...)
}

// ===================================================
// 事务封装
// ===================================================

func (c *MySQLClient) Transaction(fn func(tx *sql.Tx) error) error {
	ctx, cancel := c.defaultCtx()
	defer cancel()

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// ===================================================
// 存储过程调用
// ===================================================

// CallProcCtx：带 ctx 调用存储过程
func (c *MySQLClient) CallProcCtx(ctx context.Context, procName string, args ...any) (*sql.Rows, error) {
	placeholders := ""
	if len(args) > 0 {
		placeholders = "?" + strings.Repeat(",?", len(args)-1)
	}
	query := fmt.Sprintf("CALL %s(%s)", procName, placeholders)
	start := time.Now()
	rows, err := c.db.QueryContext(ctx, query, args...)
	c.logQuery("CallProcCtx", query, time.Since(start), err)
	return rows, err
}

// CallProc：使用默认 ctx 调用存储过程
func (c *MySQLClient) CallProc(procName string, args ...any) (*sql.Rows, error) {
	ctx, cancel := c.defaultCtx()
	defer cancel()
	return c.CallProcCtx(ctx, procName, args...)
}

// ===================================================
// 内部日志
// ===================================================

func (c *MySQLClient) logQuery(tag, query string, duration time.Duration, err error) {
	if !c.logSQL {
		return
	}
	if err != nil {
		logger.Error("[MySQL]  %s | %s | err=%v | cost=%v", tag, query, err, duration)
	} else {
		logger.Info("[MySQL]  %s | %s | cost=%v", tag, query, duration)
	}
}
func (c *MySQLClient) Close() error {
	if c.db != nil {
		err := c.db.Close()
		if err != nil {
			logger.Error("[MySQL]  Close error: %v", err)
			return err
		}
		logger.Info("[MySQL]  Connection closed.")
	}
	return nil
}
