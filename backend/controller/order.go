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

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":  orders,
		"total": total,
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

	response.SuccessWithMsg(c, "获取成功", order)
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
	endTime := time.Now()
	order.EndTime = &endTime

	if err := db.DB.Save(&order).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "完成订单失败", err)
		return
	}

	response.SuccessWithMsg(c, "完成订单成功", gin.H{"success": true})
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
	if room.Status != "空闲" {
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
	if req.Mode == "single" || req.Mode == "pay_single" || req.Mode == "pay_time" {
		Status = "已完成"
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
		Deposit:    req.Deposit,
		UserPhone:  req.UserPhone,
		MerchPhone: req.MerchPhone,
		FreeTime:   now.Add(time.Minute * time.Duration(rule.FreeTime)),
		StartTime:  now,
		CreatedAt:  now,
	}

	// 根据规则类型确定订单结束时间
	if req.Mode == "single" || req.Mode == "pay_single" || req.Mode == "pay_time" {
		order.EndTime = &now
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "创建订单失败", err)
		return
	}

	// 根据规则类型确定更新柜子状态
	if req.Mode != "single" && req.Mode != "pay_single" {
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

// UserRentAndUnlock 用户端租用并开锁（免费模式完整流程）
// 流程：创建订单 -> 更新柜子状态 -> 发送开锁指令 -> 记录开锁日志 -> 返回结果
func UserRentAndUnlock(c *gin.Context) {
	var req struct {
		RoomID   int32   `json:"roomId" binding:"required"`
		MerchsID int32   `json:"merchsId" binding:"required"`
		UsersID  int32   `json:"usersId" binding:"required"`
		GroupsID int32   `json:"groupsId" binding:"required"`
		RulesID  int32   `json:"rulesId" binding:"required"`
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
	if room.Status != "空闲" {
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
	var orderDuration int32

	// 根据规则类型确定金额
	if rule.Type == "free" || req.Type == "free" {
		orderAmount = 0
	} else {
		if req.Amount > 0 {
			orderAmount = req.Amount
		} else {
			orderAmount = float64(rule.Price)
		}
	}

	// 根据规则模式设置时长
	if rule.AutoEndTime > 0 {
		orderDuration = rule.AutoEndTime
	}

	// 创建订单
	order := model.Order{
		MerchsID:  req.MerchsID,
		UsersID:   req.UsersID,
		RoomsID:   req.RoomID,
		GroupsID:  req.GroupsID,
		Name:      room.Name,
		Code:      orderCode,
		Status:    "进行中",
		Amount:    1,
		Duration:  orderDuration,
		Price:     orderAmount,
		Deposit:   req.Deposit,
		StartTime: now,
		CreatedAt: now,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "创建订单失败", err)
		return
	}

	// 更新柜子状态
	room.Status = "租用"
	room.UsersID = req.UsersID
	room.OrdersID = int32(order.ID)
	if err := tx.Save(&room).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "更新柜子状态失败", err)
		return
	}

	// 提交事务（订单和柜子状态已更新）
	if err := tx.Commit().Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "提交订单失败", err)
		return
	}

	// 发送开锁指令（在事务外执行，避免长时间锁表）
	var unlockSuccess bool
	var unlockMsg string
	var responseData string
	var err error

	if device.Code != "" {
		// 拼接开锁指令（板号+锁号）
		unlockCommand := room.BoardNo + room.LockNo + "01" // 01表示开锁
		// 调用 EleControlCommon 发送控制指令
		responseData, err = EleControlCommon(device.Code, unlockCommand, "")
		if err != nil {
			unlockSuccess = false
			unlockMsg = "开锁失败: " + err.Error()
		} else {
			unlockSuccess = true
			unlockMsg = "开锁成功"
		}
	} else {
		// 没有设备编码，模拟开锁成功
		unlockSuccess = true
		unlockMsg = "模拟开锁成功"
	}

	// 记录开锁日志
	deviceLog := model.Devicelog{
		MerchsID:   req.MerchsID,
		UsersID:    req.UsersID,
		DevicesID:  int32(device.ID),
		RoomID:     req.RoomID,
		Code:       device.Code,
		DeviceName: device.Name,
		RoomName:   room.Name,
		Type:       "手机", // 默认设备类型
		Control:    "开锁",
		Status:     map[bool]string{true: "成功", false: "失败"}[unlockSuccess],
		Occupant:   "用户",
		CreatedAt:  time.Now(),
	}

	if err := db.DB.Create(&deviceLog).Error; err != nil {
		// 日志记录失败不影响主流程
		fmt.Println("记录开锁日志失败:", err)
	}

	if unlockSuccess {
		response.SuccessWithMsg(c, "租用成功", gin.H{
			"orderId": order.ID,
			"code":    orderCode,
			"unlock": gin.H{
				"success":  true,
				"message":  unlockMsg,
				"response": responseData,
			},
		})
	} else {
		// 开锁失败，返回订单ID，让前端跳转订单详情
		response.SuccessWithMsg(c, "订单已创建，开锁失败", gin.H{
			"orderId": order.ID,
			"code":    orderCode,
			"unlock": gin.H{
				"success":  false,
				"message":  unlockMsg,
				"response": responseData,
			},
		})
	}
}

