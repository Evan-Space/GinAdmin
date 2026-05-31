package routers

import (
	"github.com/gin-gonic/gin"
)

// SetRouters 创建 Gin 引擎并注册全部应用路由。
func SetRouters() (*gin.Engine, error) {
	engine := gin.Default()

	// 模拟注册路由
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// mock 没有错误
	return engine, nil
}

