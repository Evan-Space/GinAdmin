# Round 1: 登录认证 + 受保护接口 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现完整请求链路 `POST /api/v1/login` → JWT → `GET /api/v1/user/info`

**Architecture:** 参照 gin-layout 的声明式路由 + 中间件链 + controller/service 分层模式。全局 ParseToken 中间件静默解析 Token，controller 层按路由 AuthMode 决定是否拦截未登录请求。

**Tech Stack:** Gin, GORM, golang-jwt/jwt/v5, bcrypt, go-playground/validator/v10

---

## 文件结构规划

| 文件 | 操作 | 职责 | 对应 gin-layout 文件 |
|------|------|------|----------------------|
| `internal/validator/form/common.go` | 新建 | Paginate 通用分页表单 | `internal/validator/form/common.go` |
| `internal/pkg/errors/code.go` | 新建 | 错误码常量 | `internal/pkg/errors/code.go` |
| `internal/pkg/errors/error.go` | 新建 | BusinessError 类型 | `internal/pkg/errors/error.go` |
| `internal/pkg/errors/zh-cn.go` | 新建 | 错误码→中文消息映射 | `internal/pkg/errors/zh-cn.go` |
| `internal/pkg/request/request.go` | 修改 | 改造 GetAccessToken 对齐 gin-layout | `internal/pkg/request/request.go` |
| `internal/pkg/response/response.go` | 修改 | 激活 Success/WithDataSuccess/Err | `internal/pkg/response/response.go` |
| `internal/pkg/utils/jwt.go` | 新建 | JWT 生成/解析 | `internal/pkg/utils/token/jwt.go` |
| `internal/controller/sys_base.go` | 新建 | controller 基类 Api | `internal/controller/sys_base.go` |
| `internal/controller/login.go` | 修改 | 替换存根为真实登录逻辑 | `internal/controller/admin_v1/auth.go` |
| `internal/controller/user.go` | 修改 | 修复 List + 实现 Info | `internal/controller/admin_v1/auth_admin_user.go` |
| `internal/service/auth.go` | 新建 | 登录业务逻辑 | `internal/service/auth/login.go` |
| `internal/service/user.go` | 修改 | 新增 GetUserByID | `internal/service/admin/` |
| `internal/middleware/parse_token.go` | 新建 | Bearer Token 解析中间件 | `internal/middleware/parse_token.go` |
| `internal/validator/runtime.go` | 新建 | validator 初始化 + CheckPostParams/CheckQueryParams | `internal/validator/runtime.go` + `binding.go` |
| `internal/validator/form/auth.go` | 新建 | LoginForm 结构体 | `internal/validator/form/auth.go` |
| `internal/model/user.go` | 修改 | 新增 Password 字段 | `internal/model/admin_users.go` |
| `internal/routers/router.go` | 修改 | 接入 ParseToken 中间件 | `internal/routers/router.go` |
| `cmd/bootstrapx/bootstrap.go` | 修改 | 新增 InitializeValidator | `cmd/bootstrapx/bootstrap.go` |
| `cmd/root.go` | 修改 | PreRunE 调用 InitializeValidator | `cmd/root.go` |
| `config/config.go` | 修改 | JWTConfig 新增 TTL 字段 | `config/autoload/jwt.go` |
| `config/config.yaml` | 修改 | 新增 jwt.ttl | `config/config.yaml.example` |
| `go.mod` | 修改 | 新依赖 golang-jwt/jwt/v5, x/crypto | |

---

### Task 1: 修改 request 包 — 对齐 gin-layout

**文件:**
- 修改: `internal/pkg/request/request.go`

gin-layout 中 `GetAccessToken` 调用 `token.GetAccessToken` 来从 Authorization 头提取 Bearer Token，并返回 `(string, error)`。当前 GinAdmin 版本直接返回 string 且不区分错误原因。需要对齐。

gin-layout 没有 `BindJSON`/`BindQuery` 函数——参数校验统一走 `validator.CheckPostParams` / `validator.CheckQueryParams`（将在 Task 7 实现）。

- [x] **Step 1: 重写 `internal/pkg/request/request.go`**

参照 gin-layout `internal/pkg/request/request.go`：

