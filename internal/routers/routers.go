package routers

import (
	"GinAdmin/internal/controller"
	"GinAdmin/internal/pkg/response"
	"GinAdmin/middleware"

	"github.com/gin-gonic/gin"
)

// SetRouters 创建 Gin 引擎并注册全部应用路由。
func SetRouters() (*gin.Engine, error) {
	engine := gin.Default()

	engine.Use(middleware.RequestMeta()) // 注册中间件

	// 模拟注册路由
	engine.GET("/ping", func(c *gin.Context) {
		response.Ok(c, "pong")
	})


	userController := controller.NewUserController()

	api := engine.Group("/api/v1")
	{
		api.GET("/users", userController.List)
	}

	// mock 没有错误
	return engine, nil
}

