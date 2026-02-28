package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


type Field = zap.Field

func Str(key, val string) Field        { return zap.String(key, val) }
func Int(key string, val int) Field     { return zap.Int(key, val) }
func Uint(key string, val uint) Field   { return zap.Uint(key, val) }
func Int64(key string, val int64) Field { return zap.Int64(key, val) }
func Float64(key string, val float64) Field { return zap.Float64(key, val) }
func Bool(key string, val bool) Field   { return zap.Bool(key, val) }
func Err(err error) Field              { return zap.Error(err) }
func Any(key string, val interface{}) Field { return zap.Any(key, val) }
func Dur(key string, val time.Duration) Field { return zap.Duration(key, val) }

var l *zap.Logger
var s *zap.SugaredLogger

// Config 日志配置
type Config struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"file_path"`
}

// Init 初始化，main.go 里调用一次
func Init(cfg Config) error {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
			level = zapcore.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	var cores []zapcore.Core

	stdoutCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.AddSync(os.Stdout),
			level,
	)
	cores = append(cores, stdoutCore)

	if cfg.FilePath != "" {
			file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
			if err != nil {
					return err
			}
			fileCore := zapcore.NewCore(
					zapcore.NewJSONEncoder(encoderCfg),
					zapcore.AddSync(file),
					level,
			)
			cores = append(cores, fileCore)
	}

	l = zap.New(
			zapcore.NewTee(cores...),
			zap.AddCaller(),
			zap.AddCallerSkip(1), // 跳过本包的包装层，显示真实调用方
			zap.AddStacktrace(zap.ErrorLevel),
	)
	s = l.Sugar()

	return nil
}

// ---- 包级别日志函数，直接调用 ----

func Debug(msg string, fields ...Field) { l.Debug(msg, fields...) }
func Info(msg string, fields ...Field)  { l.Info(msg, fields...) }
func Warn(msg string, fields ...Field)  { l.Warn(msg, fields...) }
func Error(msg string, fields ...Field) { l.Error(msg, fields...) }
func Fatal(msg string, fields ...Field) { l.Fatal(msg, fields...) }

// printf 风格，偶尔用
func Debugf(template string, args ...interface{}) { s.Debugf(template, args...) }
func Infof(template string, args ...interface{})  { s.Infof(template, args...) }
func Warnf(template string, args ...interface{})  { s.Warnf(template, args...) }
func Errorf(template string, args ...interface{}) { s.Errorf(template, args...) }
func Fatalf(template string, args ...interface{}) { s.Fatalf(template, args...) }

// Sync 程序退出前调用
func Sync() {
	if l != nil {
			_ = l.Sync()
	}
}