```go
package request

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetQueryParams 提取当前请求的查询参数。
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

// GetAccessToken 从 Authorization 请求头提取 access token。
func GetAccessToken(c *gin.Context) (string, error) {
	if c == nil {
		return "", errors.New("gin context is nil")
	}
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		return "", errors.New("authorization header is missing")
	}
	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("invalid authorization header format")
	}
	return strings.TrimPrefix(authorization, "Bearer "), nil
}
```

注意：移除了 `response` 包的 import——gin-layout 的 request 包不依赖 response，只做纯粹的请求数据提取。

- [x] **Step 2: 验证编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

预期：user.go 仍然报错（缺少 errors 包引用），但 request 包本身无编译错误。

---

### Task 2: 新建错误码体系

**文件:**
- 新建: `internal/pkg/errors/code.go`
- 新建: `internal/pkg/errors/error.go`
- 新建: `internal/pkg/errors/zh-cn.go`

参照 gin-layout `internal/pkg/errors/code.go` 和 `internal/pkg/errors/error.go` 创建简化版（去掉 i18n 国际化，只用中文消息）。

- [x] **Step 1: 创建 `internal/pkg/errors/code.go`**

```go
package errors

const (
	SUCCESS          = 0
	NotLogin         = 401
	AuthorizationErr = 403
	NotFound         = 404
	ServerErr        = 500
	InvalidParameter = 10000
	UserDoesNotExist = 10001
	UserDisable      = 10002
	// 业务错误码 20000+
	UserPasswordWrong = 20001
	LoginFailed       = 20002
	TokenGenerateFailed = 20003
)
```

- [x] **Step 2: 创建 `internal/pkg/errors/error.go`**

参照 gin-layout 的 BusinessError 实现（去掉 i18n messageKey 相关逻辑，简化）。

```go
package errors

import (
	stderrors "errors"
	"fmt"
)

// BusinessError 带业务码的可控错误
type BusinessError struct {
	code    int
	message string
}

// Error 实现 error 接口
func (e *BusinessError) Error() string {
	return fmt.Sprintf("[Code]:%d [Msg]:%s", e.code, e.message)
}

// GetCode 返回业务错误码
func (e *BusinessError) GetCode() int {
	return e.code
}

// GetMessage 返回业务错误消息
func (e *BusinessError) GetMessage() string {
	return e.message
}

// NewBusinessError 创建业务错误
func NewBusinessError(code int, message ...string) *BusinessError {
	msg := ""
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	} else {
		msg = GetErrorMessage(code)
	}
	return &BusinessError{
		code:    code,
		message: msg,
	}
}

// Error 提供错误转换辅助方法
type ErrorHelper struct{}

// AsBusinessError 尝试把任意 error 转为 *BusinessError
func (e *ErrorHelper) AsBusinessError(err error) (*BusinessError, error) {
	var be *BusinessError
	if stderrors.As(err, &be) {
		return be, nil
	}
	return nil, err
}
```

- [x] **Step 3: 创建 `internal/pkg/errors/zh-cn.go`**

```go
package errors

var zhCNText = map[int]string{
	SUCCESS:             "OK",
	NotLogin:            "请先登录",
	AuthorizationErr:    "暂无权限",
	NotFound:            "资源不存在",
	ServerErr:           "服务器内部错误",
	InvalidParameter:    "参数错误",
	UserDoesNotExist:    "用户不存在",
	UserDisable:         "用户已被禁用",
	UserPasswordWrong:   "用户密码错误",
	LoginFailed:         "登录失败",
	TokenGenerateFailed: "Token 生成失败",
}

// GetErrorMessage 根据错误码获取中文错误消息
func GetErrorMessage(code int) string {
	if msg, ok := zhCNText[code]; ok {
		return msg
	}
	return "未知错误"
}
```

- [x] **Step 4: 验证编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

---

### Task 3: 修改 response 包 — 激活 Success/WithDataSuccess/Err

**文件:**

- 修改: `internal/pkg/response/response.go`

参照 gin-layout `internal/pkg/response/response.go` 中的 `Success`、`WithDataSuccess`、`FailCode` 便捷方法和 `Err` 智能错误处理。

- [x] **Step 1: 修改 `internal/pkg/response/response.go`**

在文件末尾追加以下方法（替换掉现有的注释块），同时在 import 中加入 `"GinAdmin/internal/pkg/errors"`：

