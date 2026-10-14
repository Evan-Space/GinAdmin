# Token 刷新 — 实现文档

> 参照 gin-layout 的 token refresh 逻辑，简化掉 Redis 锁 + 黑名单 + 自动续期。
> 核心：新增 refresh 接口，用当前有效的旧 Token 换一个新 Token。

---

## 原理

Token 快过期时，前端拿当前还有效的旧 Token 调 refresh 接口，后端签发一个新的 Token 返回。老 Token 在过期前仍可使用，前端收到新 Token 后替换即可。

```
前端                                     后端
 │                                        │
 ├── POST /api/v1/auth/refresh-token ───→│
 │    Authorization: Bearer <旧Token>    │  解析旧 Token → 验证有效
 │                                        │  → GenerateToken(uid, username)
 │←── {"access_token":"eyJ...新Token"}───┤
 │                                        │
```

---

## 步骤 1：config 中新增 refresh_ttl

### 修改 `config/config.go` — JWTConfig

```go
type JWTConfig struct {
    SecretKey  string `yaml:"secret_key"`
    TTL        string `yaml:"ttl"`         // access_token 有效期
    RefreshTTL string `yaml:"refresh_ttl"` // refresh 模式下新 token 有效期（可选，比 ttl 长）
}
```

### 修改 `config/config.yaml` — jwt 段

```yaml
jwt:
  secret_key: "your-256-bit-secret-key"
  ttl: "2h"
  refresh_ttl: "168h"  # 刷新后 7 天内有效
```

---

## 步骤 2：修改 JWT 工具 — 新增函数

### 修改 `internal/pkg/utils/jwt.go`

在 `ParseToken` 函数后面追加：

```go
// RefreshToken 用旧 token 中的用户信息签发新 token
// 对应 gin-layout internal/pkg/utils/token/jwt.go Refresh()
func RefreshToken(claims *CustomClaims) (string, error) {
    return GenerateToken(claims.UserID, claims.Username)
}
```

---

## 步骤 3：新建 Controller

### 新建 `internal/controller/auth.go`

```go
package controller

import (
    "GinAdmin/internal/pkg/errors"
    req "GinAdmin/internal/pkg/request"
    "GinAdmin/internal/pkg/utils"

    "github.com/gin-gonic/gin"
)

// AuthController 认证相关控制器
type AuthController struct {
    Api
}

func NewAuthController() *AuthController {
    return &AuthController{}
}

// RefreshToken 刷新 Token
// 拿当前有效的旧 Token，换一个新 Token
func (ctl *AuthController) RefreshToken(c *gin.Context) {
    // 1. 提取当前 Token
    accessToken, err := req.GetAccessToken(c)
    if err != nil || accessToken == "" {
        ctl.FailCode(c, errors.NotLogin)
        return
    }

    // 2. 解析旧 Token（能解析成功说明还合法）
    claims, err := utils.ParseToken(accessToken)
    if err != nil {
        ctl.FailCode(c, errors.NotLogin)
        return
    }

    // 3. 签发新 Token（用户信息不变，有效期重新算）
    newToken, err := utils.RefreshToken(claims)
    if err != nil {
        ctl.Err(c, err)
        return
    }

    ctl.Success(c, gin.H{
        "access_token": newToken,
        "token_type":   "Bearer",
    })
}

// CheckToken 检查 Token 是否有效
func (ctl *AuthController) CheckToken(c *gin.Context) {
    accessToken, err := req.GetAccessToken(c)
    if err != nil || accessToken == "" {
        ctl.Success(c, gin.H{"valid": false})
        return
    }

    _, err = utils.ParseToken(accessToken)
    ctl.Success(c, gin.H{"valid": err == nil})
}
```

---

## 步骤 4：注册路由

### 修改 `internal/routers/router.go`

在 `AdminRouteTree()` 中追加（放在公开接口下面）：

```go
{
    Routes: []RouteDef{
        POST("auth/refresh-token", "刷新Token", AuthLogin, authCtrl.RefreshToken),
        GET("auth/check-token", "检查Token", AuthLogin, authCtrl.CheckToken),
    },
},
```

并在顶部创建控制器变量：

```go
authCtrl := controller.NewAuthController()
```

---

## 新建/修改文件汇总

| 文件 | 操作 | 内容 |
|------|------|------|
| `config/config.go` | 修改 | JWTConfig 新增 `RefreshTTL` 字段（可选） |
| `config/config.yaml` | 修改 | 新增 `refresh_ttl` 配置项 |
| `internal/pkg/utils/jwt.go` | 修改 | 新增 `RefreshToken()` 函数 |
| `internal/controller/auth.go` | **新建** | refresh + check 两个方法 |
| `internal/routers/router.go` | 修改 | 追加两个路由 |

---

## 前端使用方式

```js
// 1. 登录成功后存 token
localStorage.setItem('token', response.data.access_token)

// 2. 每次请求前检查 token 是否快过期
function isTokenExpiring(token) {
    const payload = JSON.parse(atob(token.split('.')[1]))
    const exp = payload.exp * 1000           // JWT exp 是秒
    return Date.now() > exp - 5 * 60 * 1000  // 离过期不足 5 分钟
}

// 3. 快过期时刷新
async function refreshIfNeeded() {
    if (isTokenExpiring(localStorage.getItem('token'))) {
        const res = await fetch('/api/v1/auth/refresh-token', {
            headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
        })
        const data = await res.json()
        localStorage.setItem('token', data.data.access_token)
    }
}
```

---

## 测试验证

```bash
# 1. 登录
TOKEN=$(curl -s -X POST http://127.0.0.1:8080/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"super_admin","password":"123456"}' | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

# 2. 刷新
curl -X POST http://127.0.0.1:8080/api/v1/auth/refresh-token \
  -H "Authorization: Bearer $TOKEN"
# 预期: {"code":0,"data":{"access_token":"eyJ...","token_type":"Bearer"}}

# 3. 检查 Token 是否有效
curl http://127.0.0.1:8080/api/v1/auth/check-token \
  -H "Authorization: Bearer $TOKEN"
# 预期: {"code":0,"data":{"valid":true}}
```
