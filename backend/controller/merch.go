package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/jwt"
	"centraliz-backend/pkg/mail"
	"centraliz-backend/pkg/redis"
	"centraliz-backend/pkg/response"
	"context"
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type MerchLoginRequest struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Type     string `json:"type"` // "phone", "email", "account"
}

type MerchRegisterRequest struct {
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Account         string `json:"account"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Code            string `json:"code"`
	Type            string `json:"type"` // "phone", "email", "account"
}

type MerchProfileResponse struct {
	ID        uint32  `json:"id"`
	Account   string  `json:"account"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	CreatedAt string  `json:"createdAt"`
}

// MerchLogin 商家登录: 可以手机号、邮箱或商家账号登录。手机号和邮箱需要验证码验证，商家账号需要验证账号密码
func MerchLogin(c *gin.Context) {
	var req MerchLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	// 验证登录方式
	if req.Type == "" {
		response.Fail(c, 400, "请提供登录方式")
		return
	}

	// 根据登录方式进行不同的验证
	var merch *model.Merch
	var err error

	if req.Type == "phone" || req.Type == "email" {
		// 验证码登录
		merch, err = validateCodeLogin(req)
	} else if req.Type == "account" {
		// 账号密码登录
		merch, err = validatePasswordLogin(req)
	} else {
		response.Fail(c, 400, "不支持的登录方式")
		return
	}

	if err != nil {
		response.Fail(c, 401, err.Error())
		return
	}

	// 更新登录时间
	now := time.Now()
	db.DB.Model(&merch).Update("log_at", now)

	// 生成JWT令牌
	token, err := jwt.GenerateToken(uint(merch.ID), merch.Account, "merch")
	if err != nil {
		response.Fail(c, 500, "生成token失败")
		return
	}

	response.SuccessWithMsg(c, "登录成功", gin.H{
		"token": token,
		"merch": merch,
	})
}

// MerchRegister 商家注册: 可以手机号、邮箱或商家账号注册。手机号和邮箱需要验证码验证，商家账号需要唯一性验证，密码需要加密存储
func MerchRegister(c *gin.Context) {
	var req MerchRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误", gin.H{"request": req})
		return
	}

	// 验证注册方式
	hasAccount := req.Account != ""
	hasPhone := req.Phone != ""
	hasEmail := req.Email != ""

	if !hasAccount && !hasPhone && !hasEmail {
		response.Fail(c, 400, "请至少提供商家账号、手机号或邮箱中的一个")
		return
	}

	// 账号注册
	if req.Type == "account" {
		// 验证确认密码
		if req.ConfirmPassword != req.Password {
			response.Fail(c, 400, "两次输入的密码不一致")
			return
		}

		// 验证密码强度
		if isValidPassword(req.Password) {
			response.Fail(c, 400, "密码长度必须在6-20位之间，且不能包含中文字符")
			return
		}

		// 验证输入的有效性
		if req.Account != "" && len(req.Account) < 4 {
			response.Fail(c, 400, "商家账号长度不能少于4位")
			return
		}
	}

	// 手机号或邮箱注册
	if req.Type == "phone" || req.Type == "email" {
		if req.Phone != "" && !isValidPhone(req.Phone) {
			response.Fail(c, 400, "手机号格式错误")
			return
		}

		if req.Email != "" && !isValidEmail(req.Email) {
			response.Fail(c, 400, "邮箱格式错误")
			return
		}

		// 如果提供了手机号或邮箱，需要验证码
		if (hasPhone || hasEmail) && req.Code == "" {
			response.Fail(c, 400, "手机号或邮箱注册需要验证码")
			return
		}

		// 验证验证码（如果提供了手机号或邮箱）
		if hasPhone && req.Code != "" {
			if !verifyCode("sms_code:"+req.Phone, req.Code) {
				response.Fail(c, 400, "手机号验证码错误或已过期")
				return
			}
		}

		if hasEmail && req.Code != "" {
			if !verifyCode("email_code:"+req.Email, req.Code) {
				response.Fail(c, 400, "邮箱验证码错误或已过期")
				return
			}
		}
	}

	merch, err := createMerch(req)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	var verifyKey string
	if req.Type == "phone" {
		verifyKey = "sms_code:" + req.Phone
	} else {
		verifyKey = "email_code:" + req.Email
	}

	if !verifyCode(verifyKey, req.Code) {
		response.Fail(c, 400, "验证码错误或已过期")
	}

	// 删除已使用的验证码
	deleteUsedCode(verifyKey)

	response.SuccessWithMsg(c, "注册成功", gin.H{"id": merch.ID, "type": req.Type})
}

