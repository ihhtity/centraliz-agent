package middleware

import (
	"net/http"
	"sync"

	"centraliz-backend/pkg/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter 按IP的限流器
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

var (
	ipLimiter   *IPRateLimiter
	limiterOnce sync.Once
)

// NewIPRateLimiter 创建IP限流器
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// getLimiter 获取或创建指定IP的限流器
func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}
	return limiter
}

// Allow 检查是否允许请求
func (i *IPRateLimiter) Allow(ip string) bool {
	return i.getLimiter(ip).Allow()
}

// getIPLimiter 获取IP限流器实例，懒加载模式
func getIPLimiter() *IPRateLimiter {
	limiterOnce.Do(func() {
		ipLimiter = NewIPRateLimiter(
			rate.Limit(config.AppConfig.RateLimit.RequestsPerSecond),
			config.AppConfig.RateLimit.Burst,
		)
	})
	return ipLimiter
}

// RateLimit 限流中间件（按IP限流）
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()

		// 检查是否超过速率限制
		if !getIPLimiter().Allow(clientIP) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"msg":     "请求过多",
				"message": "请稍后再试",
			})
			return
		}
		c.Next()
	}
}
