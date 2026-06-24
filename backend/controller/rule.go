package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetRuleList 获取规则列表
func GetRuleList(c *gin.Context) {
	merchsID := c.Query("merchs_id")
	if merchsID == "" {
		response.Fail(c, http.StatusBadRequest, "商户ID不能为空", nil)
		return
	}

	var rules []model.Rule
	if err := db.DB.Where("merchs_id = ?", merchsID).Order("created_at DESC").Find(&rules).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询规则失败", err)
		return
	}

	response.SuccessWithMsg(c, "获取成功", rules)
}

// GetRuleDetail 获取规则详情
func GetRuleDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "规则ID格式错误", nil)
		return
	}

	var rule model.Rule
	if err := db.DB.Where("id = ?", id).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "规则不存在", nil)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询规则失败", err)
		return
	}

	response.SuccessWithMsg(c, "获取成功", rule)
}

// CreateRule 创建规则
func CreateRule(c *gin.Context) {
	var req struct {
		MerchsID     int32   `json:"merchsId" binding:"required"`
		Name         string  `json:"name" binding:"required"`
		Type         string  `json:"type" binding:"required"` // free/charge
		Mode         string  `json:"mode" binding:"required"` // single/deposit/pay_single/pay_deposit/pay_hourly/pay_time
		Price        float32 `json:"price"`
		Deposit      float32 `json:"deposit"`
		DurationUnit string  `json:"durationUnit"`
		AutoEndTime  int32   `json:"autoEndTime"`
		FreeTime     int32   `json:"freeTime"`
		AutoRefund   bool    `json:"autoRefund"`
		ManualRenew  bool    `json:"manualRenew"`
		Description  string  `json:"description"`
		TimeOptions  string  `json:"timeOptions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 验证规则类型和模式类型的合法性
	if !isValidRuleType(req.Type) {
		response.Fail(c, http.StatusBadRequest, "无效的规则类型", nil)
		return
	}
	if !isValidRuleMode(req.Mode) {
		response.Fail(c, http.StatusBadRequest, "无效的模式类型", nil)
		return
	}
	if !isModeTypeMatch(req.Type, req.Mode) {
		response.Fail(c, http.StatusBadRequest, "规则类型与模式类型不匹配", nil)
		return
	}

	now := time.Now()
	rule := model.Rule{
		MerchsID:     req.MerchsID,
		Name:         req.Name,
		Mode:         req.Mode,
		Price:        req.Price,
		Deposit:      req.Deposit,
		DurationUnit: req.DurationUnit,
		AutoEndTime:  req.AutoEndTime,
		FreeTime:     req.FreeTime,
		AutoRefund:   req.AutoRefund,
		ManualRenew:  req.ManualRenew,
		TimeOptions:  req.TimeOptions,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 按时付费模式需要时间单位，其他模式不需要
	if req.Mode == "pay_hourly" && rule.DurationUnit == "" {
		rule.DurationUnit = "hour"
	}

	// 如果没有提供描述，则使用默认描述
	if req.Description == "" {
		rule.Description = getDefaultDescription(req.Mode)
	} else {
		rule.Description = req.Description
	}

	// 预付费模式验证
	if req.Mode == "pay_time" && req.TimeOptions == "" {
		response.Fail(c, http.StatusBadRequest, "预付费模式必须提供时间选项", nil)
		return
	}

	if err := db.DB.Create(&rule).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建规则失败", err)
		return
	}

	response.SuccessWithMsg(c, "创建成功", rule)
}

// UpdateRule 更新规则
func UpdateRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "规则ID格式错误", nil)
		return
	}

	var req struct {
		Name         string  `json:"name"`
		Mode         string  `json:"mode"`
		Type         string  `json:"type"`
		Price        float32 `json:"price"`
		Deposit      float32 `json:"deposit"`
		DurationUnit string  `json:"durationUnit"`
		AutoEndTime  int32   `json:"autoEndTime"`
		FreeTime     int32   `json:"freeTime"`
		AutoRefund   bool    `json:"autoRefund"`
		ManualRenew  bool    `json:"manualRenew"`
		Description  string  `json:"description"`
		TimeOptions  string  `json:"timeOptions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	var rule model.Rule
	if err := db.DB.Where("id = ?", id).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "规则不存在", nil)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询规则失败", err)
		return
	}

	// 验证模式类型的合法性（如果提供了模式）
	if req.Mode != "" {
		if !isValidRuleMode(req.Mode) {
			response.Fail(c, http.StatusBadRequest, "无效的模式类型", nil)
			return
		}
	}

	// 更新字段
	rule.Name = req.Name
	rule.Type = req.Type
	rule.Mode = req.Mode
	rule.Price = req.Price
	rule.Deposit = req.Deposit

	rule.DurationUnit = req.DurationUnit
	if req.Mode == "pay_hourly" && rule.DurationUnit == "" {
		// 切换到按时付费模式时，如果没有提供时间单位，设置默认值
		rule.DurationUnit = "hour"
	}
	rule.AutoEndTime = req.AutoEndTime
	rule.FreeTime = req.FreeTime
	// 对于 bool 类型，直接赋值（因为 bool 的零值 false 也是有效值）
	rule.AutoRefund = req.AutoRefund
	rule.ManualRenew = req.ManualRenew
	rule.TimeOptions = req.TimeOptions
	rule.Description = req.Description
	if req.Mode != "" {
		// 如果更新了模式但没提供描述，自动生成默认描述
		rule.Description = getDefaultDescription(req.Mode)
	}

	// 预付费模式验证
	if rule.Mode == "pay_time" && rule.TimeOptions == "" {
		response.Fail(c, http.StatusBadRequest, "预付费模式必须提供时间选项", nil)
		return
	}

	now := time.Now()
	rule.UpdatedAt = now

	if err := db.DB.Save(&rule).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新规则失败", err)
		return
	}

	response.SuccessWithMsg(c, "更新成功", rule)
}

