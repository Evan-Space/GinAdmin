package logger

import (
	"GinAdmin/config"
	"fmt"
	"io"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)


var (
	nopLogger = zap.NewNop()
	Logger    = nopLogger

	loggerOnce sync.Once
	loggerVal  atomic.Value
	loggerMu   sync.Mutex
	initErr    error
)

/**
* 初始化日志
**/
func InitLogger() error {
	loggerOnce.Do(func() {
		logger, err := buildLogger(config.GetConfig())
		if err != nil {
			initErr = fmt.Errorf("创建zap日志包失败: %w", err)
			return
		}
		setLogger(logger)
	})
	return initErr
}


func buildLogger(cfg *config.Conf) (*zap.Logger, error) {
	if cfg == nil {
		cfg = config.GetConfig()
	}

	if cfg.Logger.Output == "stderr" {
		return zap.NewDevelopment()
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	filename := filepath.Join(cfg.BasePath, "logs", cfg.Logger.Filename)

	var writer zapcore.WriteSyncer
	if cfg.Logger.DefaultDivision == "size" {
		writer = zapcore.AddSync(getLumberJackWriter(cfg, filename))
	} else {
		rotateWriter, err := getRotateWriter(cfg, filename)
		if err != nil {
			return nil, err
		}
		writer = zapcore.AddSync(rotateWriter)
	}

	zapCore := zapcore.NewCore(encoder, writer, zap.InfoLevel)
	return zap.New(zapCore, zap.AddCaller()), nil
}


/**
* 根据新配置重新建立全局日志实例
**/
func RebuildLogger(cfg *config.Conf) error {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	next, err := buildLogger(cfg)
	if err != nil {
		return err
	}
	old := current()
	setLogger(next)
	if old != nil {
		_ = old.Sync()
	}
	return nil
}


func current() *zap.Logger {
	if logger, ok := loggerVal.Load().(*zap.Logger); ok && logger != nil {
		return logger
	}
	if Logger != nil {
		return Logger
	}

	return nopLogger
}


func setLogger(logger *zap.Logger) {
	if logger == nil {
		logger = nopLogger
	}

	Logger = logger
	loggerVal.Store(logger)
}


func ReplaceLoggerForTesting(logger *zap.Logger) func() {
	loggerMu.Lock()
	previous := current()
	setLogger(logger)
	loggerMu.Unlock()

	return func() {
		loggerMu.Lock()
		setLogger(previous)
		loggerMu.Unlock()
	}
}



/**
* 记录info级别日志
**/
func Info(msg string, fields ...zap.Field) {
	current().Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	current().Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	current().Warn(msg, fields...)
}




// getRotateWriter 按日期切割日志
func getRotateWriter(cfg *config.Conf, filename string) (io.Writer, error) {
	maxAge := time.Duration(cfg.Logger.DivisionTime.MaxAge)
	rotationTime := time.Duration(cfg.Logger.DivisionTime.RotationTime)
	hook, err := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*maxAge),
		rotatelogs.WithRotationTime(time.Hour*rotationTime),
	)
	if err != nil {
		return nil, err
	}
	return hook, nil
}

// getLumberJackWriter 按文件切割日志
func getLumberJackWriter(cfg *config.Conf, filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cfg.Logger.DivisionSize.MaxSize,
		MaxBackups: cfg.Logger.DivisionSize.MaxBackups,
		MaxAge:     cfg.Logger.DivisionSize.MaxAge,
		Compress:   cfg.Logger.DivisionSize.Compress,
	}
}
