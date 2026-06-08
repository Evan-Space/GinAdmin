package response

import (
	"net/http"
	"reflect"
	"time"

	"GinAdmin/global"

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
		r.Withdata(data[0])
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
func (r *Response) Withdata(data any) *Response {
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
	// 如果消息为空，使用错误码对应的默认消息

	// 计算请求耗时
	r.result.Cost = time.Since(c.GetTime(global.ContextKeyRequestStartTime)).String()
	r.result.RequestId = c.GetString(global.ContextKeyRequestID)
	c.AbortWithStatusJSON(r.httpCode, r.result)
}




// 成功的返回
func Ok(c *gin.Context, data any, msg ...string) {
	r := Resp().Withdata(data)
	if len(msg) > 0 && msg[0] != "" {
		r.SetMessage(msg[0])
	}
	r.json(c)
}

// Success 业务成功响应（便捷方法）
// func Success(c *gin.Context, data ...any) {
// 	if len(data) > 0 && data[0] != nil {
// 		Resp().WithDataSuccess(c, data[0])
// 		return
// 	}
// 	Resp().Success(c)
// }

// Success 正确返回
// func (r *Response) Success(c *gin.Context) {
// 	r.SetCode(errors.SUCCESS)
// 	r.json(c)
// }

// // WithDataSuccess 成功后需要返回值
// func (r *Response) WithDataSuccess(c *gin.Context, data interface{}) {
// 	r.SetCode(errors.SUCCESS)
// 	r.WithData(data)
// 	r.json(c)
// }