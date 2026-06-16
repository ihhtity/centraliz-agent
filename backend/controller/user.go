package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/jwt"
	"centraliz-backend/pkg/response"
	"centraliz-backend/pkg/utils"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Account         string `json:"account"`
	Password        string `json:"password"`
	Code            string `json:"code"`
	Role            string `json:"role" binding:"required,oneof=user merch"`
	ConfirmPassword string `json:"confirmPassword"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Type     string `json:"type" binding:"required,oneof=phone email account"` // 登录类型
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Role     string `json:"role" binding:"required,oneof=user merch"`
}

// Register 用户注册
func UserRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	// 用户注册逻辑
	user, err := createUser(req)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.SuccessWithMsg(c, "注册成功", gin.H{"id": user.ID, "username": user.Account, "role": "user"})
}

// createUser 创建用户
func createUser(req RegisterRequest) (*model.User, error) {
	// 根据注册方式处理不同的逻辑
	if req.Account != "" && req.Password != "" {
		// 账号密码注册
		if req.ConfirmPassword == "" {
			return nil, errors.New("请确认密码")
		}
		if req.Password != req.ConfirmPassword {
			return nil, errors.New("两次输入的密码不一致")
		}
		// 验证密码强度（长度大于6位，不含中文字符）
		if len(req.Password) < 7 {
			return nil, errors.New("密码长度不能少于7位")
		}
		for _, r := range req.Password {
			if r >= 0x4e00 && r <= 0x9fa5 {
				return nil, errors.New("密码不能包含中文字符")
			}
		}
	} else if req.Phone != "" || req.Email != "" {
		// 手机号或邮箱注册，生成随机密码
		if req.Code == "" {
			return nil, errors.New("请输入验证码")
		}
		// TODO: 验证验证码的有效性（这里暂时跳过验证码验证）
		// 生成随机密码
		req.Password = utils.GenerateRandomPassword(12)
	} else {
		return nil, errors.New("请提供账号密码、手机号或邮箱进行注册")
	}

	// 检查账号是否已存在
	var existingUser model.User
	if req.Account != "" {
		if err := db.DB.Where("account = ?", req.Account).First(&existingUser).Error; err == nil {
			return nil, errors.New("账号已存在")
		}
	}
	if req.Phone != "" {
		if err := db.DB.Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
			return nil, errors.New("手机号已注册")
		}
	}
	if req.Email != "" {
		if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return nil, errors.New("邮箱已注册")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 设置用户名
	username := req.Account
	if username == "" {
		if req.Phone != "" {
			username = req.Phone
		} else if req.Email != "" {
			username = req.Email
		}
	}

	// 处理 Email 和 Phone 指针类型
	var emailPtr, phonePtr *string
	if req.Email != "" {
		emailPtr = &req.Email
	}
	if req.Phone != "" {
		phonePtr = &req.Phone
	}

	user := model.User{
		Name:     username,
		Account:  req.Account,
		Password: string(hashedPassword),
		Email:    emailPtr,
		Phone:    phonePtr,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return nil, errors.New("创建用户失败")
	}

	return &user, nil
}

// Login 用户登录
func UserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	var userID uint32
	var username string

	// 用户登录逻辑
	user, err := validateUserLogin(req)
	if err != nil {
		response.Fail(c, 401, err.Error())
		return
	}
	userID = user.ID
	username = user.Account

	// 生成JWT令牌
	token, err := jwt.GenerateToken(uint(userID), username, req.Role)
	if err != nil {
		response.Error(c, "生成token失败")
		return
	}

	response.SuccessWithMsg(c, "登录成功", gin.H{
		"token": token,
		"user":  user,
	})
}

// validateUserLogin 验证用户登录
func validateUserLogin(req LoginRequest) (*model.User, error) {
	var user model.User

	// 根据登录类型查询用户
	switch req.Type {
	case "phone":
		if req.Phone == "" {
			return nil, errors.New("请输入手机号")
		}
		if err := db.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
			return nil, errors.New("用户不存在")
		}
		// 验证码登录，验证验证码
		if req.Code == "" {
			return nil, errors.New("请输入验证码")
		}
		verifyKey := "sms_code:" + req.Phone
		if !utils.VerifyCode(verifyKey, req.Code) {
			return nil, errors.New("验证码错误或已过期")
		}
		// 删除已使用的验证码
		utils.DeleteUsedCode(verifyKey)
	case "email":
		if req.Email == "" {
			return nil, errors.New("请输入邮箱")
		}
		if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
			return nil, errors.New("用户不存在")
		}
		// 验证码登录，验证验证码
		if req.Code == "" {
			return nil, errors.New("请输入验证码")
		}
		verifyKey := "email_code:" + req.Email
		if !utils.VerifyCode(verifyKey, req.Code) {
			return nil, errors.New("验证码错误或已过期")
		}
		// 删除已使用的验证码
		utils.DeleteUsedCode(verifyKey)
	case "account":
		if req.Account == "" {
			return nil, errors.New("请输入账号")
		}
		if err := db.DB.Where("account = ?", req.Account).First(&user).Error; err != nil {
			return nil, errors.New("账号不存在")
		}
		// 密码登录，验证密码
		if req.Password == "" {
			return nil, errors.New("请输入密码")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return nil, errors.New("密码错误")
		}
	default:
		return nil, errors.New("不支持的登录方式")
	}

	// 检查用户状态
	if user.Status != nil && *user.Status == "1" {
		return nil, errors.New("账号已被禁用")
	}

	return &user, nil
}

// GetProfile 获取个人资料
func GetProfile(c *gin.Context) {
	type ProfileResponse struct {
		ID        uint32  `json:"id"`
		Name      string  `json:"name"`
		Account   string  `json:"account"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
		CreatedAt string  `json:"createdAt"`
	}

	userID := c.Param("id")

	var user model.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		response.Fail(c, 404, "用户不存在")
		return
	}

	// 处理时间字段
	createdAt := ""
	if user.CreatedAt != nil {
		createdAt = user.CreatedAt.Format("2006-01-02T15:04:05Z")
	}

	// 处理 Email 和 Phone 指针类型
	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	phone := ""
	if user.Phone != nil {
		phone = *user.Phone
	}

	profile := ProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Account:   user.Account,
		Email:     &email,
		Phone:     &phone,
		CreatedAt: createdAt,
	}
	response.SuccessWithMsg(c, "获取成功", profile)
}

