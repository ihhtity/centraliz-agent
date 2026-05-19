package middleware

import (
	"strings"
	"centraliz-backend/pkg/jwt"
	"centraliz-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, 401, "未授权")
			c.Abort()
			return
		}

		// 验证Bearer token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Fail(c, 401, "无效的token格式")
			c.Abort()
			return
		}

		tokenString := parts[1]
		// 验证JWT token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Fail(c, 401, "无效的token")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}