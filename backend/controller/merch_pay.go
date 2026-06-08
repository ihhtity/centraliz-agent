package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateMerchPay 创建商家消费订单
func CreateMerchPay(c *gin.Context) {
	var req struct {
		MerchsID int32   `json:"merchs_id" binding:"required"`
		Type     string  `json:"type" binding:"required"`
		Name     string  `json:"name" binding:"required"`
		Price    float64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", nil)
		return
	}

	// 验证商家是否存在
	var merch model.Merch
	if err := db.DB.Where("id = ?", req.MerchsID).First(&merch).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, "商家不存在", nil)
		return
	}

	// 生成订单号
	orderCode := uuid.New().String()[:24]

	// 创建订单
	merchPay := model.MerchPay{
		Code:          &orderCode,
		MerchsID:      &req.MerchsID,
		Name:          &req.Name,
		Type:          &req.Type,
		Price:         &req.Price,
		OriginalPrice: &req.Price,
		Status:        func() *string { s := "0"; return &s }(),
	}

	if err := db.DB.Create(&merchPay).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建订单失败", nil)
		return
	}

	response.SuccessWithMsg(c, "创建订单成功", merchPay)
}

// GetMerchPayList 获取商家消费订单列表
func GetMerchPayList(c *gin.Context) {
	var req struct {
		MerchsID int32 `form:"merchs_id" binding:"required"`
		Page     int   `form:"page" binding:"required"`
		Size     int   `form:"size" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", nil)
		return
	}

	// 验证商家是否存在
	var merch model.Merch
	if err := db.DB.Where("id = ?", req.MerchsID).First(&merch).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, "商家不存在", nil)
		return
	}

	// 计算分页
	offset := (req.Page - 1) * req.Size

	// 查询订单列表
	var orders []model.MerchPay
	var total int64

	query := db.DB.Model(&model.MerchPay{}).Where("merchs_id = ?", req.MerchsID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败", nil)
		return
	}

	// 获取订单列表（按创建时间倒序）
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.Size).Find(&orders).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败", nil)
		return
	}

	// 计算统计数据
	var totalAmount float64
	var todayAmount float64
	today := time.Now().Format("2006-01-02")

	if err := db.DB.Model(&model.MerchPay{}).
		Where("merchs_id = ? AND status = ?", req.MerchsID, "1").
		Select("COALESCE(SUM(price), 0)").Scan(&totalAmount).Error; err != nil {
		fmt.Println("统计总金额失败:", err)
	}

	if err := db.DB.Model(&model.MerchPay{}).
		Where("merchs_id = ? AND status = ? AND DATE(created_at) = ?", req.MerchsID, "1", today).
		Select("COALESCE(SUM(price), 0)").Scan(&todayAmount).Error; err != nil {
		fmt.Println("统计今日金额失败:", err)
	}

	response.SuccessWithMsg(c, "查询成功", gin.H{
		"list":        orders,
		"total":       total,
		"totalAmount": totalAmount,
		"todayAmount": todayAmount,
	})
}

// GetMerchPayDetail 获取订单详情
func GetMerchPayDetail(c *gin.Context) {
	var req struct {
		ID int32 `form:"id" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", nil)
		return
	}

	var order model.MerchPay
	if err := db.DB.Where("id = ?", req.ID).First(&order).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "订单不存在", nil)
		return
	}

	response.SuccessWithMsg(c, "查询成功", order)
}

// CancelMerchPay 取消订单
func CancelMerchPay(c *gin.Context) {
	var req struct {
		ID int32 `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", nil)
		return
	}

	var order model.MerchPay
	if err := db.DB.Where("id = ?", req.ID).First(&order).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "订单不存在", nil)
		return
	}

	// 只能取消待支付状态的订单
	if *order.Status != "0" {
		response.Fail(c, http.StatusBadRequest, "只能取消待支付状态的订单", nil)
		return
	}

	status := "2"
	if err := db.DB.Model(&order).Update("status", &status).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "取消订单失败", nil)
		return
	}

	response.SuccessWithMsg(c, "取消订单成功", nil)
}

// PayMerchPay 支付订单
func PayMerchPay(c *gin.Context) {
	var req struct {
		ID int32 `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", nil)
		return
	}

	var order model.MerchPay
	if err := db.DB.Where("id = ?", req.ID).First(&order).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "订单不存在", nil)
		return
	}

	// 只能支付待支付状态的订单
	if *order.Status != "0" {
		response.Fail(c, http.StatusBadRequest, "只能支付待支付状态的订单", nil)
		return
	}

	status := "1"
	if err := db.DB.Model(&order).Update("status", &status).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "支付失败", nil)
		return
	}

	response.SuccessWithMsg(c, "支付成功", nil)
}
