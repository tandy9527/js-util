package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var sugar *zap.SugaredLogger
var coreLogger *zap.Logger

// Init 初始化 logger
// logFilePath: 日志文件路径（例如 "logs/gateway.log"）
// maxSizeMB: 单个日志文件最大 MB（lumberjack）
// maxBackups: 保留旧文件个数
// maxAgeDays: 保留旧文件天数
// compress: 是否压缩旧日志（gzip）
func LoggerInit(logFilePath string, maxSizeMB, maxBackups, maxAgeDays int, compress bool) {
	if logFilePath == "" {
		logFilePath = "game.log"
	}
	if maxSizeMB <= 0 {
		maxSizeMB = 100
	}
	if maxBackups <= 0 {
		maxBackups = 7
	}
	if maxAgeDays <= 0 {
		maxAgeDays = 7
	}

	level := zapcore.InfoLevel // 日志级别按需改为 DebugLevel

	// 时间编码（只显示时:分:秒）
	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	// 级别编码（大写）
	levelEncoder := func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(strings.ToUpper(l.String()))
	}

	// 自定义 caller 编码：取最后两级路径作为 pkg/file.go:line，并右填充固定宽度（便于对齐）
	const callerWidth = 25
	callerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		// 统一为 '/' 分隔
		rel := filepath.ToSlash(caller.File)
		parts := strings.Split(rel, "/")
		short := ""
		if len(parts) >= 2 {
			short = parts[len(parts)-2] + "/" + parts[len(parts)-1]
		} else {
			short = rel
		}
		short = short + ":" + strconv.Itoa(caller.Line)
		enc.AppendString(padRight(short, callerWidth))
	}

	// 控制台 encoder（彩色）
	consoleEncCfg := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		CallerKey:      "C",
		MessageKey:     "M",
		EncodeTime:     timeEncoder,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeCaller:   callerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
	}
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncCfg)

	// 文件 encoder（无色）
	fileEncCfg := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		CallerKey:      "C",
		MessageKey:     "M",
		EncodeTime:     timeEncoder,
		EncodeLevel:    levelEncoder, // 无颜色
		EncodeCaller:   callerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
	}
	fileEncoder := zapcore.NewConsoleEncoder(fileEncCfg)

	// lumberjack 切割器（写入文件）
	lumber := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    maxSizeMB,
		MaxBackups: maxBackups,
		MaxAge:     maxAgeDays,
		Compress:   compress,
	}

	// cores：console + file
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(lumber), level)

	core := zapcore.NewTee(consoleCore, fileCore)

	// 启用 caller（显示调用位置），AddStacktrace 保留 Error 以上的堆栈
	coreLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar = coreLogger.Sugar()
}

// padRight 右侧补空格对齐
func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// Sync 在程序退出时调用，确保缓冲写入磁盘
func Sync() {
	if coreLogger != nil {
		_ = coreLogger.Sync()
	}
}

func Infof(format string, args ...any) {
	if sugar != nil {
		sugar.WithOptions(zap.AddCallerSkip(1)).Infof(format, args...)
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

func Errorf(format string, args ...any) {
	if sugar != nil {
		sugar.WithOptions(zap.AddCallerSkip(1)).Errorf(format, args...)
	} else {
		fmt.Printf(format+"\n", args...)
	}
}
func Warnf(format string, args ...any) {
	if sugar != nil {
		sugar.WithOptions(zap.AddCallerSkip(1)).Warnf(format, args...)
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

func Info(format string, args ...any) {
	if coreLogger != nil {
		var fields []zap.Field
		var msgArgs []any
		for _, a := range args {
			if f, ok := a.(zap.Field); ok {
				fields = append(fields, f)
			} else {
				msgArgs = append(msgArgs, a)
			}
		}
		msg := fmt.Sprintf(format, msgArgs...)
		coreLogger.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
	} else {
		fmt.Printf("[INFO] "+format+"\n", args...)
	}
}
func Error(format string, args ...any) {
	if coreLogger != nil {
		var fields []zap.Field
		var msgArgs []any
		for _, a := range args {
			if f, ok := a.(zap.Field); ok {
				fields = append(fields, f)
			} else {
				msgArgs = append(msgArgs, a)
			}
		}
		msg := fmt.Sprintf(format, msgArgs...)
		coreLogger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
	} else {
		fmt.Printf("[ERROR] "+format+"\n", args...)
	}
}
func Warn(format string, args ...any) {
	if coreLogger != nil {
		var fields []zap.Field
		var msgArgs []any
		for _, a := range args {
			if f, ok := a.(zap.Field); ok {
				fields = append(fields, f)
			} else {
				msgArgs = append(msgArgs, a)
			}
		}
		msg := fmt.Sprintf(format, msgArgs...)
		coreLogger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
	} else {
		fmt.Printf("[WARN] "+format+"\n", args...)
	}
}

func Debug(format string, args ...any) {
	if coreLogger != nil {
		var fields []zap.Field
		var msgArgs []any
		for _, a := range args {
			if f, ok := a.(zap.Field); ok {
				fields = append(fields, f)
			} else {
				msgArgs = append(msgArgs, a)
			}
		}
		msg := fmt.Sprintf(format, msgArgs...)
		coreLogger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
	} else {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}
