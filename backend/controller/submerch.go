package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SubMerchCreateRequest 创建子账号请求
type SubMerchCreateRequest struct {
	Account  string  `json:"account"`
	Password string  `json:"password"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Role     *string `json:"role"`   // 0商家 1管理者 2代理商
	Status   *string `json:"status"` // 0白名单 1黑名单
	Rule     string  `json:"rule"`   // 使用权限
}

// SubMerchUpdateRequest 更新子账号请求
type SubMerchUpdateRequest struct {
	Account  *string `json:"account"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Role     *string `json:"role"`
	Status   *string `json:"status"`
	Rule     *string `json:"rule"`
}

// SubMerchResponse 子账号响应
type SubMerchResponse struct {
	ID        int32   `json:"id"`
	Account   string  `json:"account"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Role      *string `json:"role"`
	Status    *string `json:"status"`
	Rule      string  `json:"rule"`
	LogAt     *string `json:"logAt"`
	CreatedAt *string `json:"createdAt"`
}

// CreateSubMerch 创建子账号
func CreateSubMerch(c *gin.Context) {
	var req SubMerchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 获取当前商家ID（从token中解析）
	merchID, exists := c.Get("merch_id")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 验证必填字段
	if req.Account == "" || req.Password == "" || req.Rule == "" {
		response.Fail(c, http.StatusBadRequest, "账号、密码和权限不能为空")
		return
	}

	// 检查账号是否已存在
	var count int64
	db.DB.Model(&model.SubMerch{}).Where("account = ?", req.Account).Count(&count)
	if count > 0 {
		response.Fail(c, http.StatusConflict, "账号已存在")
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	// 设置默认值
	role := "1" // 默认管理者
	if req.Role != nil && *req.Role != "" {
		role = *req.Role
	}

	status := "0" // 默认白名单
	if req.Status != nil && *req.Status != "" {
		status = *req.Status
	}

	// 创建子账号
	subMerch := &model.SubMerch{
		Account:  req.Account,
		Password: string(hashedPassword),
		MerchID:  merchID.(int32),
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     &role,
		Status:   &status,
		Rule:     req.Rule,
	}

	if err := db.DB.Create(subMerch).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建失败")
		return
	}

	response.SuccessWithMsg(c, "创建成功", toSubMerchResponse(subMerch))
}

// GetSubMerchList 获取子账号列表
func GetSubMerchList(c *gin.Context) {
	// 获取当前商家ID（从token中解析）
	merchsID := c.Query("merchs_id")
	if merchsID == "" {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 获取查询参数中的商家ID（可选）
	var merchID int32
	if merchsID != "" {
		// 如果传入了merchs_id参数，使用参数值
		id, err := strconv.ParseInt(merchsID, 10, 32)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "商家ID格式错误")
			return
		}
		merchID = int32(id)
	}

	var subMerches []*model.SubMerch
	if err := db.DB.Where("merch_id = ?", merchID).Order("created_at DESC").Find(&subMerches).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}

	result := make([]*SubMerchResponse, 0, len(subMerches))
	for _, sm := range subMerches {
		result = append(result, toSubMerchResponse(sm))
	}

	response.Success(c, result)
}

// GetSubMerchDetail 获取子账号详情
func GetSubMerchDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID格式错误")
		return
	}

	// 获取当前商家ID
	merchID, exists := c.Get("merch_id")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	var subMerch model.SubMerch
	if err := db.DB.Where("id = ? AND merch_id = ?", id, merchID.(int32)).First(&subMerch).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "子账号不存在")
		return
	}

	response.Success(c, toSubMerchResponse(&subMerch))
}

// UpdateSubMerch 更新子账号
func UpdateSubMerch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID格式错误")
		return
	}

	var req SubMerchUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 获取当前商家ID
	merchID, exists := c.Get("merch_id")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	var subMerch model.SubMerch
	if err := db.DB.Where("id = ? AND merch_id = ?", id, merchID.(int32)).First(&subMerch).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "子账号不存在")
		return
	}

	// 如果更新账号，检查是否已存在
	if req.Account != nil && *req.Account != "" && *req.Account != subMerch.Account {
		var count int64
		db.DB.Model(&model.SubMerch{}).Where("account = ? AND id != ?", *req.Account, id).Count(&count)
		if count > 0 {
			response.Fail(c, http.StatusConflict, "账号已存在")
			return
		}
		subMerch.Account = *req.Account
	}

	// 如果更新密码，重新加密
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "密码加密失败")
			return
		}
		subMerch.Password = string(hashedPassword)
	}

	if req.Email != nil {
		subMerch.Email = req.Email
	}
	if req.Phone != nil {
		subMerch.Phone = req.Phone
	}
	if req.Role != nil {
		subMerch.Role = req.Role
	}
	if req.Status != nil {
		subMerch.Status = req.Status
	}
	if req.Rule != nil {
		subMerch.Rule = *req.Rule
	}

	if err := db.DB.Save(&subMerch).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", toSubMerchResponse(&subMerch))
}

// DeleteSubMerch 删除子账号
func DeleteSubMerch(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID格式错误")
		return
	}

	// 获取当前商家ID
	merchID, exists := c.Get("merch_id")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	var subMerch model.SubMerch
	if err := db.DB.Where("id = ? AND merch_id = ?", id, merchID.(int32)).First(&subMerch).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "子账号不存在")
		return
	}

	if err := db.DB.Delete(&subMerch).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除失败")
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

// toSubMerchResponse 转换为响应结构
func toSubMerchResponse(sm *model.SubMerch) *SubMerchResponse {
	var logAt, createdAt *string
	if sm.LogAt != nil {
		logAtStr := sm.LogAt.Format("2006-01-02 15:04:05")
		logAt = &logAtStr
	}
	if sm.CreatedAt != nil {
		createdAtStr := sm.CreatedAt.Format("2006-01-02 15:04:05")
		createdAt = &createdAtStr
	}

	return &SubMerchResponse{
		ID:        sm.ID,
		Account:   sm.Account,
		Email:     sm.Email,
		Phone:     sm.Phone,
		Role:      sm.Role,
		Status:    sm.Status,
		Rule:      sm.Rule,
		LogAt:     logAt,
		CreatedAt: createdAt,
	}
}
