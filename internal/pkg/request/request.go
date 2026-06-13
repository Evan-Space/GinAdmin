package request

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

/**
*  GetQueryParams 提取当前请求的查询参数。
**/
func GetQueryParams(c *gin.Context) map[string]any {
	if c == nil {
		return map[string]any{}
	}
	query := c.Request.URL.Query()
	var queryMap = make(map[string]any, len(query))
	for k := range query {
		queryMap[k] = c.Query(k)
	}
	return queryMap
}

func GetAccessToken(c *gin.Context) (string, error) {
	if c == nil {
		return "", errors.New("gin context is nil")
	}
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		return "", errors.New("authorization header is missing")
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("invalid Authorization header format")
	}

	return strings.TrimPrefix(authorization, "Bearer "), nil
}