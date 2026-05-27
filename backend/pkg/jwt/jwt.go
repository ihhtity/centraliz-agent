package jwt

import (
	"centraliz-backend/pkg/config"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireHours) * time.Hour)), // 根据配置的过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "centraliz",
			Subject:   "user_token",
		},
	}

	// log.Printf("生成token - 用户ID: %d, 过期时间: %s", userID, claims.ExpiresAt.Time.Format("2006-01-02 15:04:05"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	// log.Printf("开始解析token...")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// log.Printf("token签名方法: %v", token.Header["alg"])
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		// log.Printf("token解析错误: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// log.Printf("token验证成功 - 过期时间: %s, 当前时间: %s",
		// 	claims.ExpiresAt.Time.Format("2006-01-02 15:04:05"),
		// 	time.Now().Format("2006-01-02 15:04:05"))

		// 检查是否过期
		if claims.ExpiresAt.Time.Before(time.Now()) {
			log.Printf("token已过期")
			return nil, jwt.ErrTokenExpired
		}

		return claims, nil
	}

	log.Printf("token无效")
	return nil, jwt.ErrSignatureInvalid
}