// DeleteRule 删除规则
func DeleteRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "规则ID格式错误", nil)
		return
	}

	var rule model.Rule
	if err := db.DB.Where("id = ?", id).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "规则不存在", nil)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询规则失败", err)
		return
	}

	if err := db.DB.Delete(&rule).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除规则失败", err)
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

// isValidRuleType 验证规则类型是否有效
func isValidRuleType(ruleType string) bool {
	return ruleType == "free" || ruleType == "charge"
}

// isValidRuleMode 验证模式类型是否有效
func isValidRuleMode(mode string) bool {
	validModes := []string{"single", "deposit", "pay_single", "pay_deposit", "pay_hourly", "pay_time"}
	for _, m := range validModes {
		if m == mode {
			return true
		}
	}
	return false
}

// isModeTypeMatch 验证模式类型与规则类型是否匹配
func isModeTypeMatch(ruleType, mode string) bool {
	// 免费模式对应的模式
	freeModes := []string{"single", "deposit"}
	// 收费模式对应的模式
	chargeModes := []string{"pay_single", "pay_deposit", "pay_hourly", "pay_time"}

	if ruleType == "free" {
		for _, m := range freeModes {
			if m == mode {
				return true
			}
		}
	} else if ruleType == "charge" {
		for _, m := range chargeModes {
			if m == mode {
				return true
			}
		}
	}
	return false
}

// getDefaultDescription 根据模式类型获取默认描述
func getDefaultDescription(mode string) string {
	switch mode {
	case "single":
		return "免费模式:您可以免费开锁，无需支付任何费用"
	case "deposit":
		return "免费模式:您可以免费存取，无需支付任何费用"
	case "pay_single":
		return "单次付费模式:您需要在开锁时支付费用，开锁后即可使用"
	case "pay_deposit":
		return "先存后取模式:您可以先存放物品，取走时需要支付费用才能打开。如果涉及押金可以手动申请退款"
	case "pay_hourly":
		return "按时付费模式:您可以先存放物品，取走时需要支付所有时长对应的费用才能打开。如果涉及押金可以手动申请退款"
	case "pay_time":
		return "预付费模式:您需要先选择时间套餐并支付费用，支付完成后才能开锁使用"
	default:
		return ""
	}
}
