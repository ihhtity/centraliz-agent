package controller

import (
	"centraliz-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetDashboardStats 获取仪表板统计数据
func GetDashboardStats(c *gin.Context) {
	// TODO: 实现获取仪表板统计数据逻辑
	response.SuccessWithMsg(c, "获取成功", gin.H{
		"totalUsers":    150,
		"totalDevices":  100,
		"totalRooms":    10,
		"totalOrders":   500,
		"todayOrders":   25,
		"availableDevices": 60,
		"occupiedDevices": 35,
		"maintenanceDevices": 5,
	})
}

// GetAccountList 获取账户列表
func GetAccountList(c *gin.Context) {
	// TODO: 实现获取账户列表逻辑
	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list": []gin.H{
			{
				"id":        1,
				"username":  "admin",
				"role":      "admin",
				"status":    "active",
				"createdAt": "2026-05-18T16:56:47Z",
			},
			{
				"id":        2,
				"username":  "manager1",
				"role":      "manager",
				"status":    "active",
				"createdAt": "2026-05-18T16:56:47Z",
			},
		},
		"total": 2,
	})
}

// CreateAccount 创建账户
func CreateAccount(c *gin.Context) {
	// TODO: 实现创建账户逻辑
	response.SuccessWithMsg(c, "创建成功", gin.H{"id": 3, "username": "newuser"})
}