// UserOrderPayment 用户端订单支付（模拟支付）
func UserOrderPayment(c *gin.Context) {
	var req struct {
		UsersID int32   `json:"usersId" binding:"required"`
		OrderID uint64  `json:"orderId" binding:"required"`
		Amount  float64 `json:"amount"`
		Method  string  `json:"method"`
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

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, "支付失败", err)
		return
	}

	// 更新柜子状态（确保柜子处于租用状态）
	var room model.Room
	if err := tx.Where("id = ?", order.RoomsID).First(&room).Error; err == nil {
		if room.Status != "租用" {
			room.Status = "租用"
			room.UsersID = req.UsersID
			if err := tx.Save(&room).Error; err != nil {
				tx.Rollback()
				response.Fail(c, http.StatusInternalServerError, "更新柜子状态失败", err)
				return
			}
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
	})
}

// UserOrderRenew 用户端订单续费（预付费模式）
func UserOrderRenew(c *gin.Context) {
	var req struct {
		UsersID      int32   `json:"usersId" binding:"required"`
		OrderID      uint64  `json:"orderId" binding:"required"`
		Duration     int32   `json:"duration"`
		DurationUnit string  `json:"durationUnit"`
		Amount       float64 `json:"amount"`
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
		MerchsID int32   `json:"merchsId" binding:"required"`
		UsersID  int32   `json:"usersId" binding:"required"`
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
		MerchsID:  req.MerchsID,
		UsersID:   req.UsersID,
		GroupsID:  req.GroupsID,
		Name:      "押金",
		Code:      orderCode,
		Status:    "进行中",
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

// UserOrderEnd 用户端结束订单（计算费用）
func UserOrderEnd(c *gin.Context) {
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
		response.Fail(c, http.StatusBadRequest, "订单状态不正确，无法结束", nil)
		return
	}

	// 获取分组信息以获取规则ID
	var group model.Group
	var rule model.Rule
	if err := db.DB.Where("id = ?", order.GroupsID).First(&group).Error; err == nil && group.RulesID > 0 {
		// 获取规则信息
		db.DB.Where("id = ?", group.RulesID).First(&rule)
	}

	// 计算费用（根据规则类型）
	var totalCost float64
	var refundAmount float64

	// 获取使用时长（分钟）
	usageDuration := int32(time.Since(order.StartTime).Minutes())

	if rule.Type == "free" || rule.Type == "" {
		// 免费模式
		totalCost = 0
	} else if rule.Mode == "pay_single" {
		// 单次付费
		totalCost = float64(rule.Price)
	} else if rule.Mode == "pay_hourly" {
		// 按时付费
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

		// 计算实际计费时长（向上取整）
		if usageDuration <= 0 {
			usageDuration = 1
		}
		chargedDuration := (usageDuration + unitMinutes - 1) / unitMinutes
		totalCost = float64(chargedDuration) * float64(rule.Price)
	} else if rule.Mode == "pay_time" {
		// 预付费模式
		totalCost = order.Price

		// 如果支持自动退款，计算剩余时长退款
		if rule.AutoRefund && rule.AutoEndTime > 0 && order.Duration > 0 {
			remainingMinutes := order.Duration - usageDuration
			if remainingMinutes > 0 {
				// 按比例计算退款
				refundAmount = (float64(remainingMinutes) / float64(order.Duration)) * order.Price
			}
		}
	}

	tx := db.DB.Begin()

	// 更新订单状态和费用
	order.Status = "已完成"
	endTime := time.Now()
	order.EndTime = &endTime
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
		room.OrdersID = 0 // 重置订单ID
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
		"usageMinutes": usageDuration,
	})
}
