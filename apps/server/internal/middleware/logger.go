package middleware

import (
	"time"

	"GinAdmin/internal/global"
	"GinAdmin/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()       // ← 请求开始时间
		c.Next()                  // ← 执行后续中间件和 handler
		cost := time.Since(start) // 请求总耗时

		// 跳过 ping 和 404
		if c.Request.URL.Path == "/ping" || c.Writer.Status() == 404 {
			return
		}
		logger.Logger.Info("request",
			zap.String("requestId", c.GetString(global.ContextKeyRequestID)),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path), 
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("cost", cost),
		)

	}
}