```go
package response

import (
	"net/http"
	"reflect"
	"time"

	"GinAdmin/global"
	"GinAdmin/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

// 保留原有的 Result、NewResult、Response、Resp、Fail、FailCode、emptyObject、
// SetCode、SetMessage、Withdata、isNilData、isObjectData、json、Ok 不变。

// Success 正确返回
func (r *Response) Success(c *gin.Context) {
	r.SetCode(errors.SUCCESS)
	r.json(c)
}

// WithDataSuccess 成功后需要返回值
func (r *Response) WithDataSuccess(c *gin.Context, data interface{}) {
	r.SetCode(errors.SUCCESS)
	r.WithData(data)
	r.json(c)
}

// Success 业务成功响应（便捷方法）
func Success(c *gin.Context, data ...any) {
	if len(data) > 0 && data[0] != nil {
		Resp().WithDataSuccess(c, data[0])
		return
	}
	Resp().Success(c)
}

// FailCode 业务失败响应（便捷方法）
func FailCode(c *gin.Context, code int, data ...any) {
	if len(data) > 0 && data[0] != nil {
		Resp().WithData(data[0]).FailCode(c, code)
		return
	}
	Resp().FailCode(c, code)
}

// Fail 业务失败响应（便捷方法，带自定义消息）
func Fail(c *gin.Context, code int, message string, data ...any) {
	if len(data) > 0 && data[0] != nil {
		Resp().WithData(data[0]).Fail(c, code, message)
		return
	}
	Resp().Fail(c, code, message)
}

// Err 统一错误处理 — 参照 gin-layout controller.Api.Err
// 判断错误类型是 BusinessError 则自动返回携带的 code 和 message
func Err(c *gin.Context, err error) {
	helper := errors.ErrorHelper{}
	businessError, parseErr := helper.AsBusinessError(err)
	if parseErr != nil {
		// 非 BusinessError，返回服务器错误
		FailCode(c, errors.ServerErr)
		return
	}
	// BusinessError，使用其 code 和 message
	Fail(c, businessError.GetCode(), businessError.GetMessage())
}
```

- [x] **Step 2: 删除原文件末尾注释块**

删除 `internal/pkg/response/response.go` 第 171-190 行的 `// Success 业务成功响应` 注释块（已经被上面的新代码替代）。

- [ ] **Step 3: 验证编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

---

### Task 4: 修改 User Model — 新增 Password 字段

**文件:**
- 修改: `internal/model/user.go`

- [ ] **Step 1: 修改 `internal/model/user.go`**

在现有 User 结构体中加入 `Password` 字段（`json:"-"` 防止序列化泄露）：

```go
package model

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"column:id"`
	Username  string    `json:"username" gorm:"column:username"`
	Password  string    `json:"-" gorm:"column:password"`
	Nickname  string    `json:"nickname" gorm:"column:nickname"`
	Email     string    `json:"email" gorm:"column:email"`
	Status    int       `json:"status" gorm:"column:status"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}
```

- [ ] **Step 2: 验证编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

---

### Task 5: JWT 实现

**文件:**
- 新建: `internal/pkg/utils/jwt.go`
- 修改: `config/config.go` (JWTConfig 新增 TTL)
- 修改: `config/config.yaml` (新增 jwt.ttl)
- 修改: `go.mod` (新依赖)

参照 gin-layout `internal/pkg/utils/token/jwt.go` 的实现（简化版，去掉 phone/email/superadmin 字段）。

- [ ] **Step 1: 添加依赖**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go get golang.org/x/crypto@latest && go get github.com/golang-jwt/jwt/v5@latest
```

- [ ] **Step 2: 修改 `config/config.go` — JWTConfig 新增 TTL 字段**

```go
type JWTConfig struct {
	SecretKey string `yaml:"secret_key"`
	TTL       string `yaml:"ttl"` // 新增
}
```

- [ ] **Step 3: 修改 `config/config.yaml`**

```yaml
timezone: "Asia/Shanghai"

jwt:
  secret_key: "your-256-bit-secret-key"
  ttl: "2h"

mysql:
  enable: true
  host: 127.0.0.1
  port: 3306
  database: GinAdminMysql
  username: root
  password: 123456

redis:
  enable: false
```

- [ ] **Step 4: 创建 `internal/pkg/utils/jwt.go`**

参照 gin-layout `token.go`，简化版 JWT（Claims 只保留 UserID + Username + 标准字段）：

```go
package utils

