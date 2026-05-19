package controller

import (
	"context"
	"time"

	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/redis"
	"centraliz-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	response.SuccessWithMsg(c, "正常", gin.H{
		"status": "healthy",
		"timestamp": time.Now().Unix(),
	})
}

// ReadyCheck 就绪检查
func ReadyCheck(c *gin.Context) {
	// 检查数据库连接
	dbStatus := "healthy"
	if err := db.DB.Exec("SELECT 1").Error; err != nil {
		dbStatus = "unhealthy"
	}

	// 检查Redis连接
	redisStatus := "healthy"
	if _, err := redis.RDB.Ping(context.Background()).Result(); err != nil {
		redisStatus = "unhealthy"
	}

	status := "healthy"
	if dbStatus == "unhealthy" || redisStatus == "unhealthy" {
		status = "unhealthy"
	}

	response.SuccessWithMsg(c, "就绪检查", gin.H{
		"status": status,
		"database": dbStatus,
		"redis": redisStatus,
		"timestamp": time.Now().Unix(),
	})
}