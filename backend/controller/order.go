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

// GetUserOrderList 获取用户端订单列表（支持状态筛选和分页）
func GetUserOrderList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status := c.Query("status")

	// 从JWT中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权", nil)
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

	var orders []model.Order
	var total int64

	// 查询订单列表（按创建时间倒序）
	query := db.DB.Model(&model.Order{}).Where("users_id = ?", userID)

	// 状态筛选
	if status != "" && status != "全部" {
		query = query.Where("status = ?", status)
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

// GetUserOrderDetail 获取用户端订单详情
func GetUserOrderDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "订单ID格式错误", nil)
		return
	}

	// 从JWT中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权", nil)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ? AND users_id = ?", id, userID).First(&order).Error; err != nil {
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

// UserApplyRefund 用户申请退款
func UserApplyRefund(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "订单ID格式错误", nil)
		return
	}

	// 从JWT中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权", nil)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ? AND users_id = ?", id, userID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单不存在", id)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	// 检查订单状态
	if order.Status != "未完成" && order.Status != "进行中" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法申请退款", nil)
		return
	}

	// 更新状态为申请退款
	order.Status = "申请退款"

	if err := db.DB.Save(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "申请退款失败", err)
		return
	}

	response.SuccessWithMsg(c, "申请退款成功", gin.H{"success": true})
}

// UserCancelRefund 用户取消退款申请
func UserCancelRefund(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "订单ID格式错误", nil)
		return
	}

	// 从JWT中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权", nil)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ? AND users_id = ?", id, userID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单不存在", id)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	// 检查订单状态
	if order.Status != "申请退款" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法取消退款", nil)
		return
	}

	// 更新状态为未完成
	order.Status = "未完成"

	if err := db.DB.Save(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "取消退款失败", err)
		return
	}

	response.SuccessWithMsg(c, "取消退款成功", gin.H{"success": true})
}

// UserCompleteOrder 用户完成订单
func UserCompleteOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "订单ID格式错误", nil)
		return
	}

	// 从JWT中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权", nil)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ? AND users_id = ?", id, userID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单不存在", id)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	// 检查订单状态
	if order.Status != "进行中" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法完成订单", nil)
		return
	}

	// 更新状态为已完成
	order.Status = "已完成"
	order.EndTime = time.Now()

	if err := db.DB.Save(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "完成订单失败", err)
		return
	}

	response.SuccessWithMsg(c, "完成订单成功", gin.H{"success": true})
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
		MerchsID  int32   `json:"merchsId" binding:"required"`
		UsersID   int32   `json:"usersId" binding:"required"`
		DevicesID int32   `json:"devicesId" binding:"required"`
		RoomsID   int32   `json:"roomsId" binding:"required"`
		GroupsID  int32   `json:"groupsId" binding:"required"`
		Name      string  `json:"name"`
		Amount    int32   `json:"amount"`
		Duration  int32   `json:"duration"`
		Price     float64 `json:"price"`
		Deposit   float64 `json:"deposit"`
		UserPhone string  `json:"userPhone"`
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
		Name:      req.Name,
		Code:      orderCode,
		Status:    "未完成",
		Amount:    req.Amount,
		Duration:  req.Duration,
		Price:     req.Price,
		Deposit:   req.Deposit,
		UserPhone: req.UserPhone,
		CreatedAt: now,
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
		Status     string    `json:"status"`
		Price      float64   `json:"price"`
		Duration   int32     `json:"duration"`
		EndTime    time.Time `json:"endTime"`
		UserPhone  string    `json:"userPhone"`
		MerchPhone string    `json:"merchPhone"`
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
		order.Status = req.Status
	}
	if req.Price != 0 {
		order.Price = req.Price
	}
	if req.Duration != 0 {
		order.Duration = req.Duration
	}
	if !req.EndTime.IsZero() {
		order.EndTime = req.EndTime
	}
	if req.UserPhone != "" {
		order.UserPhone = req.UserPhone
	}
	if req.MerchPhone != "" {
		order.MerchPhone = req.MerchPhone
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

	query := db.DB.Model(&model.Order{}).Where("status IN (?, ?, ?)", "申请退款", "已退款", "拒绝退款")

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
	if order.Status != "申请退款" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法同意退款", nil)
		return
	}

	// 更新状态为已退款
	order.Status = "已退款"

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
	if order.Status != "申请退款" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法拒绝退款", nil)
		return
	}

	// 更新状态为拒绝退款
	order.Status = "拒绝退款"

	if err := db.DB.Save(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "拒绝退款失败", err)
		return
	}

	response.SuccessWithMsg(c, "已拒绝退款", gin.H{"success": true})
}

// ==================== 用户端订单操作接口 ====================

