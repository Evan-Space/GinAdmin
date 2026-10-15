package validator

import (
	"GinAdmin/internal/pkg/response"
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// InitValidator 初始化 Gin 的 validator 引擎
func InitValidator() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			if label := field.Tag.Get("json"); label != "" && label != "-" {
				return label
			}

			return field.Name
		})
	}
	return nil
}


// CheckPostParams 检查 POST JSON 请求参数，校验失败自动响应错误
func CheckPostParams(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBind(obj); err != nil {
		handleBindError(c, err)
		return err
	}
	return nil
}


// CheckQueryParams 检查 GET 查询参数，校验失败自动响应错误
func CheckQueryParams(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		handleBindError(c, err)
		return err
	}
	return nil
}

func handleBindError(c *gin.Context, err error) {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		msgs := make([]string, 0, len(validationErrs))
		for _, fieldErr := range validationErrs {
			msgs = append(msgs, fieldErr.Field()+": "+fieldErr.Tag())
		}
		response.Resp().Fail(c, 10000,strings.Join(msgs, ","))
		return
	}

	// JSON 格式/类型错误
	var typeErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError
	if errors.As(err, &typeErr) || errors.As(err, &syntaxErr) {
		response.Resp().FailCode(c, 10000)
		return
	}

	response.Resp().FailCode(c, 10000)
}