package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ==================== 商家订单相关接口 ====================

// GetOrderList 商家获取订单列表（支持日期范围筛选和分页）
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
	var total, refund int64
	var income, expense float64

	// 查询申请退款订单数量
	if err := db.DB.Model(&model.Order{}).Where("merchs_id = ? AND status = ?", merchsID, "申请退款").Select("COALESCE(COUNT(id), 0)").Scan(&refund).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询申请退款订单失败", err)
		return
	}
	// 查询订单列表（按创建时间倒序）
	query := db.DB.Model(&model.Order{}).Where("merchs_id = ?", merchsID)
	// 计算统计数据
	if useDateFilter {
		query = query.Where("created_at >= ? AND created_at < ?", startTime, endTime)
	}
	// 查询已退款订单的总金额（根据时间条件过滤）
	expenseQuery := db.DB.Model(&model.Order{}).Where("merchs_id = ? AND status = ?", merchsID, "已退款")
	if useDateFilter {
		expenseQuery = expenseQuery.Where("created_at >= ? AND created_at < ?", startTime, endTime)
	}
	if err := expenseQuery.Select("COALESCE(SUM(price), 0)").Scan(&expense).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询已退款订单失败", err)
		return
	}
	// 查询已完成订单的总金额（根据时间条件过滤）
	incomeQuery := db.DB.Model(&model.Order{}).Where("merchs_id = ? AND status = ?", merchsID, "已完成")
	if useDateFilter {
		incomeQuery = incomeQuery.Where("created_at >= ? AND created_at < ?", startTime, endTime)
	}
	if err := incomeQuery.Select("COALESCE(SUM(price), 0)").Scan(&income).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询已完成订单失败", err)
		return
	}
	// 查询订单列表
	if err := query.
		Count(&total).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&orders).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":    orders,
		"total":   total,
		"income":  income,
		"expense": expense,
		"refund":  refund,
	})
}

// GetOrderDetail 商家获取订单详情
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

	// 获取房间名称
	var roomName string
	var room model.Room
	if order.RoomsID > 0 {
		if err := db.DB.Where("id = ?", order.RoomsID).First(&room).Error; err == nil {
			roomName = room.Name
		}
	}

	// 获取设备名称
	var deviceName string
	if room.DevicesID > 0 {
		var device model.Device
		if err := db.DB.Where("id = ?", room.DevicesID).First(&device).Error; err == nil {
			deviceName = device.Name
		}
	}

	// 返回包含设备名称和房间名称的订单详情
	response.SuccessWithMsg(c, "获取成功", gin.H{
		"id":         order.ID,
		"code":       order.Code,
		"name":       order.Name,
		"status":     order.Status,
		"amount":     order.Amount,
		"duration":   order.Duration,
		"price":      order.Price,
		"deposit":    order.Deposit,
		"userPhone":  order.UserPhone,
		"merchPhone": order.MerchPhone,
		"reqDate":    order.ReqDate,
		"freeTime":   order.FreeTime,
		"startTime":  order.StartTime,
		"endTime":    order.EndTime,
		"createdAt":  order.CreatedAt,
		"deviceName": deviceName,
		"roomName":   roomName,
	})
}

// CreateOrder 商家创建订单
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

// UpdateOrder 商家更新订单
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
		order.EndTime = &req.EndTime
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

// GetRefundList 商家获取退款列表
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

// ApproveRefund 商家同意退款
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

// RejectRefund 商家拒绝退款
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

// GetUserOrderList 获取用户端订单列表（支持状态筛选和分页）
func GetUserOrderList(c *gin.Context) {
	var req struct {
		UsersID int32  `json:"usersId"`
		Page    int32  `json:"page"`
		Size    int32  `json:"size"`
		Status  string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 20
	}
	if req.Size > 100 {
		req.Size = 100
	}

	offset := (req.Page - 1) * req.Size

	var orders []model.Order
	var total int64

	// 查询订单列表（按创建时间倒序）
	query := db.DB.Model(&model.Order{}).Where("users_id = ?", req.UsersID)

	// 状态筛选
	if req.Status != "" && req.Status != "全部" {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.
		Count(&total).
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(req.Size)).
		Find(&orders).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":  orders,
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

	var order model.Order
	if err := db.DB.Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单不存在", id)
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单失败", err)
		return
	}

	response.SuccessWithMsg(c, "获取成功", order)
}

