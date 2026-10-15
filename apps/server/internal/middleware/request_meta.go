package middleware

import (
	"GinAdmin/internal/global"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


func RequestMeta() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(global.ContextKeyRequestID, uuid.New().String())
		c.Set(global.ContextKeyRequestStartTime, time.Now())
		c.Next()
	}
}
