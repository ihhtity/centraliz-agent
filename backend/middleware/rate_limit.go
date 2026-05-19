package middleware

import (
	"net/http"
	"sync"

	"centraliz-backend/pkg/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	limiter     *rate.Limiter
	limiterOnce sync.Once
)

// getLimiter 获取限流器，懒加载模式
func getLimiter() *rate.Limiter {
	limiterOnce.Do(func() {
		limiter = rate.NewLimiter(
			rate.Limit(config.AppConfig.RateLimit.RequestsPerSecond),
			config.AppConfig.RateLimit.Burst,
		)
	})
	return limiter
}

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !getLimiter().Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "请求过多",
				"message": "请稍后再试",
			})
			return
		}
		c.Next()
	}
}