// UserApplyRefund 用户申请退款
func UserApplyRefund(c *gin.Context) {
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

	// 更新标签为申请退款
	var time = time.Now()
	order.Status = "已完成"
	order.Tag = "申请退款"
	if order.EndTime == nil {
		order.EndTime = &time
		order.Duration = int32(time.Sub(order.StartTime).Minutes())
	}

	// 添加事务，确保数据库操作原子性
	tx := db.DB.Begin()
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "申请退款失败", err)
		return
	}

	// 创建新的退款订单
	newOrder := model.Order{
		MerchsID:   order.MerchsID,
		UsersID:    order.UsersID,
		RoomsID:    order.RoomsID,
		GroupsID:   order.GroupsID,
		Name:       order.Name + "-申请退款",
		Code:       fmt.Sprintf("ORD%s%04d", time.Format("20060102150405"), time.UnixNano()%10000),
		Status:     "申请退款",
		Tag:        "申请退款",
		Amount:     order.Amount,
		Price:      order.Price,
		Deposit:    order.Deposit,
		UserPhone:  order.UserPhone,
		MerchPhone: order.MerchPhone,
		ReqDate:    time.Format("2006-01-02 15:04:05"),
		FreeTime:   time,
		StartTime:  time,
		EndTime:    &time,
	}

	if err := tx.Create(&newOrder).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "创建退款订单失败", err)
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "提交退款订单失败", err)
		return
	}

	response.SuccessWithMsg(c, "申请成功", gin.H{"success": true})
}

// UserCompleteOrder 用户完成订单
func UserCompleteOrder(c *gin.Context) {
	var req struct {
		RoomID  int32  `json:"roomId" binding:"required"`
		OrderID int32  `json:"orderId" binding:"required"`
		Mode    string `json:"mode" binding:"required"`
	}

	// 解析请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ?", req.OrderID).First(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "订单不存在", err)
		return
	}

	// 检查订单状态
	if order.Status != "进行中" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法完成订单", nil)
		return
	}

	// 获取柜子信息
	var room model.Room
	if err := db.DB.Where("id = ?", req.RoomID).First(&room).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, "柜子不存在", nil)
		return
	}

	// 开始事务
	tx := db.DB.Begin()

	// 更新状态为已完成
	order.Status = "已完成"
	if req.Mode != "pay_time" {
		order.Price = 0.00
	}

	endTime := time.Now()
	order.EndTime = &endTime
	order.FreeTime = endTime
	order.Duration = int32(endTime.Sub(order.StartTime).Minutes())

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "完成订单失败", err)
		return
	}

	// 更新柜子状态
	room.Status = "空闲"
	room.Combo = ""
	room.UsersID = 0
	room.OrdersID = 0
	room.FreeTime = endTime
	if err := tx.Save(&room).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "更新柜子状态失败", err)
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "完成订单失败", err)
		return
	}

	// 提交事务成功后，返回成功响应
	response.SuccessWithMsg(c, "完成订单成功", gin.H{"orderId": order.ID, "price": order.Price, "duration": order.Duration})
}