// UpdateProfile 更新个人资料（支持修改昵称、密码、手机号绑定、解绑、换绑）
func UpdateProfile(c *gin.Context) {
	type UpdateProfileRequest struct {
		ID          uint32 `json:"id"`
		Name        string `json:"name"`
		Phone       string `json:"phone"`
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	// 优先使用请求体中的ID，否则从JWT中获取
	userID := req.ID
	if userID == 0 {
		response.Fail(c, 400, "用户ID不能为空")
		return
	}

	var user model.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		response.Fail(c, 404, "用户不存在")
		return
	}

	// 修改昵称
	if req.Name != "" {
		user.Name = strings.TrimSpace(req.Name)
	}

	// 修改手机号（绑定/换绑/解绑）
	if req.Phone != "" {
		// 检查手机号是否已被其他用户绑定
		var existingUser model.User
		if err := db.DB.Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
			if existingUser.ID != req.ID {
				response.Fail(c, 400, "该手机号已被其他账号绑定")
				return
			}
		}
		user.Phone = &req.Phone
	} else if req.Phone == "" {
		// 解绑手机号
		user.Phone = nil
	}

	// 修改密码
	if req.OldPassword != "" && req.NewPassword != "" {
		// 验证旧密码
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
			response.Fail(c, 400, "原密码错误")
			return
		}

		// 加密新密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			response.Error(c, "密码加密失败")
			return
		}
		user.Password = string(hashedPassword)
	}

	// 保存更新
	if err := db.DB.Save(&user).Error; err != nil {
		response.Error(c, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", gin.H{"name": user.Name})
}

// UserBindEmail 绑定/换绑邮箱
func UserBindEmail(c *gin.Context) {
	type BindEmailRequest struct {
		ID    uint32 `json:"id"`
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	var req BindEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	// 验证邮箱格式
	if !utils.IsValidEmail(req.Email) {
		response.Fail(c, 400, "邮箱格式错误")
		return
	}

	// 验证验证码
	verifyKey := "email_code:" + req.Email
	if !utils.VerifyCode(verifyKey, req.Code) {
		response.Fail(c, 400, "验证码错误或已过期")
		return
	}

	// 检查用户是否存在
	var user model.User
	if err := db.DB.Where("id = ?", req.ID).First(&user).Error; err != nil {
		response.Fail(c, 404, "用户不存在")
		return
	}

	// 检查邮箱是否已被其他用户绑定
	var existingUser model.User
	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		if existingUser.ID != req.ID {
			response.Fail(c, 400, "该邮箱已被其他账号绑定")
			return
		}
	}

	// 更新用户邮箱
	user.Email = &req.Email
	if err := db.DB.Save(&user).Error; err != nil {
		response.Error(c, "邮箱绑定失败")
		return
	}

	// 删除已使用的验证码
	utils.DeleteUsedCode(verifyKey)

	response.SuccessWithMsg(c, "邮箱绑定成功", gin.H{"email": req.Email})
}

// UserUnbindEmail 解绑邮箱
func UserUnbindEmail(c *gin.Context) {
	// UserUnbindEmailRequest 解绑邮箱请求
	type UserUnbindEmailRequest struct {
		ID uint32 `json:"id"`
	}

	var req UserUnbindEmailRequest
	_ = c.ShouldBindJSON(&req)

	// 优先使用请求体中的ID，否则从JWT中获取
	userID := req.ID
	if userID == 0 {
		response.Fail(c, 400, "用户ID不能为空")
		return
	}

	var user model.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		response.Fail(c, 404, "用户不存在")
		return
	}

	if user.Email == nil || *user.Email == "" {
		response.Fail(c, 400, "未绑定邮箱")
		return
	}

	// 解绑邮箱（设置为空）
	user.Email = nil
	if err := db.DB.Save(&user).Error; err != nil {
		response.Error(c, "邮箱解绑失败")
		return
	}

	response.SuccessWithMsg(c, "邮箱解绑成功", nil)
}
