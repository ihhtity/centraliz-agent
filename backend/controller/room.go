package controller

import (
	"centraliz-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetRoomList 获取房间列表
func GetRoomList(c *gin.Context) {
	// TODO: 实现获取房间列表逻辑
	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list": []gin.H{
			{
				"id":       1,
				"name":     "A区储物间",
				"location": "一楼东侧",
				"capacity": 50,
				"used":     25,
				"status":   "active", // active: 启用, inactive: 停用
				"createdAt": "2026-05-18T16:56:47Z",
			},
			{
				"id":       2,
				"name":     "B区储物间",
				"location": "二楼西侧",
				"capacity": 30,
				"used":     15,
				"status":   "active",
				"createdAt": "2026-05-18T16:56:47Z",
			},
		},
		"total": 2,
	})
}

// GetRoomDetail 获取房间详情
func GetRoomDetail(c *gin.Context) {
	// TODO: 实现获取房间详情逻辑
	response.SuccessWithMsg(c, "获取成功", gin.H{
		"id":        1,
		"name":      "A区储物间",
		"location":  "一楼东侧",
		"capacity":  50,
		"used":      25,
		"status":    "active",
		"manager":   "张三",
		"phone":     "13800138001",
		"createdAt": "2026-05-18T16:56:47Z",
		"updatedAt": "2026-05-18T16:56:47Z",
	})
}

// CreateRoom 创建房间
func CreateRoom(c *gin.Context) {
	// TODO: 实现创建房间逻辑
	response.SuccessWithMsg(c, "创建成功", gin.H{"id": 3, "name": "C区储物间"})
}

// UpdateRoom 更新房间
func UpdateRoom(c *gin.Context) {
	// TODO: 实现更新房间逻辑
	response.SuccessWithMsg(c, "更新成功", gin.H{"success": true})
}

// DeleteRoom 删除房间
func DeleteRoom(c *gin.Context) {
	// TODO: 实现删除房间逻辑
	response.SuccessWithMsg(c, "删除成功", gin.H{"success": true})
}