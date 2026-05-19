package controller

import (
	"centraliz-backend/pkg/jwt"
	"centraliz-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// Login 用户登录
func Login(c *gin.Context) {
	// TODO: 实现真实的登录验证逻辑
	// 这里简化处理，实际项目中需要验证用户名和密码
	
	// 生成JWT令牌
	token, err := jwt.GenerateToken(1, "test", "user")
	if err != nil {
		response.Error(c, "生成token失败")
		return
	}
	
	response.SuccessWithMsg(c, "登录成功", gin.H{
		"token": token,
		"user":  gin.H{"id": 1, "username": "test", "role": "user"},
	})
}

// Register 用户注册
func Register(c *gin.Context) {
	// TODO: 实现注册逻辑
	response.SuccessWithMsg(c, "注册成功", gin.H{"id": 1, "username": "test"})
}

// GetProfile 获取用户个人资料
func GetProfile(c *gin.Context) {
	// TODO: 实现获取个人资料逻辑
	response.SuccessWithMsg(c, "获取成功", gin.H{
		"id":        1,
		"username":  "test",
		"email":     "test@example.com",
		"phone":     "13800138000",
		"realName":  "测试用户",
		"avatar":    "",
		"createdAt": "2026-05-18T16:56:47Z",
	})
}

// UpdateProfile 更新用户个人资料
func UpdateProfile(c *gin.Context) {
	// TODO: 实现更新个人资料逻辑
	response.SuccessWithMsg(c, "更新成功", gin.H{"success": true})
}