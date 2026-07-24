package middleware

import (
	"centraliz-backend/pkg/jwt"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

const (
	UserIDKey   = "user_id"
	UsernameKey = "username"
	RoleKey     = "role"
	TokenKey    = "token"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未授权",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "无效的token格式",
			})
			return
		}

		tokenString := parts[1]

		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			errorMsg := "无效的token"
			if errors.Is(err, jwtv4.ErrTokenExpired) {
				errorMsg = "token已过期"
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  errorMsg,
			})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UsernameKey, claims.Username)
		c.Set(RoleKey, claims.Role)
		c.Set(TokenKey, tokenString)

		c.Next()
	}
}

func GetUserID(c *gin.Context) uint32 {
	userID, ok := c.Get(UserIDKey)
	if !ok {
		return 0
	}
	return userID.(uint32)
}

func GetUsername(c *gin.Context) string {
	username, ok := c.Get(UsernameKey)
	if !ok {
		return ""
	}
	return username.(string)
}

func GetRole(c *gin.Context) string {
	role, ok := c.Get(RoleKey)
	if !ok {
		return ""
	}
	return role.(string)
}

func GetToken(c *gin.Context) string {
	token, ok := c.Get(TokenKey)
	if !ok {
		return ""
	}
	return token.(string)
}