// UserCreateOrder 用户端创建订单（支持不同规则类型）
func UserCreateOrder(c *gin.Context) {
	var req struct {
		RoomID     int32   `json:"roomId" binding:"required"`
		MerchsID   int32   `json:"merchsId" binding:"required"`
		UsersID    int32   `json:"usersId" binding:"required"`
		GroupsID   int32   `json:"groupsId" binding:"required"`
		RulesID    int32   `json:"rulesId" binding:"required"`
		Mode       string  `json:"mode"`
		Type       string  `json:"type"`
		OrderType  string  `json:"orderType"`
		Amount     float64 `json:"amount"`
		Deposit    float64 `json:"deposit"`
		UserPhone  string  `json:"userPhone"`
		MerchPhone string  `json:"merchPhone"`
	}

	// 解析请求参数
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

	// 获取设备信息
	var device model.Device
	if err := db.DB.Where("id = ?", room.DevicesID).First(&device).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, "设备不存在", nil)
		return
	}

	// 检查柜子状态（使用数据库事务确保并发安全）
	tx := db.DB.Begin()

	// 重新获取柜子状态（确保最新）
	if err := tx.Where("id = ?", req.RoomID).First(&room).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusBadRequest, "柜子不存在", nil)
		return
	}

	// 检查柜子状态
	if room.Status != "空闲" && req.OrderType != "renew" {
		tx.Rollback()
		response.Fail(c, http.StatusBadRequest, "柜子当前不可用", nil)
		return
	}

	// 检查柜子所属分组是否匹配
	if room.GroupsID != req.GroupsID {
		tx.Rollback()
		response.Fail(c, http.StatusBadRequest, "柜子与分组不匹配", nil)
		return
	}

	// 生成订单编号
	now := time.Now()
	orderCode := fmt.Sprintf("ORD%s%04d", now.Format("20060102150405"), time.Now().UnixNano()%10000)

	// 计算订单金额和时长
	var orderAmount float64
	// 订单状态默认为进行中
	var Status = "进行中"

	// 根据规则类型确定订单状态
	if req.Mode == "single" {
		Status = "已完成"
	}
	if req.Mode == "pay_single" {
		Status = "未完成"
	}

	// 根据规则类型确定金额
	if rule.Type == "free" || req.Type == "free" {
		// 免费模式
		orderAmount = 0
	} else {
		// 收费模式：优先使用传入的金额，否则使用规则配置的价格
		if req.Amount > 0 {
			orderAmount = req.Amount
		} else {
			orderAmount = float64(rule.Price)
		}
	}

	// 创建订单
	order := model.Order{
		MerchsID:   req.MerchsID,
		UsersID:    req.UsersID,
		RoomsID:    req.RoomID,
		GroupsID:   req.GroupsID,
		Name:       room.Name,
		Code:       orderCode,
		Status:     Status,
		Amount:     1, // 商品数量默认为1
		Duration:   0,
		Price:      orderAmount,
		Deposit:    0.00,
		UserPhone:  req.UserPhone,
		MerchPhone: req.MerchPhone,
		FreeTime:   now.Add(time.Minute * time.Duration(rule.FreeTime)),
		StartTime:  now,
		CreatedAt:  now,
	}

	// 根据规则类型确定订单结束时间
	if req.Mode == "single" || req.Mode == "pay_single" {
		order.EndTime = &now
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "创建订单失败", err)
		return
	}

	// 根据规则类型确定更新柜子状态
	if req.Mode != "single" && req.Mode != "pay_single" {
		room.Combo = now.Format("2006-01-02 15:04:05")
		room.Status = "租用"
		room.UsersID = req.UsersID
		room.OrdersID = int32(order.ID) // 设置订单ID
		room.FreeTime = now.Add(time.Minute * time.Duration(rule.FreeTime))
		if err := tx.Save(&room).Error; err != nil {
			tx.Rollback()
			response.Fail(c, http.StatusInternalServerError, "更新柜子状态失败", err)
			return
		}
	}

	// 提交事务（订单和柜子状态已更新）
	if err := tx.Commit().Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "提交订单失败", err)
		return
	}

	response.SuccessWithMsg(c, "创建订单成功", gin.H{
		"orderId": order.ID,
		"code":    orderCode,
	})
}

