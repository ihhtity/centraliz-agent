package controller

import (
	"centraliz-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetDeviceList 获取设备列表
func GetDeviceList(c *gin.Context) {
	// TODO: 实现获取设备列表逻辑
	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list": []gin.H{
			{
				"id":          1,
				"name":        "智能储物柜001",
				"roomId":      1,
				"roomName":    "A区储物间",
				"status":      "available", // available: 空闲, occupied: 租用, maintenance: 维修
				"battery":     95,
				"lastActive":  "2026-05-18T16:56:47Z",
				"createdAt":   "2026-05-18T16:56:47Z",
			},
			{
				"id":          2,
				"name":        "智能储物柜002",
				"roomId":      1,
				"roomName":    "A区储物间",
				"status":      "occupied",
				"battery":     80,
				"lastActive":  "2026-05-18T16:56:47Z",
				"createdAt":   "2026-05-18T16:56:47Z",
			},
		},
		"total": 2,
	})
}

// GetDeviceDetail 获取设备详情
func GetDeviceDetail(c *gin.Context) {
	// TODO: 实现获取设备详情逻辑
	response.SuccessWithMsg(c, "获取成功", gin.H{
		"id":         1,
		"name":       "智能储物柜001",
		"roomId":     1,
		"roomName":   "A区储物间",
		"status":     "available",
		"battery":    95,
		"ipAddress":  "192.168.1.100",
		"macAddress": "AA:BB:CC:DD:EE:FF",
		"firmware":   "v1.2.3",
		"createdAt":  "2026-05-18T16:56:47Z",
		"updatedAt":  "2026-05-18T16:56:47Z",
	})
}

// CreateDevice 创建设备
func CreateDevice(c *gin.Context) {
	// TODO: 实现创建设备逻辑
	response.SuccessWithMsg(c, "创建成功", gin.H{"id": 3, "name": "智能储物柜003"})
}

// UpdateDevice 更新设备
func UpdateDevice(c *gin.Context) {
	// TODO: 实现更新设备逻辑
	response.SuccessWithMsg(c, "更新成功", gin.H{"success": true})
}

// DeleteDevice 删除设备
func DeleteDevice(c *gin.Context) {
	// TODO: 实现删除设备逻辑
	response.SuccessWithMsg(c, "删除成功", gin.H{"success": true})
}