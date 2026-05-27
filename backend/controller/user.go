package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/jwt"
	"centraliz-backend/pkg/response"
	"errors"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Type     string `json:"type" binding:"required,oneof=phone email account"` // 登录类型
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Role     string `json:"role" binding:"required,oneof=user merch"`
}

type RegisterRequest struct {
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Account         string `json:"account"`
	Password        string `json:"password"`
	Code            string `json:"code"`
	Role            string `json:"role" binding:"required,oneof=user merch"`
	ConfirmPassword string `json:"confirmPassword"`
}

type ProfileResponse struct {
	ID        uint32  `json:"id"`
	Name      string  `json:"name"`
	Account   string  `json:"account"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	CreatedAt string  `json:"createdAt"`
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
		"user":  gin.H{"id": userID, "username": username, "role": req.Role},
	})
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

// 修改密码
func UserResetPassword(c *gin.Context) {

}

// GetProfile 获取个人资料
func GetProfile(c *gin.Context) {
	userID, _ := c.Get("userId")

	var user model.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		response.Fail(c, 404, "用户不存在")
		return
	}
	profile := ProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Account:   user.Account,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	response.SuccessWithMsg(c, "获取成功", profile)
}

// UpdateProfile 更新个人资料
func UpdateProfile(c *gin.Context) {
	// TODO: 实现更新个人资料逻辑
	// 这里简化处理，实际项目中需要根据角色更新对应的表
	response.SuccessWithMsg(c, "更新成功", gin.H{"success": true})
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
		if err := validateCode(req.Phone, "", req.Code, 1); err != nil {
			return nil, err
		}
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
		if err := validateCode("", req.Email, req.Code, 1); err != nil {
			return nil, err
		}
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

// validateCode 验证验证码
func validateCode(phone, email, inputCode string, codeType int8) error {
	var code model.Code

	if phone == "" && email == "" {
		return errors.New("请提供手机号或邮箱")
	}

	query := db.DB
	if phone != "" {
		query = query.Where("phone = ?", phone)
	} else {
		query = query.Where("email = ?", email)
	}

	if err := query.Where("type = ?", codeType).
		Where("status = ?", 0).
		Where("expire_at > ?", time.Now()).
		Order("created_at DESC").
		First(&code).Error; err != nil {
		return errors.New("验证码无效或已过期")
	}

	// 验证验证码
	if code.Code != inputCode {
		return errors.New("验证码错误")
	}

	// 标记验证码为已使用
	code.Status = 1
	if err := db.DB.Save(&code).Error; err != nil {
		return errors.New("更新验证码状态失败")
	}

	return nil
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
		req.Password = generateRandomPassword(12)
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

	user := model.User{
		Name:     username,
		Account:  req.Account,
		Password: string(hashedPassword),
		Email:    &req.Email,
		Phone:    &req.Phone,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return nil, errors.New("创建用户失败")
	}

	return &user, nil
}

// generateRandomPassword 生成随机密码
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// 在init函数中设置随机种子
func init() {
	rand.Seed(time.Now().UnixNano())
}
