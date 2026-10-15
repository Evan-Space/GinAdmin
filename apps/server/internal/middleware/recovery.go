package middleware

import (
	"GinAdmin/internal/global"
	"GinAdmin/internal/pkg/errors"
	"GinAdmin/internal/pkg/logger"
	"GinAdmin/internal/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CustomRecovery() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(
		&panicWriter{},
		func(c *gin.Context, err any) {
			// ----- 1. 记录结构化日志 -----
			logger.Logger.Error("panic recovered",
				zap.String("requestId", c.GetString(global.ContextKeyRequestID)),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("query", c.Request.URL.RawQuery),
				zap.String("ip", c.ClientIP()),
				zap.Any("error", err), // zap.Any 可以打印任意类型的 panic 值
			)

			// ----- 2. 返回统一错误 JSON -----
			response.Resp().
				FailCode(c, errors.ServerErr) // {"code":500,"msg":"服务器内部错误"}
		},
	)
}

// panicWriter 实现 io.Write 接口
// Gin 会把 panic 的完整调用栈（stack trace） 写到这个  Write
type panicWriter struct{}

func (p *panicWriter) Write(b []byte) (int, error) {
	// 调用战也通过 logger 输出
	logger.Logger.Error(strings.TrimSpace(string(b)))
	return len(b), nil
}
