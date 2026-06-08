package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetOrderList 获取订单列表（支持日期范围筛选和分页）
func GetOrderList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	merchsID := c.Query("merch_id")
	if merchsID == "" {
		response.Fail(c, http.StatusBadRequest, "商户ID不能为空", nil)
		return
	}

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	offset := (page - 1) * size

	// 根据日期范围确定查询时间
	var startTime, endTime time.Time
	useDateFilter := false

	if endDate != "" {
		// 使用日期范围过滤
		useDateFilter = true
		if startDate == "" {
			startDate = endDate // 如果没有开始日期，默认等于结束日期
		}
		startTime, _ = time.Parse("2006-01-02", startDate)
		endTime, _ = time.Parse("2006-01-02", endDate)
		endTime = endTime.Add(24 * time.Hour) // 结束日期包含当天
	}

	var orders []model.Order
	var total int64

	// 查询订单列表（按创建时间倒序）
	query := db.DB.Model(&model.Order{}).Where("merchs_id = ?", merchsID)

	if useDateFilter {
		query = query.Where("created_at >= ? AND created_at < ?", startTime, endTime)
	}

	if err := query.
		Count(&total).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&orders).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	// 转换为响应格式
	var list []gin.H
	for _, order := range orders {
		list = append(list, gin.H{
			"id":         order.ID,
			"code":       order.Code,
			"name":       order.Name,
			"status":     order.Status,
			"price":      order.Price,
			"amount":     order.Amount,
			"duration":   order.Duration,
			"userPhone":  order.UserPhone,
			"merchPhone": order.MerchPhone,
			"createdAt":  order.CreatedAt,
			"reqDate":    order.ReqDate,
			"startTime":  order.StartTime,
			"endTime":    order.EndTime,
		})
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":  list,
		"total": total,
	})
}

// GetOrderDetail 获取订单详情
func GetOrderDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "订单ID格式错误", nil)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单不存在", id)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"id":         order.ID,
		"code":       order.Code,
		"name":       order.Name,
		"status":     order.Status,
		"price":      order.Price,
		"deposit":    order.Deposit,
		"amount":     order.Amount,
		"duration":   order.Duration,
		"userPhone":  order.UserPhone,
		"merchPhone": order.MerchPhone,
		"createdAt":  order.CreatedAt,
		"reqDate":    order.ReqDate,
		"startTime":  order.StartTime,
		"endTime":    order.EndTime,
	})
}

// CreateOrder 创建订单
func CreateOrder(c *gin.Context) {
	var req struct {
		MerchsID  int32    `json:"merchsId" binding:"required"`
		UsersID   int32    `json:"usersId" binding:"required"`
		DevicesID int32    `json:"devicesId" binding:"required"`
		RoomsID   int32    `json:"roomsId" binding:"required"`
		GroupsID  int32    `json:"groupsId" binding:"required"`
		Name      string   `json:"name"`
		Amount    *int64   `json:"amount"`
		Duration  *int64   `json:"duration"`
		Price     *float64 `json:"price"`
		Deposit   *float64 `json:"deposit"`
		UserPhone string   `json:"userPhone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 生成订单编号
	now := time.Now()
	orderCode := fmt.Sprintf("ORD%s%04d", now.Format("20060102150405"), time.Now().UnixNano()%10000)

	order := model.Order{
		MerchsID:  req.MerchsID,
		UsersID:   req.UsersID,
		DevicesID: req.DevicesID,
		RoomsID:   req.RoomsID,
		GroupsID:  req.GroupsID,
		Name:      &req.Name,
		Code:      &orderCode,
		Status:    strPtr("0"), // 0: 未支付
		Duration:  req.Duration,
		Price:     req.Price,
		Deposit:   req.Deposit,
		UserPhone: strPtr(req.UserPhone),
		CreatedAt: &now,
	}

	if err := db.DB.Create(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建订单失败", err)
		return
	}

	response.SuccessWithMsg(c, "创建成功", gin.H{
		"id":   order.ID,
		"code": orderCode,
	})
}

// UpdateOrder 更新订单
func UpdateOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "订单ID格式错误", nil)
		return
	}

	var req struct {
		Status     string     `json:"status"`
		Price      *float64   `json:"price"`
		Duration   *int64     `json:"duration"`
		EndTime    *time.Time `json:"endTime"`
		UserPhone  string     `json:"userPhone"`
		MerchPhone string     `json:"merchPhone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单不存在", id)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	// 更新字段
	if req.Status != "" {
		order.Status = &req.Status
	}
	if req.Price != nil {
		order.Price = req.Price
	}
	if req.Duration != nil {
		order.Duration = req.Duration
	}
	if req.EndTime != nil {
		order.EndTime = req.EndTime
	}
	if req.UserPhone != "" {
		order.UserPhone = &req.UserPhone
	}
	if req.MerchPhone != "" {
		order.MerchPhone = &req.MerchPhone
	}

	if err := db.DB.Save(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新订单失败", err)
		return
	}

	response.SuccessWithMsg(c, "更新成功", gin.H{"success": true})
}

// GetRefundList 获取退款列表
func GetRefundList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	offset := (page - 1) * size

	var orders []model.Order
	var total int64

	query := db.DB.Model(&model.Order{}).Where("status IN (?, ?, ?)", "3", "4", "5") // 申请退款、已退款、拒绝退款

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.
		Count(&total).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&orders).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询退款列表失败", err)
		return
	}

	var list []gin.H
	for _, order := range orders {
		list = append(list, gin.H{
			"id":        order.ID,
			"code":      order.Code,
			"name":      order.Name,
			"status":    order.Status,
			"price":     order.Price,
			"userPhone": order.UserPhone,
			"createdAt": order.CreatedAt,
		})
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":  list,
		"total": total,
	})
}

// ApproveRefund 同意退款
func ApproveRefund(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "订单ID格式错误", nil)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单不存在", id)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	// 检查订单状态
	if order.Status == nil || *order.Status != "3" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法同意退款", nil)
		return
	}

	// 更新状态为已退款
	status := "4"
	order.Status = &status

	if err := db.DB.Save(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "同意退款失败", err)
		return
	}

	response.SuccessWithMsg(c, "退款成功", gin.H{"success": true})
}

// RejectRefund 拒绝退款
func RejectRefund(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "订单ID格式错误", nil)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单不存在", id)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	// 检查订单状态
	if order.Status == nil || *order.Status != "3" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法拒绝退款", nil)
		return
	}

	// 更新状态为拒绝退款
	status := "5"
	order.Status = &status

	if err := db.DB.Save(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "拒绝退款失败", err)
		return
	}

	response.SuccessWithMsg(c, "已拒绝退款", gin.H{"success": true})
}

// strPtr 辅助函数：将字符串转换为指针
func strPtr(s string) *string {
	return &s
}