import (
	"time"

	"GinAdmin/config"
	"GinAdmin/internal/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CustomClaims 自定义 JWT Claims
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, username string) (string, error) {
	cfg := config.GetConfig()

	// 解析 TTL，默认 2 小时
	ttl, parseErr := time.ParseDuration(cfg.JWT.TTL)
	if parseErr != nil {
		ttl = 2 * time.Hour
	}

	now := time.Now().UTC()
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "GinAdmin",
			Subject:   "admin",
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(cfg.JWT.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// ParseToken 解析 JWT Token
func ParseToken(accessToken string) (*CustomClaims, error) {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(accessToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.JWT.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.NewBusinessError(errors.NotLogin)
}
```

- [x] **Step 5: 验证编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

---

### Task 6: 新建 Service 层 — AuthService 登录业务

**文件:**
- 新建: `internal/service/auth.go`
- 修改: `internal/service/user.go` (新增 GetUserByID)

参照 gin-layout `internal/service/auth/login.go` 的 `validateUser` 逻辑（简化版，去掉 login log / security 等）。

- [x] **Step 1: 新建 `internal/service/auth.go`**

```go
package service

import (
	"GinAdmin/data"
	"GinAdmin/internal/model"
	"GinAdmin/internal/pkg/errors"
	"GinAdmin/internal/pkg/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// Login 用户登录，成功返回 JWT Token
func (s *AuthService) Login(username, password string) (string, *model.User, error) {
	// 1. 查询用户（含 password 字段）
	var user model.User
	err := data.GetDB().Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil, errors.NewBusinessError(errors.UserDoesNotExist)
		}
		return "", nil, errors.NewBusinessError(errors.ServerErr)
	}

	// 2. 检查用户状态（Status == 0 已禁用）
	if user.Status == 0 {
		return "", nil, errors.NewBusinessError(errors.UserDisable)
	}

	// 3. bcrypt 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.NewBusinessError(errors.UserPasswordWrong)
	}

	// 4. 生成 JWT
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, errors.NewBusinessError(errors.TokenGenerateFailed)
	}

	return token, &user, nil
}
```

- [ ] **Step 2: 修改 `internal/service/user.go` — 新增 GetUserByID**

在原有文件末尾追加：

```go
// GetUserByID 按 ID 查询用户信息（含密码字段，内部用）
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := data.GetDB().
		Select("id", "username", "nickname", "email", "status", "created_at", "updated_at").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
```

- [ ] **Step 3: 验证编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

---

### Task 7: 新建 Validator 包

**文件:**

- 新建: `internal/validator/form/common.go`（Paginate 分页表单 — 对应 gin-layout 同名文件）
- 新建: `internal/validator/form/auth.go`（LoginForm — 对应 gin-layout 同名文件）
- 新建: `internal/validator/runtime.go`（InitValidator + CheckPostParams/CheckQueryParams — 对应 gin-layout `runtime.go` + `binding.go`）

- [x] **Step 1: 新建 `internal/validator/form/common.go`**

参照 gin-layout `internal/validator/form/common.go`，定义通用表单结构体：

```go
package form

// Paginate 分页查询参数
type Paginate struct {
	Page    int `form:"page" json:"page" binding:"omitempty,gt=0"`
	PerPage int `form:"per_page" json:"per_page" binding:"omitempty,gt=0"`
}

// NewPaginate 创建一个新的分页查询
func NewPaginate() *Paginate {
	return &Paginate{}
}
```

- [ ] **Step 2: 新建 `internal/validator/form/auth.go`**

参照 gin-layout `internal/validator/form/auth.go`，去掉 Captcha 字段：

```go
package form

// LoginForm 登录表单
type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required,min=2,max=64"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=128"`
}
```

- [ ] **Step 3: 新建 `internal/validator/runtime.go`**

参照 gin-layout `internal/validator/runtime.go` + `binding.go`（简化版，不做多语言翻译）：