// UserOrderPayment 用户端订单支付开始
func UserOrderPayment(c *gin.Context) {
	var req struct {
		UsersID int32   `json:"usersId" binding:"required"`
		OrderID uint64  `json:"orderId" binding:"required"`
		Amount  float64 `json:"amount"`
		Method  string  `json:"method"`
		Mode    string  `json:"mode"`
		Combo   string  `json:"combo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ? AND users_id = ?", req.OrderID, req.UsersID).First(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "订单不存在", err)
		return
	}

	// 检查订单状态
	if order.Status != "未完成" && order.Status != "进行中" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法支付", nil)
		return
	}

	// 验证支付金额
	if req.Amount <= 0 {
		response.Fail(c, http.StatusBadRequest, "支付金额必须大于0", nil)
		return
	}

	// 验证支付方式
	if req.Method == "" {
		req.Method = "wechat" // 默认微信支付
	}

	tx := db.DB.Begin()

	// 更新订单状态和金额
	order.Status = "进行中"
	order.Price = req.Amount
	order.ReqDate = time.Now().Format("2006-01-02 15:04:05")

	if req.Mode == "pay_single" {
		order.Status = "已完成"
	}

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "支付失败", err)
		return
	}

	// 更新柜子状态（确保柜子处于租用状态）
	var room model.Room
	if err := tx.Where("id = ?", order.RoomsID).First(&room).Error; err == nil {
		if req.Mode == "pay_single" {
			room.Status = "空闲"
		} else {
			room.Status = "租用"
			room.UsersID = req.UsersID
			room.OrdersID = int32(order.ID)
		}
		room.Combo = req.Combo

		if err := tx.Save(&room).Error; err != nil {
			tx.Rollback()
			response.Fail(c, http.StatusInternalServerError, "更新柜子状态失败", err)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "支付失败", err)
		return
	}

	response.SuccessWithMsg(c, "支付成功", gin.H{
		"orderId": order.ID,
		"status":  order.Status,
		"amount":  req.Amount,
		"method":  req.Method,
		"mode":    req.Mode,
	})
}

// UserOrderEnd 用户端订单支付结束
func UserOrderEnd(c *gin.Context) {
	var req struct {
		OrderID uint64  `json:"orderId" binding:"required"`
		Amount  float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ?", req.OrderID).First(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "订单不存在", err)
		return
	}

	// 检查订单状态
	if order.Status != "进行中" {
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法结束", nil)
		return
	}

	// 验证支付金额
	if req.Amount <= 0 {
		response.Fail(c, http.StatusBadRequest, "支付金额必须大于0", nil)
		return
	}

	tx := db.DB.Begin()

	// 更新订单状态和费用
	order.Status = "已完成"
	endTime := time.Now()
	order.ReqDate = time.Now().Format("2006-01-02 15:04:05")
	order.EndTime = &endTime
	order.Duration = int32(endTime.Sub(order.StartTime).Minutes())
	order.Price = req.Amount

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "更新订单失败", err)
		return
	}

	// 更新柜子状态
	var room model.Room
	if err := tx.Where("id = ?", order.RoomsID).First(&room).Error; err == nil {
		room.Status = "空闲"
		room.Combo = ""
		room.UsersID = 0
		room.OrdersID = 0
		room.FreeTime = endTime
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
		"orderId":  order.ID,
		"duration": order.Duration,
		"price":    order.Price,
	})
}

// UserOrderRenew 用户端订单续费（预付费模式）
func UserOrderRenew(c *gin.Context) {
	var req struct {
		OrderID uint64  `json:"orderId" binding:"required"`
		Amount  float64 `json:"amount" binding:"required"`
		Method  string  `json:"method"`
		Combo   string  `json:"combo"` // combo JSON字符串
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 查询订单
	var order model.Order
	if err := db.DB.Where("id = ?", req.OrderID).First(&order).Error; err != nil {
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

	// 获取柜子信息
	var room model.Room
	if err := db.DB.Where("id = ?", order.RoomsID).First(&room).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, "柜子不存在", nil)
		return
	}

	// 解析 combo 数据
	var comboData map[string]interface{}
	if req.Combo != "" {
		if err := json.Unmarshal([]byte(req.Combo), &comboData); err != nil {
			response.Fail(c, http.StatusBadRequest, "续费数据格式错误", err)
			return
		}
	}

	// 结束旧订单、创建新订单、更新柜子（使用数据库事务确保并发安全）
	tx := db.DB.Begin()

	// 1. 更新旧订单为已完成
	oldOrderEndTime := time.Now()
	order.Status = "已完成"
	order.EndTime = &oldOrderEndTime
	order.Duration = int32(oldOrderEndTime.Sub(order.StartTime).Minutes())
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "完成订单失败", err)
		return
	}

	// 2. 创建新的续费订单
	newOrder := model.Order{
		MerchsID:   order.MerchsID,
		UsersID:    order.UsersID,
		RoomsID:    order.RoomsID,
		GroupsID:   order.GroupsID,
		Name:       room.Name + "-续费",
		Code:       fmt.Sprintf("ORD%s%04d", oldOrderEndTime.Format("20060102150405"), time.Now().UnixNano()%10000),
		Status:     "进行中",
		Amount:     order.Amount,
		Price:      req.Amount,
		UserPhone:  order.UserPhone,
		MerchPhone: order.MerchPhone,
		ReqDate:    oldOrderEndTime.Format("2006-01-02 15:04:05"),
		FreeTime:   oldOrderEndTime,
		StartTime:  order.StartTime,
	}

	// 解析新订单的结束时间（前端传来的格式：2026-06-24 01:02:32）
	if endTimeStr, ok := comboData["endTime"].(string); ok {
		if parsedEndTime, err := time.ParseInLocation("2006-01-02 15:04:05", endTimeStr, time.Local); err == nil {
			newOrder.EndTime = &parsedEndTime
		}
	}

	if err := tx.Create(&newOrder).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "创建续费订单失败", err)
		return
	}

	// 3. 更新柜子
	room.OrdersID = int32(newOrder.ID)
	room.Combo = req.Combo
	room.FreeTime = oldOrderEndTime
	if err := tx.Save(&room).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "更新柜子状态失败", err)
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "提交订单失败", err)
		return
	}

	// 格式化时间为 RFC3339 格式（如：2026-06-23T11:55:32.919+08:00）
	var endTimeStr, startTimeStr string
	if newOrder.EndTime != nil {
		endTimeStr = newOrder.EndTime.Format(time.RFC3339)
	}
	startTimeStr = newOrder.StartTime.Format(time.RFC3339)

	response.SuccessWithMsg(c, "续费成功", gin.H{
		"orderId":   newOrder.ID,
		"price":     newOrder.Price,
		"endTime":   endTimeStr,
		"startTime": startTimeStr,
	})
}

// ==================== 押金相关接口 ====================

// CheckDepositStatus 检查押金状态
func CheckDepositStatus(c *gin.Context) {
	var req struct {
		MerchsID int32 `json:"merchsId"`
		UsersID  int32 `json:"usersId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 查询用户的押金订单
	var depositOrder model.Order
	var err error

	if req.MerchsID > 0 {
		// 按商家ID查询
		err = db.DB.Where("users_id = ? AND merchs_id = ? AND status = ? AND deposit > ?",
			req.UsersID, req.MerchsID, "进行中", 0).First(&depositOrder).Error
	} else {
		// 查询用户所有押金订单
		err = db.DB.Where("users_id = ? AND status = ? AND deposit > ?",
			req.UsersID, "进行中", 0).First(&depositOrder).Error
	}

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
		"code":       depositOrder.Code,
	})
}