// UserCreateOrder 用户端创建订单（支持不同规则类型）
func UserCreateOrder(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权", nil)
		return
	}

	var req struct {
		RoomID   uint64  `json:"roomId"`
		MerchsID int32   `json:"merchsId"`
		GroupsID int32   `json:"groupsId"`
		RulesID  int32   `json:"rulesId"`
		Mode     string  `json:"mode"`
		Type     string  `json:"type"`
		Amount   float64 `json:"amount"`
		Deposit  float64 `json:"deposit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 获取规则信息
	var rule model.Rule
	if err := db.DB.Where("id = ?", req.RulesID).First(&rule).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, "规则不存在", nil)
		return
	}

	// 获取柜子信息
	var room model.Room
	if err := db.DB.Where("id = ?", req.RoomID).First(&room).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, "柜子不存在", nil)
		return
	}

	// 检查柜子状态
	if room.Status != "空闲" {
		response.Fail(c, http.StatusBadRequest, "柜子当前不可用", nil)
		return
	}

	// 生成订单编号
	now := time.Now()
	orderCode := fmt.Sprintf("ORD%s%04d", now.Format("20060102150405"), time.Now().UnixNano()%10000)

	// 计算订单金额和时长
	var orderAmount float64
	var orderDuration int32

	if req.Type == "free" {
		// 免费模式
		orderAmount = 0
	} else {
		// 收费模式
		orderAmount = req.Amount
	}

	// 根据规则设置时长
	if rule.AutoEndTime > 0 {
		orderDuration = rule.AutoEndTime
	}

	// 创建订单
	order := model.Order{
		MerchsID:  req.MerchsID,
		UsersID:   int32(userID.(uint32)),
		DevicesID: 0,
		RoomsID:   int32(req.RoomID),
		GroupsID:  req.GroupsID,
		Name:      room.Name,
		Code:      orderCode,
		Status:    "进行中",
		Amount:    int32(req.Amount), // 使用传入的金额
		Duration:  orderDuration,
		Price:     orderAmount, // 使用计算的订单金额
		Deposit:   req.Deposit,
		StartTime: now,
		CreatedAt: now,
	}

	tx := db.DB.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "创建订单失败", err)
		return
	}

	// 更新柜子状态
	room.Status = "租用"
	room.UsersID = int32(userID.(uint32))
	if err := tx.Save(&room).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "更新柜子状态失败", err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "提交订单失败", err)
		return
	}

	response.SuccessWithMsg(c, "创建订单成功", gin.H{
		"orderId": order.ID,
		"code":    orderCode,
	})
}

// UserOrderPayment 用户端订单支付（模拟支付）
func UserOrderPayment(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权", nil)
		return
	}

	var req struct {
		OrderID uint64  `json:"orderId"`
		Amount  float64 `json:"amount"`
		Method  string  `json:"method"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ? AND users_id = ?", req.OrderID, userID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单不存在", nil)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	// 检查订单状态
	if order.Status != "未完成" && order.Status != "进行中" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法支付", nil)
		return
	}

	// 模拟支付处理（实际项目中这里应该调用第三方支付接口）
	// 模拟支付成功
	order.Status = "进行中"
	order.Price = req.Amount
	order.ReqDate = time.Now().Format("2006-01-02 15:04:05")

	if err := db.DB.Save(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "支付失败", err)
		return
	}

	response.SuccessWithMsg(c, "支付成功", gin.H{
		"orderId": order.ID,
		"status":  order.Status,
	})
}

// UserOrderRenew 用户端订单续费（预付费模式）
func UserOrderRenew(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权", nil)
		return
	}

	var req struct {
		OrderID      uint64  `json:"orderId"`
		Duration     int32   `json:"duration"`
		DurationUnit string  `json:"durationUnit"`
		Amount       float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ? AND users_id = ?", req.OrderID, userID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单不存在", nil)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	// 检查订单状态
	if order.Status != "进行中" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法续费", nil)
		return
	}

	// 更新订单时长和金额
	order.Duration += req.Duration
	order.Price += req.Amount

	if err := db.DB.Save(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "续费失败", err)
		return
	}

	response.SuccessWithMsg(c, "续费成功", gin.H{
		"orderId":  order.ID,
		"duration": order.Duration,
		"price":    order.Price,
	})
}

// ==================== 押金相关接口 ====================