```go
package validator

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"GinAdmin/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// InitValidator 初始化 Gin 的 validator 引擎
// 参照 gin-layout InitValidatorTrans 的 initOnce() → init()
func InitValidator() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 参照 gin-layout registerTagNameFunc: 注册 json tag 作为字段名（用于错误提示）
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
// 参照 gin-layout CheckPostParams
func CheckPostParams(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBind(obj); err != nil {
		handleBindError(c, err)
		return err
	}
	return nil
}

// CheckQueryParams 检查 GET 查询参数，校验失败自动响应错误
// 参照 gin-layout CheckQueryParams
func CheckQueryParams(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		handleBindError(c, err)
		return err
	}
	return nil
}

// handleBindError 参照 gin-layout ResponseError → handleValidationError / handleBindingError
func handleBindError(c *gin.Context, err error) {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		msgs := make([]string, 0, len(validationErrs))
		for _, fieldErr := range validationErrs {
			msgs = append(msgs, fieldErr.Field()+": "+fieldErr.Tag())
		}
		response.Resp().Fail(c, 10000, strings.Join(msgs, "; "))
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
```

- [x] **Step 4: 验证编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

---

### Task 8: 新建 Controller 基类

**文件:**
- 新建: `internal/controller/sys_base.go`

参照 gin-layout `internal/controller/sys_base.go` 的 Api 结构体方法（简化版）。

- [ ] **Step 1: 创建 `internal/controller/sys_base.go`**

```go
package controller

import (
	"GinAdmin/internal/pkg/errors"
	response "GinAdmin/internal/pkg/response"
	"GinAdmin/global"

	"github.com/gin-gonic/gin"
)

// Api 控制器基类
type Api struct {
}

// Success 业务成功响应
func (api Api) Success(c *gin.Context, data ...any) {
	if len(data) > 0 && data[0] != nil {
		response.Success(c, data[0])
		return
	}
	response.Success(c)
}

// FailCode 业务失败响应（按错误码返回）
func (api Api) FailCode(c *gin.Context, code int, data ...any) {
	response.FailCode(c, code, data...)
}

// Fail 业务失败响应（自定义消息）
func (api Api) Fail(c *gin.Context, code int, message string, data ...any) {
	response.Fail(c, code, message, data...)
}

// Err 统一错误处理 — 参照 gin-layout controller.Api.Err
func (api Api) Err(c *gin.Context, err error) {
	helper := errors.ErrorHelper{}
	businessError, parseErr := helper.AsBusinessError(err)
	if parseErr != nil {
		response.FailCode(c, errors.ServerErr)
		return
	}
	response.Fail(c, businessError.GetCode(), businessError.GetMessage())
}

// GetCurrentUserID 获取当前登录用户 ID
func (api Api) GetCurrentUserID(c *gin.Context) uint {
	return c.GetUint(global.ContextKeyUID)
}
```

- [ ] **Step 2: 验证编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

---

### Task 9: 重写 Login Controller 和 User Controller

**文件:**
- 修改: `internal/controller/login.go`
- 修改: `internal/controller/user.go`

参照 gin-layout `internal/controller/admin_v1/auth.go` 的 Login 方法（简化版）。

- [x] **Step 1: 重写 `internal/controller/login.go`**

参照 gin-layout Login 逻辑：bind 表单 → 查用户验密 → 生成 token → 返回。

```go
package controller

import (
	"GinAdmin/internal/pkg/response"
	"GinAdmin/internal/service"
	"GinAdmin/internal/validator"
	"GinAdmin/internal/validator/form"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	Api
	authService *service.AuthService
}

func NewLoginController() *LoginController {
	return &LoginController{
		authService: service.NewAuthService(),
	}
}

// Login 用户登录
func (ctl *LoginController) Login(c *gin.Context) {
	params := &form.LoginForm{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}

	token, user, err := ctl.authService.Login(params.Username, params.Password)
	if err != nil {
		ctl.Err(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"email":    user.Email,
		},
	})
}
```

- [ ] **Step 2: 重写 `internal/controller/user.go`**

List 方法参照 gin-layout 的模式（`form.NewPaginate()` + `validator.CheckQueryParams`），Info 方法获取当前登录用户信息。

