# Round 1: Login Auth + Protected API Design

## Goal

Implement a complete request flow: `POST /api/v1/login` → get JWT → `GET /api/v1/user/info` with Bearer token.

## Non-goals

- Captcha, token refresh, logout blacklist
- Casbin RBAC permission checking
- Redis, config hot-reload
- Logger initialization (deferred)

---

## 1. Fix Compilation Errors

### 1.1 New package: `internal/dto/`

Create `internal/dto/common.go` with `PageReq` struct:

```go
type PageReq struct {
    Page     int `json:"page" form:"page"`
    PageSize int `json:"page_size" form:"page_size"`
}

func (p *PageReq) Normalize() {
    if p.Page <= 0    { p.Page = 1 }
    if p.PageSize <= 0 { p.PageSize = 10 }
}
```

### 1.2 Add functions to `internal/pkg/request/request.go`

Add `BindQuery` and `BindJSON` helpers wrapping `c.ShouldBindQuery` / `c.ShouldBindJSON` with unified error responses.

### 1.3 Fix `internal/controller/user.go`

Replace the broken imports and references so it compiles. `List` method uses local `dto.PageReq` and `request.BindQuery`.

---

## 2. Error Code System

### New package: `internal/pkg/errors/`

```
errors/
├── code.go       // error code constants
├── error.go      // BusinessError type
└── message.go    // error code → human-readable message map
```

**Error codes:**

| Code | Constant | Description |
|------|----------|-------------|
| 0    | SUCCESS  | Success |
| 10001 | UNAUTHORIZED | Not logged in / token missing or invalid |
| 10002 | TOKEN_EXPIRED | Token has expired |
| 10003 | PARAM_ERROR | Request parameter validation failed |
| 10004 | NOT_FOUND | Resource not found |
| 10005 | INTERNAL_ERROR | Internal server error |
| 10006 | LOGIN_FAILED | Bad credentials |
| 10007 | ACCOUNT_DISABLED | Account is disabled |

`BusinessError` wraps code + message, implements `error` interface. `Resp().Err(c, err)` in response.go unboxes it.

### Update `internal/pkg/response/response.go`

- Uncomment `Success` / `WithDataSuccess` methods, using `errors.SUCCESS` constant.
- Wire up `Err(c, err)` for smart BusinessError handling.

---

## 3. JWT Implementation

### New file: `internal/pkg/utils/jwt.go`

- HS256 symmetric signing via `golang-jwt/jwt/v5`
- Custom claims: `UserID` (uint), `Username` (string), plus `jwt.RegisteredClaims` (ExpiresAt, IssuedAt, Issuer)
- API surface: `GenerateToken(userID uint, username string) (string, error)` and `ParseToken(tokenString string) (*Claims, error)`
- Token TTL read from config; default 2 hours if not configured

### Config changes

Add `ttl` field to `JWTConfig` struct and `config.yaml`:

```yaml
jwt:
  secret_key: "your-256-bit-secret-key"
  ttl: 2h
```

---

## 4. Login Form Validation

### New package: `internal/validator/`

```
validator/
├── runtime.go    // InitValidator(): register gin validator, set up translation
└── form/
    └── auth.go  // LoginForm struct
```

**LoginForm:**

```go
type LoginForm struct {
    Username string `json:"username" binding:"required,min=2,max=64"`
    Password string `json:"password" binding:"required,min=6,max=128"`
}
```

`runtime.go` wires up `go-playground/validator/v10` with Gin's binding. On validation failure, returns structured error with field-level messages.

### Call InitValidator in bootstrap

Add `InitializeValidator()` to `cmd/bootstrapx/bootstrap.go` and call it from `cmd/root.go`'s `PersistentPreRunE`.

---

## 5. Login Business Logic

### New file: `internal/service/auth.go`

```go
type AuthService struct{}

func (s *AuthService) Login(username, password string) (token string, user *model.User, err error)
```

