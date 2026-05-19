package middleware

import (
	"strings"
	"centraliz-backend/pkg/config"
	"github.com/gin-gonic/gin"
)

// CORS 跨域资源共享中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		corsConfig := config.AppConfig.CORS
		
		origin := c.GetHeader("Origin")
		if origin != "" {
			// 检查是否在允许的源列表中
			for _, allowedOrigin := range corsConfig.AllowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin || strings.HasSuffix(origin, allowedOrigin) {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		// 设置其他CORS头
		c.Header("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ","))
		c.Header("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ","))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}