```go
package controller

import (
	"GinAdmin/internal/pkg/errors"
	"GinAdmin/internal/service"
	"GinAdmin/internal/validator"
	"GinAdmin/internal/validator/form"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Api
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// List 用户列表
func (ctl *UserController) List(c *gin.Context) {
	// 参照 gin-layout 风格：用 validator/form 包中的结构体 + CheckQueryParams
	params := form.NewPaginate()
	if err := validator.CheckQueryParams(c, params); err != nil {
		return
	}

	users, err := ctl.userService.ListUsers()
	if err != nil {
		ctl.Fail(c, 500, "查询用户列表失败")
		return
	}

	ctl.Success(c, users)
}

// Info 获取当前登录用户信息
func (ctl *UserController) Info(c *gin.Context) {
	uid := ctl.GetCurrentUserID(c)
	if uid == 0 {
		ctl.FailCode(c, errors.NotLogin)
		return
	}

	user, err := ctl.userService.GetUserByID(uid)
	if err != nil {
		ctl.FailCode(c, errors.NotFound)
		return
	}

	ctl.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"email":      user.Email,
		"status":     user.Status,
		"created_at": user.CreatedAt,
	})
}
```

- [ ] **Step 3: 验证编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

---

### Task 10: 新建 ParseToken 中间件

**文件:**
- 新建: `internal/middleware/parse_token.go`

参照 gin-layout `internal/middleware/parse_token.go` 实现（简化版，直接在中间件内解析 JWT）。

- [ ] **Step 1: 创建 `internal/middleware/parse_token.go`**

```go
package middleware

import (
	req "GinAdmin/internal/pkg/request"
	"GinAdmin/internal/pkg/utils"
	"GinAdmin/global"

	"github.com/gin-gonic/gin"
)

// ParseToken 全局 Token 解析中间件
// - 从 Authorization 头提取 Bearer Token
// - Token 有效：将 UserID 和 Username 注入 context
// - Token 无效或缺失：静默继续（不阻断请求）
func ParseToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := req.GetAccessToken(c)
		if err != nil || accessToken == "" {
			c.Next()
			return
		}

		claims, err := utils.ParseToken(accessToken)
		if err != nil || claims == nil {
			c.Next()
			return
		}

		// Token 有效，设置用户信息到 context
		c.Set(global.ContextKeyUID, claims.UserID)
		c.Set(global.ContextKeyAuthPrincipal, claims)

		c.Next()
	}
}
```

- [ ] **Step 2: 验证编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

---

### Task 11: 路由编排 — 接入 ParseToken 中间件

**文件:**
- 修改: `internal/routers/router.go`

- [x] **Step 1: 修改 `internal/routers/router.go` — SetRouters 中加入 ParseToken**

将 `SetRouters()` 函数改为：

```go
func SetRouters() *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.RequestMeta(), middleware.ParseToken())

	RegisterRoutes(engine, AppRouteTree())
	return engine
}
```

同时修改 import，加入 middleware 包的引用（当前已有）并确认包含 ParseToken。

- [ ] **Step 2: 将 LoginController 嵌入 Api 基类**

修改 `AdminRouteTree` 中 loginCtrl 的 Login handler —— 由于 LoginController 现在嵌入 `Api` 基类，Login 方法签名是 `func (ctl *LoginController) Login(c *gin.Context)`，与之前一致无需改动。

确认路由注册代码无需修改：

```go
func AdminRouteTree() RouteGroupDef {
	loginCtrl := controller.NewLoginController()
	userCtrl  := controller.NewUserController()

	return RouteGroupDef{
		Prefix: "api/v1",
		Children: []RouteGroupDef{
			{
				Routes: []RouteDef{
					POST("login", "登录", AuthNone, loginCtrl.Login),
				},
			},
			{
				Routes: []RouteDef{
					GET("user/info", "用户信息", AuthLogin, userCtrl.Info),
					GET("user/list", "用户列表", AuthPerm,  userCtrl.List),
				},
			},
		},
	}
}
```

- [ ] **Step 3: 验证编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

---

### Task 12: Bootstrap 接入 Validator 初始化

**文件:**
- 修改: `cmd/bootstrapx/bootstrap.go`
- 修改: `cmd/root.go`

- [ ] **Step 1: 修改 `cmd/bootstrapx/bootstrap.go`**

在文件末尾追加：

```go
import "GinAdmin/internal/validator"

// InitializeValidator 初始化验证器
func InitializeValidator() error {
	return validator.InitValidator()
}
```

需要修改现有 import 块（当前只 import config 和 time），新增 validator：

```go
package bootstrapx

import (
	"GinAdmin/config"
	"GinAdmin/internal/validator"
	"time"
)
```

保留原有的 `InitializeConfig`、`InitializeTimezone`、`InitializeLogger` 不变。