// PayDeposit 支付押金
func PayDeposit(c *gin.Context) {
	var req struct {
		MerchsID   int32   `json:"merchsId" binding:"required"`
		UsersID    int32   `json:"usersId" binding:"required"`
		GroupsID   int32   `json:"groupsId" binding:"required"`
		RulesID    int32   `json:"rulesId" binding:"required"`
		Amount     float64 `json:"amount" binding:"required"`
		MerchPhone string  `json:"merchPhone"`
		UserPhone  string  `json:"userPhone"`
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

	// 验证押金金额是否与规则匹配
	if rule.Deposit > 0 && req.Amount != float64(rule.Deposit) {
		response.Fail(c, http.StatusBadRequest, "押金金额与规则不符", nil)
		return
	}

	// 生成订单编号
	now := time.Now()
	orderCode := fmt.Sprintf("DEP%s%04d", now.Format("20060102150405"), time.Now().UnixNano()%10000)

	// 创建押金订单
	order := model.Order{
		MerchsID:   req.MerchsID,
		UsersID:    req.UsersID,
		GroupsID:   req.GroupsID,
		Name:       "押金",
		Code:       orderCode,
		Status:     "进行中",
		Tag:        "申请退款",
		Deposit:    req.Amount,
		UserPhone:  req.UserPhone,
		MerchPhone: req.MerchPhone,
		ReqDate:    now.Format("2006-01-02 15:04:05"),
		FreeTime:   now,
		StartTime:  now,
		CreatedAt:  now,
		EndTime:    nil,
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
	var req struct {
		UsersID int32  `json:"usersId" binding:"required"`
		OrderID uint64 `json:"orderId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	var order model.Order
	if err := db.DB.Where("id = ? AND users_id = ?", req.OrderID, req.UsersID).First(&order).Error; err != nil {
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
	endTime := time.Now()
	order.EndTime = &endTime

	if err := db.DB.Save(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "退还押金失败", err)
		return
	}

	response.SuccessWithMsg(c, "押金已退还", gin.H{
		"orderId": order.ID,
		"deposit": order.Deposit,
	})
}
