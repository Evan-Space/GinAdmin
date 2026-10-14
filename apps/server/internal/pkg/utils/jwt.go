package utils

import (
	"GinAdmin/config"
	"GinAdmin/internal/pkg/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// / GenerateToken 生成 jwt token
func GenerateToken(userID uint, username string) (string, error) {
	cfg := config.GetConfig()

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

// ParseToken 解析 jwt token
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

// RefreshToken 用旧 token 中的信息签发新的 token
func RefreshToken(claims *CustomClaims) (string, error) {
	return GenerateToken(claims.UserID, claims.Username)
}