Logic:
1. Query `users` table by username (include `password` column)
2. If not found → `BusinessError{Code: LOGIN_FAILED}`
3. Check `user.Status == 0` → `BusinessError{Code: ACCOUNT_DISABLED}`
4. `bcrypt.CompareHashAndPassword` → `BusinessError{Code: LOGIN_FAILED}`
5. `jwt.GenerateToken(user.ID, user.Username)` → return token

**Dependency:** `go get golang.org/x/crypto/bcrypt`

### Update `internal/model/user.go`

Add `Password` field with `json:"-" gorm:"column:password"` (hidden from JSON, readable by GORM).

### Rewrite `internal/controller/login.go`

`Login` method: bind JSON body → call `AuthService.Login` → return `{token, user_info}` on success.

---

## 6. Token Parse Middleware

### New file: `internal/middleware/parse_token.go`

```go
func ParseToken() gin.HandlerFunc
```

Logic:
1. Extract Bearer token from `Authorization` header
2. `jwt.ParseToken(token)` → if valid, set `uid` and `auth_principal` in Gin context
3. If token missing or invalid → do NOT block the request (just leave context empty)
4. Downstream handlers or middleware decide whether to reject unauthenticated requests

---

## 7. User Info Endpoint

### Update `internal/controller/user.go`

`Info` method:
1. Read `uid` from Gin context (set by ParseToken middleware)
2. If no uid → return `UNAUTHORIZED`
3. Query DB for user by id (columns: id, username, nickname, email, status, created_at)
4. Return user info via `response.Ok`

### New method in `internal/service/user.go`

```go
func (s *UserService) GetUserByID(id uint) (*model.User, error)
```

---

## 8. Route Wiring

### Update `internal/routers/router.go`

- Add `ParseToken` middleware to `SetRouters` engine: `engine.Use(middleware.RequestMeta(), middleware.ParseToken())`
- `POST /api/v1/login` stays `AuthNone`
- `GET /api/v1/user/info` stays `AuthLogin`
- `GET /api/v1/user/list` stays `AuthPerm` (compiles but no RBAC enforcement yet)

**Final middleware chain:** `RequestMeta → ParseToken → [route handlers]`

---

## 9. Model Update

### Update `internal/model/user.go`

Add `Password` field:
```go
Password  string    `json:"-" gorm:"column:password"`
```

---

## Summary of New / Modified Files

| Action | File |
|--------|------|
| **New** | `internal/dto/common.go` |
| **New** | `internal/pkg/errors/code.go` |
| **New** | `internal/pkg/errors/error.go` |
| **New** | `internal/pkg/errors/message.go` |
| **New** | `internal/pkg/utils/jwt.go` |
| **New** | `internal/validator/runtime.go` |
| **New** | `internal/validator/form/auth.go` |
| **New** | `internal/service/auth.go` |
| **New** | `internal/middleware/parse_token.go` |
| **Modify** | `internal/pkg/request/request.go` |
| **Modify** | `internal/pkg/response/response.go` |
| **Modify** | `internal/controller/user.go` |
| **Modify** | `internal/controller/login.go` |
| **Modify** | `internal/routers/router.go` |
| **Modify** | `internal/model/user.go` |
| **Modify** | `internal/service/user.go` |
| **Modify** | `cmd/bootstrapx/bootstrap.go` |
| **Modify** | `cmd/root.go` |
| **Modify** | `config/config.go` |
| **Modify** | `config/config.yaml` |
| **Modify** | `go.mod` (new dependency: `x/crypto`) |

---

## Verification

1. `go build ./...` compiles cleanly
2. Start MySQL via Docker: `docker compose -f config/docker/docker-compose.yaml up -d`
3. Create tables & seed data: run `SQL/createTableDemoSql.sql`
4. `go run main.go server` starts the API
5. `curl -X POST localhost:8080/api/v1/login -H 'Content-Type: application/json' -d '{"username":"admin","password":"123456"}'` returns a JWT token
6. `curl localhost:8080/api/v1/user/info -H 'Authorization: Bearer <token>'` returns user info
7. `curl localhost:8080/api/v1/user/info` (no token) returns 401 unauthorized
8. `curl localhost:8080/ping` returns "pong"
