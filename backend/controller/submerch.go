package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"centraliz-backend/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SubMerchController struct{}

// CreateSubMerch 创建子账号
func (s *SubMerchController) Create(c *gin.Context) {
	var req struct {
		MerchsID int64  `json:"merchs_id"`
		Account  string `json:"account"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Role     string `json:"role"`
		Status   string `json:"status"`
		Rule     string `json:"rule"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 验证必填字段
	if req.Account == "" || req.Password == "" || req.Rule == "" {
		response.Fail(c, http.StatusBadRequest, "账号、密码和权限不能为空")
		return
	}

	// 检查账号是否已存在
	var count int64
	if err := db.DB.Model(&model.SubMerch{}).Where("merchs_id = ? AND account = ?", req.MerchsID, req.Account).Count(&count).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "检查账号失败")
		return
	}
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
	if req.Role == "" {
		req.Role = "1" // 默认管理者
	}
	if req.Status == "" {
		req.Status = "0" // 默认白名单
	}

	// 创建子账号
	subMerch := &model.SubMerch{
		Account:  req.Account,
		Password: string(hashedPassword),
		MerchsID: int32(req.MerchsID),
		Email:    utils.StringPtr(req.Email),
		Phone:    utils.StringPtr(req.Phone),
		Role:     utils.StringPtr(req.Role),
		Status:   utils.StringPtr(req.Status),
		Rule:     req.Rule,
	}

	if err := db.DB.Create(subMerch).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", toSubMerchResponse(subMerch))
}

// GetSubMerchList 获取子账号列表
func (s *SubMerchController) GetList(c *gin.Context) {
	merchsID := c.Query("merchs_id")
	if merchsID == "" {
		response.Fail(c, http.StatusBadRequest, "商家ID不能为空")
		return
	}

	var subMerches []*model.SubMerch
	if err := db.DB.Where("merchs_id = ?", merchsID).Order("created_at DESC").Find(&subMerches).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	result := make([]*SubMerchResponse, 0, len(subMerches))
	for _, sm := range subMerches {
		result = append(result, toSubMerchResponse(sm))
	}

	response.Success(c, result)
}

// GetSubMerchDetail 获取子账号详情
func (s *SubMerchController) GetDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "子账号ID不能为空")
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID格式错误")
		return
	}

	merchID, exists := c.Get("merch_id")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	var subMerch model.SubMerch
	if err := db.DB.Where("id = ? AND merch_id = ?", idInt, merchID.(int32)).First(&subMerch).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "子账号不存在")
		return
	}

	response.Success(c, toSubMerchResponse(&subMerch))
}

// UpdateSubMerch 更新子账号
func (s *SubMerchController) Update(c *gin.Context) {
	var req struct {
		ID       int32  `json:"id"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Role     string `json:"role"`
		Status   string `json:"status"`
		Rule     string `json:"rule"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if req.ID == 0 {
		response.Fail(c, http.StatusBadRequest, "子账号ID不能为空")
		return
	}

	var subMerch model.SubMerch
	if err := db.DB.Where("id = ?", req.ID).First(&subMerch).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "子账号不存在")
		return
	}

	if req.Email != "" {
		subMerch.Email = utils.StringPtr(req.Email)
	}
	if req.Phone != "" {
		subMerch.Phone = utils.StringPtr(req.Phone)
	}
	if req.Role != "" {
		subMerch.Role = utils.StringPtr(req.Role)
	}
	if req.Status != "" {
		subMerch.Status = utils.StringPtr(req.Status)
	}
	if req.Rule != "" {
		subMerch.Rule = req.Rule
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "密码加密失败")
			return
		}
		subMerch.Password = string(hashedPassword)
	}

	if err := db.DB.Save(&subMerch).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", toSubMerchResponse(&subMerch))
}

// DeleteSubMerch 删除子账号
func (s *SubMerchController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "子账号ID不能为空")
		return
	}

	if err := db.DB.Delete(&model.SubMerch{}, id).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除失败")
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
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
