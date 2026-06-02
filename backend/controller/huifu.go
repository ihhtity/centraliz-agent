package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"centraliz-backend/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type HuifuController struct{}

type HuifuCreateRequest struct {
	MerchsID  int32   `json:"merchs_id" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	Account   string  `json:"account" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Phone     string  `json:"phone" binding:"required"`
	Identity  string  `json:"identity" binding:"required"`
	Card      string  `json:"card" binding:"required"`
	Storename string  `json:"storename"`
	Encrypt   string  `json:"encrypt"`
	Area      string  `json:"area"`
	Remarks   string  `json:"remarks"`
	Share     string  `json:"share"`
	Rate      float64 `json:"rate"`
}

type HuifuUpdateRequest struct {
	Type      string  `json:"type" binding:"required"`
	Account   string  `json:"account" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Phone     string  `json:"phone" binding:"required"`
	Identity  string  `json:"identity" binding:"required"`
	Card      string  `json:"card" binding:"required"`
	Storename string  `json:"storename"`
	Encrypt   string  `json:"encrypt"`
	Area      string  `json:"area"`
	Remarks   string  `json:"remarks"`
	Share     string  `json:"share"`
	Rate      float64 `json:"rate"`
}

// GetList 获取收款账号列表
func (h *HuifuController) GetList(c *gin.Context) {
	merchsID := c.Query("merchs_id")

	if merchsID == "" {
		response.Fail(c, 400, "商家ID不能为空")
		return
	}

	var accounts []model.HuifuAccount
	if err := db.DB.Where("merchs_id = ?", merchsID).Order("id desc").Find(&accounts).Error; err != nil {
		response.Fail(c, 500, "获取收款账号列表失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":  accounts,
		"total": len(accounts),
	})
}

// GetDetail 获取收款账号详情
func (h *HuifuController) GetDetail(c *gin.Context) {
	// 获取账号ID
	id := c.Query("id")
	if id == "" {
		response.Fail(c, 400, "收款ID不能为空")
		return
	}

	var account model.HuifuAccount
	if err := db.DB.Where("id = ?", id).First(&account).Error; err != nil {
		response.Fail(c, 404, "收款账号不存在")
		return
	}

	response.Success(c, account)
}

// Create 创建收款账号
func (h *HuifuController) Create(c *gin.Context) {
	var req HuifuCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	// 验证必填字段
	if err := h.validateRequiredFields(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	// 验证手机号格式
	if err := h.validatePhone(req.Phone); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	// 验证身份证号格式
	if err := h.validateIdentity(req.Identity); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	// 验证银行卡号格式
	if err := h.validateCard(req.Card); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	// 验证企业账号必填字段
	if err := h.validateCompanyAccount(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	// 检查账号是否已存在
	if exists, err := h.isAccountExists(req.MerchsID, req.Account); err != nil {
		response.Fail(c, 500, "检查账号重复失败")
		return
	} else if exists {
		response.Fail(c, 400, "该账号已存在")
		return
	}

	now := time.Now()
	account := model.HuifuAccount{
		MerchsID:  req.MerchsID,
		Type:      &req.Type,
		Account:   &req.Account,
		Name:      &req.Name,
		Phone:     &req.Phone,
		Identity:  &req.Identity,
		Card:      &req.Card,
		Storename: utils.StringPtr(req.Storename),
		Encrypt:   utils.StringPtr(req.Encrypt),
		Area:      req.Area,
		Remarks:   req.Remarks,
		Share:     req.Share,
		Rate:      &req.Rate,
		Choose:    utils.StringPtr("0"),
		CreatedAt: &now,
	}

	if err := db.DB.Create(&account).Error; err != nil {
		response.Fail(c, 500, "创建收款账号失败")
		return
	}

	response.Success(c, account)
}

// isAccountExists 检查账号是否已存在
func (h *HuifuController) isAccountExists(merchsID int32, account string) (bool, error) {
	var count int64
	if err := db.DB.Model(&model.HuifuAccount{}).Where("merchs_id = ? AND account = ?", merchsID, account).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update 更新收款账号
func (h *HuifuController) Update(c *gin.Context) {
	// 获取账号ID
	id := c.Query("id")
	if id == "" {
		response.Fail(c, 400, "收款账号ID不能为空")
		return
	}

	var account model.HuifuAccount
	var req HuifuUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := h.validateCompanyAccountUpdate(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	account.Type = &req.Type
	account.Account = &req.Account
	account.Name = &req.Name
	account.Phone = &req.Phone
	account.Identity = &req.Identity
	account.Card = &req.Card
	account.Storename = utils.StringPtr(req.Storename)
	account.Encrypt = utils.StringPtr(req.Encrypt)
	account.Area = req.Area
	account.Remarks = req.Remarks
	account.Share = req.Share
	account.Rate = &req.Rate

	if err := db.DB.Save(account).Error; err != nil {
		response.Fail(c, 500, "更新收款账号失败")
		return
	}

	response.Success(c, account)
}

// Delete 删除收款账号
func (h *HuifuController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, 400, "收款账号ID不能为空")
		return
	}

	var account model.HuifuAccount
	if err := db.DB.Where("id = ?", id).First(&account).Error; err != nil {
		response.Fail(c, 404, "收款账号不存在")
		return
	}

	if err := db.DB.Delete(&account).Error; err != nil {
		response.Fail(c, 500, "删除收款账号失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", gin.H{"success": true})
}

// SetChoose 设置默认收款账号
func (h *HuifuController) SetChoose(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, 400, "收款账号ID不能为空")
		return
	}

	var account model.HuifuAccount
	if err := db.DB.Where("id = ?", id).First(&account).Error; err != nil {
		response.Fail(c, 404, "收款账号不存在")
		return
	}

	account.Choose = utils.StringPtr("1")
	if err := db.DB.Save(&account).Error; err != nil {
		response.Fail(c, 500, "设置默认收款账号失败")
		return
	}

	response.Success(c, account)
}

// validateCompanyAccount 验证企业账号
func (h *HuifuController) validateCompanyAccount(req *HuifuCreateRequest) error {
	if req.Type == "company" {
		if req.Storename == "" {
			return &ValidationError{Msg: "企业账号需要填写店名"}
		}
		if req.Encrypt == "" {
			return &ValidationError{Msg: "企业账号需要填写营业执照编码"}
		}
	}
	return nil
}

// validateCompanyAccountUpdate 验证企业账号更新
func (h *HuifuController) validateCompanyAccountUpdate(req *HuifuUpdateRequest) error {
	if req.Type == "company" {
		if req.Storename == "" {
			return &ValidationError{Msg: "企业账号需要填写店名"}
		}
		if req.Encrypt == "" {
			return &ValidationError{Msg: "企业账号需要填写营业执照编码"}
		}
	}
	return nil
}

// validateRequiredFields 验证必填字段
func (h *HuifuController) validateRequiredFields(req *HuifuCreateRequest) error {
	if req.Type == "" {
		return &ValidationError{Msg: "请选择账号类型"}
	}
	if req.Account == "" {
		return &ValidationError{Msg: "请输入账号"}
	}
	if req.Name == "" {
		return &ValidationError{Msg: "请输入姓名"}
	}
	if req.Phone == "" {
		return &ValidationError{Msg: "请输入手机号"}
	}
	if req.Identity == "" {
		return &ValidationError{Msg: "请输入身份证号"}
	}
	if req.Card == "" {
		return &ValidationError{Msg: "请输入银行卡号"}
	}
	return nil
}

// validatePhone 验证手机号格式
func (h *HuifuController) validatePhone(phone string) error {
	if len(phone) != 11 {
		return &ValidationError{Msg: "手机号长度必须为11位"}
	}
	for _, c := range phone {
		if c < '0' || c > '9' {
			return &ValidationError{Msg: "手机号只能包含数字"}
		}
	}
	if phone[0] != '1' {
		return &ValidationError{Msg: "手机号必须以1开头"}
	}
	return nil
}

// validateIdentity 验证身份证号格式
func (h *HuifuController) validateIdentity(identity string) error {
	length := len(identity)
	if length != 15 && length != 18 {
		return &ValidationError{Msg: "身份证号长度必须为15位或18位"}
	}
	if length == 18 {
		for i := 0; i < 17; i++ {
			if identity[i] < '0' || identity[i] > '9' {
				return &ValidationError{Msg: "身份证号前17位必须为数字"}
			}
		}
		lastChar := identity[17]
		if !((lastChar >= '0' && lastChar <= '9') || lastChar == 'X' || lastChar == 'x') {
			return &ValidationError{Msg: "身份证号最后一位必须为数字或X"}
		}
	} else {
		for _, c := range identity {
			if c < '0' || c > '9' {
				return &ValidationError{Msg: "身份证号必须为数字"}
			}
		}
	}
	return nil
}

// validateCard 验证银行卡号格式
func (h *HuifuController) validateCard(card string) error {
	length := len(card)
	if length < 16 || length > 19 {
		return &ValidationError{Msg: "银行卡号长度必须在16-19位之间"}
	}
	for _, c := range card {
		if c < '0' || c > '9' {
			return &ValidationError{Msg: "银行卡号只能包含数字"}
		}
	}
	return nil
}

type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string {
	return e.Msg
}
