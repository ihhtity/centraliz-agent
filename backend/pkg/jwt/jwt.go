package jwt

import (
	"centraliz-backend/pkg/config"
	"centraliz-backend/pkg/redis"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID   uint32 `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

const (
	TokenBlacklistPrefix = "token:blacklist:"
)

func GenerateToken(userID uint32, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "centraliz",
			Subject:   "user_token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, jwt.ErrTokenExpired
		}

		isBlacklisted, _ := IsTokenBlacklisted(tokenString)
		if isBlacklisted {
			return nil, jwt.ErrTokenInvalidId
		}

		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	timeUntilExpire := claims.ExpiresAt.Time.Sub(time.Now())
	if timeUntilExpire > time.Duration(config.AppConfig.JWT.ExpireHours)*time.Hour/2 {
		return tokenString, nil
	}

	err = BlacklistToken(tokenString, timeUntilExpire)
	if err != nil {
		log.Printf("blacklist token error: %v", err)
	}

	return GenerateToken(claims.UserID, claims.Username, claims.Role)
}

func IsTokenExpired(tokenString string) bool {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return true
	}
	return claims.ExpiresAt.Time.Before(time.Now())
}

func BlacklistToken(tokenString string, expiration time.Duration) error {
	key := TokenBlacklistPrefix + tokenString
	return redis.Set(key, "1", expiration)
}

func IsTokenBlacklisted(tokenString string) (bool, error) {
	key := TokenBlacklistPrefix + tokenString
	return redis.Exists(key)
}
