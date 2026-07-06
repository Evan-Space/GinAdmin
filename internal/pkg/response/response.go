package response

import (
	"net/http"
	"reflect"
	"time"

	"GinAdmin/internal/global"
	"GinAdmin/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)



type Result struct {
	Code    int         `json:"code"`
	Msg 	string      `json:"msg"`
	Data    interface{} `json:"data"`
	Cost    string      `json:"cost"`
	RequestId string    `json:"request_id"`
}


func NewResult() *Result {
	return &Result{
		Code: 0,
		Msg: "",
		Data: emptyObject(),
		Cost: "",
		RequestId: "",
	}
}

// 响应处理器类型
type Response struct {
	httpCode int
	result *Result
	msgKey string
	msgArgs []any
}


// 响应处理器构造函数
func Resp() *Response {
	return &Response{
		httpCode: http.StatusOK,
		result: NewResult(),
	}
}


// Fail 错误的返回
func (r *Response) Fail(c *gin.Context, code int, msg string, data ...any) {
	r.SetCode(code)
	r.SetMessage(msg)
	if len(data) > 0 && data[0] != nil {
		r.WithData(data[0])
	}

	r.json(c)
}


/**
* 自定义返回错误码
**/
func (r *Response) FailCode(c *gin.Context, code int, msg ...string) {
	r.SetCode(code)
	if len(msg) > 0 && msg[0] != "" {
		r.SetMessage(msg[0])
	}
	r.json(c)
}


/**
* 空对象
**/
func emptyObject() map[string]any {
	return map[string]any{}
}


func (r *Response) SetCode(code int) *Response {
	r.result.Code = code
	return r
}

func (r *Response) SetMessage(message string) *Response {
	r.result.Msg = message
	r.msgKey = ""
	r.msgArgs = nil
	return r
}


type defaultRes struct {
	Result any `json:"result"`
}

/**
* 设置返回data数据
**/
func (r *Response) WithData(data any) *Response {
	if isNilData(data) {
		r.result.Data = emptyObject()
		return r
	}

	if !isObjectData(data) {
		r.result.Data = &defaultRes{ Result: data }
		return r
	}

	r.result.Data = data
	return r
}

func isNilData(data any) bool {
	if data == nil {
		return true
	}
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func isObjectData(data any) bool {
	value := reflect.ValueOf(data)

	for value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return false
		}
		value = value.Elem()
	}

	switch value.Kind() {
		case reflect.Struct, reflect.Map:
			return true
		default:
			return false
	}

}



// json 返回 gin 框架的 JSON 响应
func (r *Response) json(c *gin.Context) {

	// 计算请求耗时
	r.result.Cost = time.Since(c.GetTime(global.ContextKeyRequestStartTime)).String()
	r.result.RequestId = c.GetString(global.ContextKeyRequestID)
	c.AbortWithStatusJSON(r.httpCode, r.result)
}

// 成功的返回
func Ok(c *gin.Context, data any, msg ...string) {
	r := Resp().WithData(data)
	if len(msg) > 0 && msg[0] != "" {
		r.SetMessage(msg[0])
	}
	r.json(c)
}


func (r *Response) Success(c *gin.Context) {
	r.SetCode(errors.SUCCESS)
	r.json(c)
}

func (r *Response) WithDataSuccess(c *gin.Context, data interface{}) {
	r.SetCode(errors.SUCCESS)
	r.WithData(data)
	r.json(c)
}


// Success 业务响应成功（便捷响应方法）
func Success(c *gin.Context, data ...any) {
	if len(data) > 0 && data[0] != nil {
		Resp().WithDataSuccess(c, data[0])
		return
	}
	Resp().Success(c)
}


// FailCode 业务失败响应（便捷响应方法）
func FailCode(c *gin.Context, code int, data ...any) {
	if len(data) > 0 && data[0] != nil {
		Resp().WithData(data[0]).FailCode(c, code)
		return
	}
	Resp().FailCode(c, code)
}


func Fail(c *gin.Context, code int, message string, data ...any) {
	if len(data) > 0 && data[0] != nil {
		Resp().WithData(data[0]).Fail(c, code, message)
		return
	}
	Resp().Fail(c, code, message)
}

// Err 统一错误处理 
// 判断错误类型是自定义类型则自动返回错误中携带的code和message，否则返回服务器错误
func Err(c *gin.Context, err error) {
	helper := errors.ErrorHelper{}
	businessError, parseErr := helper.AsBusinessError(err)
	if parseErr != nil {
		FailCode(c, errors.ServerErr)
		return
	}
	// BusinessError，使用其 code 和 message
	Fail(c, businessError.GetCode(), businessError.GetMessage())
}