// CheckDepositStatus 检查押金状态
func CheckDepositStatus(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权", nil)
		return
	}

	var req struct {
		MerchsID int32 `json:"merchsId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 查询用户在该商家下是否有已支付押金的订单
	var depositOrder model.Order
	err := db.DB.Where("users_id = ? AND merchs_id = ? AND status = ? AND deposit > ?",
		userID, req.MerchsID, "进行中", 0).First(&depositOrder).Error

	if err == gorm.ErrRecordNotFound {
		// 没有找到押金订单，返回未支付状态
		response.SuccessWithMsg(c, "获取成功", gin.H{
			"hasDeposit": false,
			"deposit":    0,
			"orderId":    0,
		})
		return
	}

	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败", err)
		return
	}

	// 找到押金订单，返回已支付状态
	response.SuccessWithMsg(c, "获取成功", gin.H{
		"hasDeposit": true,
		"deposit":    depositOrder.Deposit,
		"orderId":    depositOrder.ID,
	})
}

// PayDeposit 支付押金
func PayDeposit(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权", nil)
		return
	}

	var req struct {
		MerchsID int32   `json:"merchsId" binding:"required"`
		GroupsID int32   `json:"groupsId" binding:"required"`
		RulesID  int32   `json:"rulesId" binding:"required"`
		Amount   float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 获取规则信息
	var rule model.Rule
	if err := db.DB.Where("id = ?", req.RulesID).First(&rule).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, "规则不存在", nil)
		return
	}

	// 生成订单编号
	now := time.Now()
	orderCode := fmt.Sprintf("DEP%s%04d", now.Format("20060102150405"), time.Now().UnixNano()%10000)

	// 创建押金订单
	order := model.Order{
		MerchsID:  req.MerchsID,
		UsersID:   int32(userID.(uint32)),
		DevicesID: 0,
		RoomsID:   0,
		GroupsID:  req.GroupsID,
		Name:      "押金",
		Code:      orderCode,
		Status:    "进行中",
		Amount:    1,
		Duration:  0,
		Price:     0,
		Deposit:   req.Amount,
		StartTime: now,
		CreatedAt: now,
	}

	if err := db.DB.Create(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建押金订单失败", err)
		return
	}

	response.SuccessWithMsg(c, "押金支付成功", gin.H{
		"orderId": order.ID,
		"code":    orderCode,
		"deposit": req.Amount,
	})
}

// RefundDeposit 退还押金
func RefundDeposit(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权", nil)
		return
	}

	var req struct {
		OrderID uint64 `json:"orderId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ? AND users_id = ?", req.OrderID, userID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单不存在", nil)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	// 检查订单状态
	if order.Status != "进行中" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法退还押金", nil)
		return
	}

	// 更新订单状态为已完成
	order.Status = "已完成"
	order.EndTime = time.Now()

	if err := db.DB.Save(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "退还押金失败", err)
		return
	}

	response.SuccessWithMsg(c, "押金已退还", gin.H{
		"orderId": order.ID,
		"deposit": order.Deposit,
	})
}

// UserOrderEnd 用户端结束订单（计算费用）
func UserOrderEnd(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权", nil)
		return
	}

	var req struct {
		OrderID uint64 `json:"orderId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ? AND users_id = ?", req.OrderID, userID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单不存在", nil)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	// 检查订单状态
	if order.Status != "进行中" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法结束", nil)
		return
	}

	// 获取规则信息
	var rule model.Rule
	if err := db.DB.Where("id = ?", order.GroupsID).First(&rule).Error; err != nil {
		// 规则不存在，使用默认计算方式
	}

	// 计算费用（根据规则类型）
	var totalCost float64
	var refundAmount float64

	// 获取使用时长
	usageDuration := int32(time.Since(order.StartTime).Minutes())

	if rule.Type == "free" || rule.Type == "" {
		// 免费模式
		totalCost = 0
	} else if rule.Mode == "pay_single" {
		// 单次付费
		totalCost = float64(rule.Price)
	} else if rule.Mode == "pay_hourly" {
		// 按时付费
		// 计算实际使用时长（向上取整到计费单位）
		var unitMinutes int32
		switch rule.DurationUnit {
		case "minute":
			unitMinutes = 1
		case "hour":
			unitMinutes = 60
		case "day":
			unitMinutes = 1440
		default:
			unitMinutes = 60
		}

		// 计算实际计费时长
		chargedDuration := (usageDuration + unitMinutes - 1) / unitMinutes
		totalCost = float64(chargedDuration) * float64(rule.Price)
	} else if rule.Mode == "pay_time" {
		// 预付费模式
		totalCost = order.Price

		// 如果支持自动退款，计算剩余时长退款
		if rule.AutoRefund && rule.AutoEndTime > 0 {
			remainingMinutes := rule.AutoEndTime - usageDuration
			if remainingMinutes > 0 {
				// 按比例计算退款
				refundAmount = (float64(remainingMinutes) / float64(rule.AutoEndTime)) * order.Price
			}
		}
	}

	tx := db.DB.Begin()

	// 更新订单状态和费用
	order.Status = "已完成"
	order.EndTime = time.Now()
	order.Price = totalCost

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "更新订单失败", err)
		return
	}

	// 更新柜子状态
	var room model.Room
	if err := tx.Where("id = ?", order.RoomsID).First(&room).Error; err == nil {
		room.Status = "空闲"
		room.UsersID = 0
		if err := tx.Save(&room).Error; err != nil {
			tx.Rollback()
			response.Fail(c, http.StatusInternalServerError, "更新柜子状态失败", err)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "结束订单失败", err)
		return
	}

	response.SuccessWithMsg(c, "结束订单成功", gin.H{
		"orderId":      order.ID,
		"totalCost":    totalCost,
		"refundAmount": refundAmount,
		"deposit":      order.Deposit,
	})
}