// createMerch 创建商家
func createMerch(req MerchRegisterRequest) (*model.Merch, error) {
	// 检查商家账号是否已存在
	var existingMerch model.Merch

	// 如果提供了账号，检查账号是否已存在
	if req.Account != "" {
		if err := db.DB.Where("account = ?", req.Account).First(&existingMerch).Error; err == nil {
			return nil, errors.New("商家账号已存在")
		}
	}

	// 检查手机号是否已注册
	if req.Phone != "" {
		if err := db.DB.Where("phone = ?", req.Phone).First(&existingMerch).Error; err == nil {
			return nil, errors.New("手机号已注册为商家")
		}
	}

	// 检查邮箱是否已注册
	if req.Email != "" {
		if err := db.DB.Where("email = ?", req.Email).First(&existingMerch).Error; err == nil {
			return nil, errors.New("邮箱已注册为商家")
		}
	}

	// 如果没有提供账号，根据手机号或邮箱生成账号
	account := req.Account
	if account == "" {
		if req.Phone != "" {
			account = req.Phone
		} else if req.Email != "" {
			// 从邮箱提取用户名部分作为账号
			atIndex := strings.Index(req.Email, "@")
			if atIndex > 0 {
				account = req.Email[:atIndex]
			} else {
				account = req.Email
			}
		}
	}

	// 确保生成的账号唯一性
	if req.Account == "" {
		originalAccount := account
		counter := 1
		for {
			var checkMerch model.Merch
			if err := db.DB.Where("account = ?", account).First(&checkMerch).Error; err != nil {
				// 账号不存在，可以使用
				break
			}
			// 账号已存在，添加后缀
			account = originalAccount + strconv.Itoa(counter)
			counter++
			if counter > 100 { // 防止无限循环
				return nil, errors.New("无法生成唯一账号，请手动指定账号")
			}
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 根据模型定义，角色字段应该是"0"表示商家（参考模型注释）
	role := "0"
	status := "0" // 默认白名单状态

	merch := model.Merch{
		Account:  account,
		Password: string(hashedPassword),
		Email:    &req.Email,
		Phone:    &req.Phone,
		Role:     &role,
		Status:   &status,
	}

	if err := db.DB.Create(&merch).Error; err != nil {
		return nil, errors.New("创建商家失败")
	}

	// 创建默认分组
	defaultGroupName := "默认"
	defaultGroup := model.Group{
		Name:     &defaultGroupName,
		MerchsID: int32(merch.ID),
	}

	if err := db.DB.Create(&defaultGroup).Error; err != nil {
		// 如果创建默认分组失败，记录错误但不中断商家注册流程
		// 在实际生产环境中，建议使用事务来保证数据一致性
		return &merch, nil
	}

	return &merch, nil
}

// 修改密码
func MerchResetPassword(c *gin.Context) {
	type ResetPasswordRequest struct {
		ID              int    `json:"merchs_id"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误", req)
		return
	}

	// 验证新密码
	if req.NewPassword == "" {
		response.Fail(c, 400, "请提供新密码")
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		response.Fail(c, 400, "两次输入的密码不一致")
		return
	}

	// 验证密码强度
	if !isValidPassword(req.NewPassword) {
		response.Fail(c, 400, "密码长度必须在6-20位之间，且不能包含中文字符")
		return
	}

	// 查找商家是否存在
	var merch model.Merch
	var err error

	// 根据ID查找商家
	if err := db.DB.Where("id = ?", req.ID).First(&merch).Error; err != nil {
		response.Fail(c, 404, "商家不存在", req.ID)
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, "密码加密失败")
		return
	}

	// 更新密码
	if err := db.DB.Model(&merch).Update("password", string(hashedPassword)).Error; err != nil {
		response.Error(c, "密码更新失败")
		return
	}

	response.SuccessWithMsg(c, "密码重置成功", nil)
}

// GetMerchProfile 获取商家个人资料
func GetMerchProfile(c *gin.Context) {
	merchsID := c.Query("merchs_id")

	var merch model.Merch
	if err := db.DB.Where("id = ?", merchsID).First(&merch).Error; err != nil {
		response.Fail(c, 404, "商家不存在", merchsID)
		return
	}

	profile := MerchProfileResponse{
		ID:        merch.ID,
		Account:   merch.Account,
		Email:     merch.Email,
		Phone:     merch.Phone,
		CreatedAt: merch.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	response.SuccessWithMsg(c, "获取成功", profile)
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	type ChangePasswordRequest struct {
		MerchsID    string `json:"merchs_id" binding:"required"`
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	merchsID, err := strconv.Atoi(req.MerchsID)
	if err != nil {
		response.Fail(c, 400, "商家ID格式错误")
		return
	}

	var merch model.Merch
	if err := db.DB.Where("id = ?", merchsID).First(&merch).Error; err != nil {
		response.Fail(c, 404, "商家不存在", merchsID)
		return
	}

	// 验证原密码
	if err := bcrypt.CompareHashAndPassword([]byte(merch.Password), []byte(req.OldPassword)); err != nil {
		response.Fail(c, 400, "原密码错误")
		return
	}

	// 验证新密码强度
	if !isValidPassword(req.NewPassword) {
		response.Fail(c, 400, "密码长度必须在6-20位之间，且不能包含中文字符")
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, "密码加密失败")
		return
	}

	// 更新密码
	if err := db.DB.Model(&merch).Update("password", string(hashedPassword)).Error; err != nil {
		response.Error(c, "密码更新失败")
		return
	}

	response.SuccessWithMsg(c, "密码修改成功", nil)
}

// BindEmail 绑定/换绑邮箱
func BindEmail(c *gin.Context) {
	type BindEmailRequest struct {
		MerchsID string `json:"merchs_id"`
		Email    string `json:"email" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}

	var req BindEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	merchsID := req.MerchsID
	if merchsID == "" {
		response.Fail(c, 400, "商家ID不能为空")
		return
	}

	// 验证邮箱格式
	if !isValidEmail(req.Email) {
		response.Fail(c, 400, "邮箱格式错误")
		return
	}

	// 验证验证码
	verifyKey := "email_code:" + req.Email
	if !verifyCode(verifyKey, req.Code) {
		response.Fail(c, 400, "验证码错误或已过期")
		return
	}

	var merch model.Merch
	if err := db.DB.Where("id = ?", merchsID).First(&merch).Error; err != nil {
		response.Fail(c, 404, "商家不存在", merchsID)
		return
	}

	// 更新邮箱
	merch.Email = &req.Email
	if err := db.DB.Save(&merch).Error; err != nil {
		response.Error(c, "邮箱绑定失败")
		return
	}

	// 删除已使用的验证码
	deleteUsedCode(verifyKey)

	response.SuccessWithMsg(c, "邮箱绑定成功", gin.H{"email": req.Email})
}

// UnbindEmail 解绑邮箱
func UnbindEmail(c *gin.Context) {
	type UnbindRequest struct {
		MerchsID string `json:"merchs_id"`
	}

	var req UnbindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果 JSON 绑定失败，尝试从 URL 参数获取
		req.MerchsID = c.Query("merchs_id")
	}

	merchsID := req.MerchsID
	if merchsID == "" {
		response.Fail(c, 400, "商家ID不能为空")
		return
	}

	var merch model.Merch
	if err := db.DB.Where("id = ?", merchsID).First(&merch).Error; err != nil {
		response.Fail(c, 404, "商家不存在")
		return
	}

	if merch.Email == nil || *merch.Email == "" {
		response.Fail(c, 400, "未绑定邮箱")
		return
	}

	// 解绑邮箱（设置为空）
	emptyEmail := ""
	merch.Email = &emptyEmail
	if err := db.DB.Save(&merch).Error; err != nil {
		response.Error(c, "邮箱解绑失败")
		return
	}

	response.SuccessWithMsg(c, "邮箱解绑成功", nil)
}

// BindPhone 绑定/换绑手机号
func BindPhone(c *gin.Context) {
	type BindPhoneRequest struct {
		MerchsID string `json:"merchs_id"`
		Phone    string `json:"phone" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}

	var req BindPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	merchsID := req.MerchsID
	if merchsID == "" {
		response.Fail(c, 400, "商家ID不能为空")
		return
	}

	// 验证手机号格式
	if !isValidPhone(req.Phone) {
		response.Fail(c, 400, "手机号格式错误")
		return
	}

	// 验证验证码
	verifyKey := "sms_code:" + req.Phone
	if !verifyCode(verifyKey, req.Code) {
		response.Fail(c, 400, "验证码错误或已过期")
		return
	}

	var merch model.Merch
	if err := db.DB.Where("id = ?", merchsID).First(&merch).Error; err != nil {
		response.Fail(c, 404, "商家不存在", merchsID)
		return
	}

	merch.Phone = &req.Phone
	if err := db.DB.Save(&merch).Error; err != nil {
		response.Error(c, "手机号绑定失败")
		return
	}

	// 删除已使用的验证码
	deleteUsedCode(verifyKey)

	response.SuccessWithMsg(c, "手机号绑定成功", gin.H{"phone": req.Phone})
}

// UnbindPhone 解绑手机号
func UnbindPhone(c *gin.Context) {
	type UnbindRequest struct {
		MerchsID string `json:"merchs_id"`
	}

	var req UnbindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果 JSON 绑定失败，尝试从 URL 参数获取
		req.MerchsID = c.Query("merchs_id")
	}

	merchsID := req.MerchsID
	if merchsID == "" {
		response.Fail(c, 400, "商家ID不能为空")
		return
	}

	var merch model.Merch
	if err := db.DB.Where("id = ?", merchsID).First(&merch).Error; err != nil {
		response.Fail(c, 404, "商家不存在", merchsID)
		return
	}

	if merch.Phone == nil || *merch.Phone == "" {
		response.Fail(c, 400, "未绑定手机号")
		return
	}

	// 解绑手机号（设置为空）
	emptyPhone := ""
	merch.Phone = &emptyPhone
	if err := db.DB.Save(&merch).Error; err != nil {
		response.Error(c, "手机号解绑失败")
		return
	}

	response.SuccessWithMsg(c, "手机号解绑成功", nil)
}

// 发送验证码
func SendCode(c *gin.Context) {
	type SendCodeRequest struct {
		Phone   string `json:"phone"`
		Email   string `json:"email"`
		Type    int    `json:"type"`    // 1: 手机验证码, 2: 邮箱验证码
		Purpose string `json:"purpose"` // 发送验证码的用途，注册或重置密码等
	}

	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	// 设置默认用途为登录
	if req.Purpose == "" {
		req.Purpose = "login"
	}

	// 验证参数
	if req.Phone == "" && req.Email == "" {
		response.Fail(c, 400, "请提供手机号或邮箱")
		return
	}

	if req.Phone != "" && req.Email != "" {
		response.Fail(c, 400, "请只提供手机号或邮箱中的一个")
		return
	}

	if req.Phone != "" {
		// 验证手机号格式
		if !isValidPhone(req.Phone) {
			response.Fail(c, 400, "手机号格式错误")
			return
		}

		// 如果用途不是登录，检查手机号是否已注册为商家
		if req.Purpose != "login" {
			// 检查手机号是否已注册为商家
			var existingMerch model.Merch
			if err := db.DB.Where("phone = ?", req.Phone).First(&existingMerch).Error; err == nil {
				response.Fail(c, 400, "该手机号已注册")
				return
			}
		}
	}

	if req.Email != "" {
		// 验证邮箱格式
		if !isValidEmail(req.Email) {
			response.Fail(c, 400, "邮箱格式错误")
			return
		}

		// 如果用途不是登录，检查邮箱是否已注册为商家
		if req.Purpose != "login" {
			// 检查邮箱是否已注册为商家
			var existingMerch model.Merch
			if err := db.DB.Where("email = ?", req.Email).First(&existingMerch).Error; err == nil {
				response.Fail(c, 400, "该邮箱已注册")
				return
			}
		}
	}

	// 生成6位随机验证码
	rand.Seed(time.Now().UnixNano())
	code := strconv.Itoa(100000 + rand.Intn(900000))

	// 设置Redis键和过期时间（5分钟）
	var redisKey string
	expiration := 5 * time.Minute

	if req.Phone != "" {
		redisKey = "sms_code:" + req.Phone
		// 开发模式下，短信验证码默认返回123456
		code = "123456"
	} else if req.Email != "" {
		redisKey = "email_code:" + req.Email
	}

	// 将验证码存储到Redis中
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redis.RDB.Set(ctx, redisKey, code, expiration).Err(); err != nil {
		response.Error(c, "验证码存储失败")
		return
	}

	// 发送验证码
	if req.Email != "" {
		// 发送邮箱验证码
		mailSender := mail.GetMailSender()
		if mailSender == nil {
			response.Error(c, "邮件服务未初始化")
			return
		}

		subject := "验证码"
		body := code // 只传递验证码本身，邮件格式由mail.go处理

		// 创建带超时的上下文用于邮件发送
		mailCtx, mailCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer mailCancel()

		// 使用goroutine异步发送邮件，避免阻塞HTTP请求
		done := make(chan error, 1)
		go func() {
			err := mailSender.Send([]string{req.Email}, subject, body)
			done <- err
		}()

		select {
		case err := <-done:
			if err != nil {
				// 删除Redis中的验证码（发送失败）
				deleteUsedCode(redisKey)
				response.Error(c, "发送邮件失败，请检查邮箱地址")
				return
			}
		case <-mailCtx.Done():
			// 邮件发送超时，删除Redis中的验证码
			deleteUsedCode(redisKey)
			response.Error(c, "邮件发送超时")
			return
		}
	}

	response.SuccessWithMsg(c, "验证码已发送", gin.H{"success": true})
}

// validateCodeLogin 验证码登录验证
func validateCodeLogin(req MerchLoginRequest) (*model.Merch, error) {
	var merch model.Merch
	var err error

	// 根据登录方式查询商家
	if req.Type == "phone" {
		if req.Phone == "" {
			return nil, errors.New("请提供手机号")
		}
		if !isValidPhone(req.Phone) {
			return nil, errors.New("手机号格式错误")
		}
		err = db.DB.Where("phone = ?", req.Phone).First(&merch).Error
	} else if req.Type == "email" {
		if req.Email == "" {
			return nil, errors.New("请提供邮箱")
		}
		if !isValidEmail(req.Email) {
			return nil, errors.New("邮箱格式错误")
		}
		err = db.DB.Where("email = ?", req.Email).First(&merch).Error
	}

	if err != nil {
		return nil, errors.New("账号不存在")
	}

	// 验证验证码
	if req.Code == "" {
		return nil, errors.New("请提供验证码")
	}

	var verifyKey string
	if req.Type == "phone" {
		verifyKey = "sms_code:" + req.Phone
	} else {
		verifyKey = "email_code:" + req.Email
	}

	if !verifyCode(verifyKey, req.Code) {
		return nil, errors.New("验证码错误或已过期")
	}

	// 验证码验证成功后立即删除，防止二次使用
	deleteUsedCode(verifyKey)

	return &merch, nil
}

// validatePasswordLogin 账号密码登录验证
func validatePasswordLogin(req MerchLoginRequest) (*model.Merch, error) {
	var merch model.Merch

	if req.Account == "" {
		return nil, errors.New("请提供账号")
	}

	if req.Password == "" {
		return nil, errors.New("请提供密码")
	}

	// 根据账号查询商家
	if err := db.DB.Where("account = ?", req.Account).First(&merch).Error; err != nil {
		return nil, errors.New("商家账号不存在")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(merch.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	return &merch, nil
}

// isValidPhone 验证手机号格式
func isValidPhone(phone string) bool {
	// 简单的中国手机号验证
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

// isValidEmail 验证邮箱格式
func isValidEmail(email string) bool {
	matched, _ := regexp.MatchString(`^\S+@\S+\.\S+$`, email)
	return matched
}

// isValidPassword 验证密码强度
func isValidPassword(password string) bool {
	if len(password) < 6 || len(password) > 20 {
		return false
	}
	matched, _ := regexp.MatchString(`\p{Han}`, password)
	return !matched
}

// verifyCode 验证验证码
func verifyCode(key, code string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storedCode, err := redis.RDB.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	return storedCode == code
}

// deleteUsedCode 删除已使用的验证码
func deleteUsedCode(key string) {
	if key == "account" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	redis.RDB.Del(ctx, key)
}