- [ ] **Step 2: 修改 `cmd/root.go`**

在 `PersistentPreRunE` 中调用 `InitializeValidator()`。找到 bootstrap 调用处：

```go
PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
    if shouldSkipBootstrap(cmd) {
        return nil
    }
    configPath, _ := cmd.Flags().GetString("config")
    if err := bootstrapx.InitializeConfig(configPath); err != nil {
        return err
    }
    bootstrapx.InitializeTimezone()
    // 新增
    if err := bootstrapx.InitializeValidator(); err != nil {
        return err
    }
    // ...
}
```

（具体行数取决于当前 root.go 的实际代码，需要精确匹配。）

- [x] **Step 3: 验证编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

---

### Task 13: 修复 `cmd/service/service.go` 中的 fmt.Errorf Bug

**文件:**
- 修改: `cmd/service/service.go`

- [ ] **Step 1: 将 `fmt.Errorf` 改为 `fmt.Printf`**

找到 `waitForShutdown` 函数中的这段：

```go
case sig := <-sigChan:
    fmt.Errorf("received shutdown signal: %s", sig.String())
```

改为：

```go
case sig := <-sigChan:
    fmt.Printf("received shutdown signal: %s\n", sig.String())
```

`fmt.Errorf` 只创建 error 但不输出，`fmt.Printf` 才会打印到控制台。

- [ ] **Step 2: 验证编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

预期：编译通过，零报错。

---

### Task 14: 最终编译 + 集成验证

- [ ] **Step 1: 完整编译**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...
```

预期：编译全部通过。

- [ ] **Step 2: 完整编译 + 静默检查**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go vet ./...
```

预期：无 vet 警告。

- [ ] **Step 3: 提交**

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin
git add -A
git commit -m "feat: 实现登录认证 + JWT + 受保护用户信息接口

- 新增 internal/pkg/errors 错误码体系（BusinessError + 中文消息）
- 新增 internal/pkg/utils/jwt.go JWT 生成/解析
- 新增 internal/validator 参数校验包（Paginate + LoginForm + CheckPostParams/CheckQueryParams）
- 新增 internal/controller/sys_base.go 控制器基类（参照 gin-layout Api 结构体）
- 新增 internal/middleware/parse_token.go Token 解析中间件
- 新增 internal/service/auth.go 登录业务逻辑
- 修改 request 包 GetAccessToken 返回 error（对齐 gin-layout）
- 修改 response 包激活 Success/FailCode/Err 方法
- 修改 user model 新增 Password 字段
- 修改 user service 新增 GetUserByID
- 完善 login controller（validator CheckPostParams → authService.Login → JWT）
- 完善 user controller（List 使用 form.Paginate + Info 从 context 获取 uid）
- 路由接入 ParseToken 全局中间件
- Bootstrap 接入 InitializeValidator
- 修复 service.go fmt.Errorf bug"
```

---

## 验证步骤

启动服务后依次验证：

```bash
# 0. 启动 MySQL
docker compose -f config/docker/docker-compose.yaml up -d

# 1. 初始化数据库表和数据
mysql -h 127.0.0.1 -u root -p123456 GinAdminMysql < SQL/createTableDemoSql.sql

# 2. 启动服务
go run main.go server

# 3. 测试健康检查
curl localhost:8080/ping
# 预期: "pong"

# 4. 测试登录（正确密码 123456）
curl -X POST localhost:8080/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"123456"}'
# 预期: {"code":0,"msg":"OK","data":{"access_token":"eyJ...","token_type":"Bearer","user":{...}}}

# 5. 测试登录（错误密码）
curl -X POST localhost:8080/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"wrong"}'
# 预期: {"code":20001,"msg":"用户密码错误",...}

# 6. 获取用户信息（携带 Token）
curl localhost:8080/api/v1/user/info \
  -H 'Authorization: Bearer <上面返回的access_token>'
# 预期: {"code":0,"msg":"OK","data":{"id":1,"username":"admin",...}}

# 7. 获取用户信息（无 Token）
curl localhost:8080/api/v1/user/info
# 预期: {"code":401,"msg":"请先登录",...}

# 8. 获取用户信息（伪造 Token）
curl localhost:8080/api/v1/user/info \
  -H 'Authorization: Bearer invalid_token_here'
# 预期: {"code":401,"msg":"请先登录",...}
```
