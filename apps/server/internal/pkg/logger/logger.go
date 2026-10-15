package logger

import (
	"GinAdmin/config"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	nopLogger = zap.NewNop() // 空操作 logger， 未初始化时不会 panic
	Logger = nopLogger // 包级全局 logger 所有代码通过它写入日志
	once sync.Once // 保证只初始化一次
	initErr error
)

// InitLogger 初始化全局 zap logger (只执行一次)
func InitLogger() error {
	once.Do(func() {
		logger, err := buildLogger(config.GetConfig())
		if err != nil {
			initErr = fmt.Errorf("初始化 logger 失败", err)
			return
		}
		Logger = logger
	})
	return initErr
}


func buildLogger(cfg *config.Conf) (*zap.Logger, error) {
	// stderr 模式，直接输出到终端，开发环境下使用
	if cfg.Logger.Output == "stderr" {
		return zap.NewDevelopment()
	}

	// ====== 1. 定义日志格式 ======
	encoderConfig := zap.NewProductionEncoderConfig()
	// 时间格式为： 2006-01-02 15:04:05.000
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	

	// ====== 2. 确定日志文件路径 ======
	// 确定日志文件路径
	logDir := filepath.Join(cfg.BasePath, "logs")
	os.MkdirAll(logDir, 0755)
	filename := filepath.Join(logDir, cfg.Logger.Filename)

	// ====== 3. 确定日志文件切割策略 ======
	var writer io.Writer
	if cfg.Logger.DefaultDivision == "size" {
		// 按大小切割：超过 MaxSize MB 自动切断文件
		writer = &lumberjack.Logger {
			Filename: filename,
			MaxSize: cfg.Logger.DivisionSize.MaxSize,
			MaxBackups: cfg.Logger.DivisionSize.MaxBackups,
			MaxAge: cfg.Logger.DivisionSize.MaxAge,
			Compress: cfg.Logger.DivisionSize.Compress,
		}
	} else {
		// 按时间切割， 每天一个文件，格式 filename.YYYY-MM-DD.log
		hook, err := rotatelogs.New(
			filename+".%Y%m%d",
			rotatelogs.WithLinkName(filename),
			rotatelogs.WithMaxAge(time.Hour*24*time.Duration(cfg.Logger.DivisionTime.MaxAge)),
			rotatelogs.WithRotationTime(time.Hour*time.Duration(cfg.Logger.DivisionTime.RotationTime)),
		)
		if err != nil {
			return nil, fmt.Errorf("创建时间切割文件失败: %w", err)
		}
		writer = hook
	}
	// ====== 4. 组合 encoder + writer → 完整的 logger ======
	core := zapcore.NewCore(encoder, zapcore.AddSync(writer), zap.InfoLevel)
	return zap.New(core, zap.AddCaller()), nil // AddCaller 在每条日志后附加调用位置
}