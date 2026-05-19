package controller

import (
	"centraliz-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetOrderList 获取订单列表
func GetOrderList(c *gin.Context) {
	// TODO: 实现获取订单列表逻辑
	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list": []gin.H{
			{
				"id":          1,
				"orderId":     "ORD20260518001",
				"userId":      1,
				"userName":    "测试用户",
				"deviceId":    1,
				"deviceName":  "智能储物柜001",
				"roomId":      1,
				"roomName":    "A区储物间",
				"status":      "completed", // pending: 待处理, using: 使用中, completed: 已完成, cancelled: 已取消
				"startTime":   "2026-05-18T16:56:47Z",
				"endTime":     "2026-05-18T18:56:47Z",
				"duration":    120, // 分钟
				"amount":      10.0, // 金额
				"createdAt":   "2026-05-18T16:56:47Z",
			},
			{
				"id":          2,
				"orderId":     "ORD20260518002",
				"userId":      1,
				"userName":    "测试用户",
				"deviceId":    2,
				"deviceName":  "智能储物柜002",
				"roomId":      1,
				"roomName":    "A区储物间",
				"status":      "using",
				"startTime":   "2026-05-18T16:56:47Z",
				"endTime":     "",
				"duration":    0,
				"amount":      0.0,
				"createdAt":   "2026-05-18T16:56:47Z",
			},
		},
		"total": 2,
	})
}

// GetOrderDetail 获取订单详情
func GetOrderDetail(c *gin.Context) {
	// TODO: 实现获取订单详情逻辑
	response.SuccessWithMsg(c, "获取成功", gin.H{
		"id":           1,
		"orderId":      "ORD20260518001",
		"userId":       1,
		"userName":     "测试用户",
		"userPhone":    "13800138000",
		"deviceId":     1,
		"deviceName":   "智能储物柜001",
		"roomId":       1,
		"roomName":     "A区储物间",
		"roomLocation": "一楼东侧",
		"status":       "completed",
		"startTime":    "2026-05-18T16:56:47Z",
		"endTime":      "2026-05-18T18:56:47Z",
		"duration":     120,
		"amount":       10.0,
		"paymentMethod": "wechat",
		"paymentStatus": "paid",
		"createdAt":    "2026-05-18T16:56:47Z",
		"updatedAt":    "2026-05-18T18:56:47Z",
	})
}

// CreateOrder 创建订单
func CreateOrder(c *gin.Context) {
	// TODO: 实现创建订单逻辑
	response.SuccessWithMsg(c, "创建成功", gin.H{"orderId": "ORD20260518003"})
}

// UpdateOrder 更新订单
func UpdateOrder(c *gin.Context) {
	// TODO: 实现更新订单逻辑
	response.SuccessWithMsg(c, "更新成功", gin.H{"success": true})
}