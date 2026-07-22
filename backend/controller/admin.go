package controller

import (
	"centraliz-backend/logic"
	"centraliz-backend/model"
	"centraliz-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ==================== 房间管理 ====================

type RoomSearchRequest struct {
	MerchsID int32  `form:"merchs_id"`
	GroupsID int32  `form:"groups_id"`
	Name     string `form:"name"`
	Status   string `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

func AdminGetRoomList(c *gin.Context) {
	merchsIDStr := c.Query("merchs_id")
	groupsIDStr := c.Query("groups_id")
	name := c.Query("name")
	status := c.Query("status")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	merchsID := int32(0)
	groupsID := int32(0)
	page := 1
	pageSize := 10

	if merchsIDStr != "" {
		if v, err := strconv.Atoi(merchsIDStr); err == nil {
			merchsID = int32(v)
		}
	}

	if groupsIDStr != "" {
		if v, err := strconv.Atoi(groupsIDStr); err == nil {
			groupsID = int32(v)
		}
	}

	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			page = v
		}
	}

	if pageSizeStr != "" {
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	rooms, total, err := logic.GetRoomListFiltered(merchsID, groupsID, name, status, page, pageSize)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  rooms,
		"total": total,
	})
}

func AdminGetRoomDetail(c *gin.Context) {
	id := c.Param("id")
	room, err := logic.GetRoomByID(id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "获取成功", room)
}

func AdminCreateRoom(c *gin.Context) {
	var req struct {
		Name      string  `json:"name" binding:"required"`
		MerchsID  int32   `json:"merchs_id" binding:"required"`
		GroupsID  *int32  `json:"groups_id"`
		RulesID   int32   `json:"rules_id"`
		DevicesID int32   `json:"devices_id"`
		Tag       string  `json:"tag"`
		BoardNo   string  `json:"board_no"`
		LockNo    string  `json:"lock_no"`
		Price     float32 `json:"price"`
		Image     string  `json:"image"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	room := &model.Room{
		Name:      req.Name,
		MerchsID:  req.MerchsID,
		DevicesID: req.DevicesID,
		RulesID:   req.RulesID,
		Tag:       req.Tag,
		BoardNo:   req.BoardNo,
		LockNo:    req.LockNo,
		Price:     req.Price,
		Image:     req.Image,
		Status:    "空闲",
	}

	if req.GroupsID != nil && *req.GroupsID > 0 {
		room.GroupsID = *req.GroupsID
	}

	if err := logic.CreateRoom(room); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", room)
}

func AdminUpdateRoom(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name   string  `json:"name"`
		Tag    string  `json:"tag"`
		Status string  `json:"status"`
		LockNo string  `json:"lock_no"`
		Price  float32 `json:"price"`
		Image  string  `json:"image"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	room, err := logic.GetRoomByID(id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}

	if req.Name != "" {
		room.Name = req.Name
	}
	if req.Tag != "" {
		room.Tag = req.Tag
	}
	if req.Status != "" {
		room.Status = req.Status
	}
	if req.LockNo != "" {
		room.LockNo = req.LockNo
	}
	if req.Price >= 0 {
		room.Price = req.Price
	}
	if req.Image != "" {
		room.Image = req.Image
	}

	if err := logic.UpdateRoom(room); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", room)
}

func AdminDeleteRoom(c *gin.Context) {
	id := c.Param("id")
	if err := logic.DeleteRoom(id); err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteRoom(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要删除的房间")
		return
	}
	if err := logic.BatchDeleteRoom(req.IDs); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "批量删除成功", nil)
}

func AdminBatchUpdateRoom(c *gin.Context) {
	var reqs []struct {
		ID     uint64  `json:"id"`
		Name   string  `json:"name"`
		Tag    string  `json:"tag"`
		Status string  `json:"status"`
		LockNo string  `json:"lock_no"`
		Price  float32 `json:"price"`
		Image  string  `json:"image"`
	}
	if err := c.ShouldBindJSON(&reqs); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if err := logic.BatchUpdateRoom(reqs); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "批量更新成功", nil)
}

// ==================== 设备管理 ====================

func AdminGetDeviceList(c *gin.Context) {
	merchsIDStr := c.Query("merchs_id")
	groupsIDStr := c.Query("groups_id")
	name := c.Query("name")
	status := c.Query("status")
	deviceType := c.Query("type")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	merchsID := int32(0)
	groupsID := int32(0)
	page := 1
	pageSize := 10

	if merchsIDStr != "" {
		if v, err := strconv.Atoi(merchsIDStr); err == nil {
			merchsID = int32(v)
		}
	}

	if groupsIDStr != "" {
		if v, err := strconv.Atoi(groupsIDStr); err == nil {
			groupsID = int32(v)
		}
	}

	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			page = v
		}
	}

	if pageSizeStr != "" {
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	devices, total, err := logic.GetDeviceListFiltered(merchsID, groupsID, name, status, deviceType, page, pageSize)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  devices,
		"total": total,
	})
}

func AdminGetDeviceDetail(c *gin.Context) {
	id := c.Param("id")
	device, err := logic.GetDeviceByID(id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "获取成功", device)
}

func AdminCreateDevice(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		Code      string `json:"code" binding:"required"`
		MerchsID  int32  `json:"merchs_id" binding:"required"`
		GroupsID  int32  `json:"groups_id"`
		Type      string `json:"type"`
		LockCount int32  `json:"lock_count"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	deviceType := req.Type
	if deviceType == "" {
		deviceType = "集控"
	}

	device := &model.Device{
		Name:      req.Name,
		Code:      req.Code,
		MerchsID:  req.MerchsID,
		GroupsID:  req.GroupsID,
		Type:      deviceType,
		LockCount: req.LockCount,
		Status:    "在线",
	}

	if err := logic.CreateDevice(device); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", device)
}

func AdminUpdateDevice(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name      string `json:"name"`
		Code      string `json:"code"`
		Status    string `json:"status"`
		Type      string `json:"type"`
		LockCount int32  `json:"lock_count"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	device, err := logic.GetDeviceByID(id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}

	if req.Name != "" {
		device.Name = req.Name
	}
	if req.Code != "" {
		device.Code = req.Code
	}
	if req.Status != "" {
		device.Status = req.Status
	}
	if req.Type != "" {
		device.Type = req.Type
	}
	if req.LockCount >= 0 {
		device.LockCount = req.LockCount
	}

	if err := logic.UpdateDevice(device); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", device)
}

func AdminDeleteDevice(c *gin.Context) {
	id := c.Param("id")
	if err := logic.DeleteDevice(id); err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteDevice(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要删除的设备")
		return
	}
	if err := logic.BatchDeleteDevice(req.IDs); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "批量删除成功", nil)
}

func AdminBatchUpdateDevice(c *gin.Context) {
	var reqs []struct {
		ID        uint64 `json:"id"`
		Name      string `json:"name"`
		Code      string `json:"code"`
		Status    string `json:"status"`
		Type      string `json:"type"`
		LockCount int32  `json:"lock_count"`
	}
	if err := c.ShouldBindJSON(&reqs); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if err := logic.BatchUpdateDevice(reqs); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "批量更新成功", nil)
}

// ==================== 分组管理 ====================

func AdminGetGroupList(c *gin.Context) {
	merchsIDStr := c.Query("merchs_id")
	name := c.Query("name")
	groupType := c.Query("type")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	merchsID := int32(0)
	page := 1
	pageSize := 10

	if merchsIDStr != "" {
		if v, err := strconv.Atoi(merchsIDStr); err == nil {
			merchsID = int32(v)
		}
	}

	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			page = v
		}
	}

	if pageSizeStr != "" {
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	groups, total, err := logic.GetGroupListFiltered(merchsID, name, groupType, page, pageSize)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  groups,
		"total": total,
	})
}

func AdminGetGroupDetail(c *gin.Context) {
	id := c.Param("id")
	group, err := logic.GetGroupByID(id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "获取成功", group)
}

func AdminCreateGroup(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		MerchsID    int32  `json:"merchs_id" binding:"required"`
		RulesID     int32  `json:"rules_id"`
		Phone       string `json:"phone"`
		Type        string `json:"type"`
		Location    string `json:"location"`
		BindNumber  string `json:"bind_number"`
		ConsumePush string `json:"consume_push"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	groupType := req.Type
	if groupType == "" {
		groupType = "存柜"
	}

	bindNumber := req.BindNumber
	if bindNumber == "" {
		bindNumber = "关闭"
	}

	consumePush := req.ConsumePush
	if consumePush == "" {
		consumePush = "关闭"
	}

	group := &model.Group{
		Name:        req.Name,
		MerchsID:    req.MerchsID,
		RulesID:     req.RulesID,
		Phone:       req.Phone,
		Type:        groupType,
		Location:    req.Location,
		BindNumber:  bindNumber,
		ConsumePush: consumePush,
	}

	if err := logic.CreateGroup(group); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", group)
}

func AdminUpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name        string `json:"name"`
		RulesID     int32  `json:"rules_id"`
		Phone       string `json:"phone"`
		Type        string `json:"type"`
		Location    string `json:"location"`
		BindNumber  string `json:"bind_number"`
		ConsumePush string `json:"consume_push"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	group, err := logic.GetGroupByID(id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}

	if req.Name != "" {
		group.Name = req.Name
	}
	if req.RulesID > 0 {
		group.RulesID = req.RulesID
	}
	if req.Phone != "" {
		group.Phone = req.Phone
	}
	if req.Type != "" {
		group.Type = req.Type
	}
	if req.Location != "" {
		group.Location = req.Location
	}
	if req.BindNumber != "" {
		group.BindNumber = req.BindNumber
	}
	if req.ConsumePush != "" {
		group.ConsumePush = req.ConsumePush
	}

	if err := logic.UpdateGroup(group); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", group)
}

func AdminDeleteGroup(c *gin.Context) {
	id := c.Param("id")
	if err := logic.DeleteGroup(id); err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteGroup(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要删除的分组")
		return
	}
	if err := logic.BatchDeleteGroup(req.IDs); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "批量删除成功", nil)
}

func AdminBatchUpdateGroup(c *gin.Context) {
	var reqs []struct {
		ID          uint64 `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Location    string `json:"location"`
		Phone       string `json:"phone"`
		BindNumber  string `json:"bind_number"`
		ConsumePush string `json:"consume_push"`
	}
	if err := c.ShouldBindJSON(&reqs); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if err := logic.BatchUpdateGroup(reqs); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "批量更新成功", nil)
}

// ==================== 规则管理 ====================

func AdminGetRuleList(c *gin.Context) {
	merchsIDStr := c.Query("merchs_id")
	name := c.Query("name")
	ruleType := c.Query("type")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	merchsID := int32(0)
	page := 1
	pageSize := 10

	if merchsIDStr != "" {
		if v, err := strconv.Atoi(merchsIDStr); err == nil {
			merchsID = int32(v)
		}
	}

	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			page = v
		}
	}

	if pageSizeStr != "" {
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	rules, total, err := logic.GetRuleListFiltered(merchsID, name, ruleType, page, pageSize)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  rules,
		"total": total,
	})
}

func AdminGetRuleDetail(c *gin.Context) {
	id := c.Param("id")
	rule, err := logic.GetRuleByID(id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "获取成功", rule)
}

func AdminCreateRule(c *gin.Context) {
	var req struct {
		Name         string  `json:"name" binding:"required"`
		MerchsID     int32   `json:"merchs_id" binding:"required"`
		Type         string  `json:"type"`
		Mode         string  `json:"mode"`
		Price        float32 `json:"price"`
		Deposit      float32 `json:"deposit"`
		DurationUnit string  `json:"duration_unit"`
		AutoEndTime  int32   `json:"auto_end_time"`
		Description  string  `json:"description"`
		FreeTime     int32   `json:"free_time"`
		AutoRefund   bool    `json:"auto_refund"`
		ManualRenew  bool    `json:"manual_renew"`
		TimeOptions  string  `json:"time_options"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	ruleType := req.Type
	if ruleType == "" {
		ruleType = "free"
	}

	mode := req.Mode
	if mode == "" {
		mode = "single"
	}

	durationUnit := req.DurationUnit
	if durationUnit == "" {
		durationUnit = "hour"
	}

	rule := &model.Rule{
		Name:         req.Name,
		MerchsID:     req.MerchsID,
		Type:         ruleType,
		Mode:         mode,
		Price:        req.Price,
		Deposit:      req.Deposit,
		DurationUnit: durationUnit,
		AutoEndTime:  req.AutoEndTime,
		Description:  req.Description,
		FreeTime:     req.FreeTime,
		AutoRefund:   req.AutoRefund,
		ManualRenew:  req.ManualRenew,
		TimeOptions:  req.TimeOptions,
	}

	if err := logic.CreateRule(rule); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", rule)
}

func AdminUpdateRule(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name         string  `json:"name"`
		Type         string  `json:"type"`
		Mode         string  `json:"mode"`
		Price        float32 `json:"price"`
		Deposit      float32 `json:"deposit"`
		DurationUnit string  `json:"duration_unit"`
		AutoEndTime  int32   `json:"auto_end_time"`
		Description  string  `json:"description"`
		FreeTime     int32   `json:"free_time"`
		AutoRefund   bool    `json:"auto_refund"`
		ManualRenew  bool    `json:"manual_renew"`
		TimeOptions  string  `json:"time_options"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	rule, err := logic.GetRuleByID(id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}

	if req.Name != "" {
		rule.Name = req.Name
	}
	if req.Type != "" {
		rule.Type = req.Type
	}
	if req.Mode != "" {
		rule.Mode = req.Mode
	}
	if req.Price >= 0 {
		rule.Price = req.Price
	}
	if req.Deposit >= 0 {
		rule.Deposit = req.Deposit
	}
	if req.DurationUnit != "" {
		rule.DurationUnit = req.DurationUnit
	}
	if req.AutoEndTime >= 0 {
		rule.AutoEndTime = req.AutoEndTime
	}
	if req.Description != "" {
		rule.Description = req.Description
	}
	if req.FreeTime >= 0 {
		rule.FreeTime = req.FreeTime
	}
	rule.AutoRefund = req.AutoRefund
	rule.ManualRenew = req.ManualRenew
	if req.TimeOptions != "" {
		rule.TimeOptions = req.TimeOptions
	}

	if err := logic.UpdateRule(rule); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", rule)
}

func AdminDeleteRule(c *gin.Context) {
	id := c.Param("id")
	if err := logic.DeleteRule(id); err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteRule(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要删除的规则")
		return
	}
	if err := logic.BatchDeleteRule(req.IDs); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "批量删除成功", nil)
}

func AdminBatchUpdateRule(c *gin.Context) {
	var reqs []struct {
		ID          uint32  `json:"id"`
		Name        string  `json:"name"`
		Type        string  `json:"type"`
		Mode        string  `json:"mode"`
		Price       float32 `json:"price"`
		Deposit     float32 `json:"deposit"`
		Description string  `json:"description"`
	}
	if err := c.ShouldBindJSON(&reqs); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if err := logic.BatchUpdateRule(reqs); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "批量更新成功", nil)
}

// ==================== 订单管理 ====================

func AdminGetOrderList(c *gin.Context) {
	merchsIDStr := c.Query("merchs_id")
	usersIDStr := c.Query("users_id")
	roomsIDStr := c.Query("rooms_id")
	status := c.Query("status")
	orderCode := c.Query("code")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	merchsID := int32(0)
	usersID := int32(0)
	roomsID := int32(0)
	page := 1
	pageSize := 10

	if merchsIDStr != "" {
		if v, err := strconv.Atoi(merchsIDStr); err == nil {
			merchsID = int32(v)
		}
	}

	if usersIDStr != "" {
		if v, err := strconv.Atoi(usersIDStr); err == nil {
			usersID = int32(v)
		}
	}

	if roomsIDStr != "" {
		if v, err := strconv.Atoi(roomsIDStr); err == nil {
			roomsID = int32(v)
		}
	}

	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			page = v
		}
	}

	if pageSizeStr != "" {
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	orders, total, err := logic.GetOrderListFiltered(merchsID, usersID, roomsID, status, orderCode, page, pageSize)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  orders,
		"total": total,
	})
}

func AdminGetOrderDetail(c *gin.Context) {
	id := c.Param("id")
	order, err := logic.GetOrderByID(id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "获取成功", order)
}

func AdminUpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status"`
		Remark string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	order, err := logic.GetOrderByID(id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}

	if req.Status != "" {
		order.Status = req.Status
	}
	if req.Remark != "" {
		order.Remark = req.Remark
	}

	if err := logic.UpdateOrder(order); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", order)
}

func AdminDeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := logic.DeleteOrder(id); err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteOrder(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要删除的订单")
		return
	}
	if err := logic.BatchDeleteOrder(req.IDs); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "批量删除成功", nil)
}

func AdminBatchUpdateOrder(c *gin.Context) {
	var reqs []struct {
		ID     uint64 `json:"id"`
		Status string `json:"status"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&reqs); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if err := logic.BatchUpdateOrder(reqs); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "批量更新成功", nil)
}

// ==================== 商家管理 ====================

func AdminGetMerchList(c *gin.Context) {
	account := c.Query("account")
	phone := c.Query("phone")
	role := c.Query("role")
	status := c.Query("status")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	page := 1
	pageSize := 10

	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			page = v
		}
	}

	if pageSizeStr != "" {
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	merchs, total, err := logic.GetMerchListFiltered(account, phone, role, status, page, pageSize)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  merchs,
		"total": total,
	})
}

func AdminGetMerchDetail(c *gin.Context) {
	id := c.Param("id")
	merch, err := logic.GetMerchByID(id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "获取成功", merch)
}

func AdminCreateMerch(c *gin.Context) {
	var req struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Role     string `json:"role"`
		Status   string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	role := req.Role
	if role == "" {
		role = "商家"
	}

	status := req.Status
	if status == "" {
		status = "0"
	}

	merch := &model.Merch{
		Account:  req.Account,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     role,
		Status:   status,
	}

	if err := logic.CreateMerch(merch); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", merch)
}

func AdminUpdateMerch(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Account  string `json:"account"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Role     string `json:"role"`
		Status   string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	merch, err := logic.GetMerchByID(id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}

	if req.Account != "" {
		merch.Account = req.Account
	}
	if req.Password != "" {
		merch.Password = req.Password
	}
	if req.Email != "" {
		merch.Email = req.Email
	}
	if req.Phone != "" {
		merch.Phone = req.Phone
	}
	if req.Role != "" {
		merch.Role = req.Role
	}
	if req.Status != "" {
		merch.Status = req.Status
	}

	if err := logic.UpdateMerch(merch); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", merch)
}

func AdminDeleteMerch(c *gin.Context) {
	id := c.Param("id")
	if err := logic.DeleteMerch(id); err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteMerch(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要删除的商家")
		return
	}
	if err := logic.BatchDeleteMerch(req.IDs); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "批量删除成功", nil)
}

func AdminBatchUpdateMerch(c *gin.Context) {
	var reqs []struct {
		ID      uint64 `json:"id"`
		Account string `json:"account"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Role    string `json:"role"`
		Status  string `json:"status"`
	}
	if err := c.ShouldBindJSON(&reqs); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if err := logic.BatchUpdateMerch(reqs); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "批量更新成功", nil)
}

// ==================== 设备日志管理 ====================

func AdminGetDeviceLogList(c *gin.Context) {
	merchsIDStr := c.Query("merchs_id")
	devicesIDStr := c.Query("devices_id")
	roomIDStr := c.Query("room_id")
	code := c.Query("code")
	deviceType := c.Query("type")
	control := c.Query("control")
	status := c.Query("status")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	merchsID := int32(0)
	devicesID := int32(0)
	roomID := int32(0)
	page := 1
	pageSize := 10

	if merchsIDStr != "" {
		if v, err := strconv.Atoi(merchsIDStr); err == nil {
			merchsID = int32(v)
		}
	}

	if devicesIDStr != "" {
		if v, err := strconv.Atoi(devicesIDStr); err == nil {
			devicesID = int32(v)
		}
	}

	if roomIDStr != "" {
		if v, err := strconv.Atoi(roomIDStr); err == nil {
			roomID = int32(v)
		}
	}

	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			page = v
		}
	}

	if pageSizeStr != "" {
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	logs, total, err := logic.GetDeviceLogListFiltered(merchsID, devicesID, roomID, code, deviceType, control, status, page, pageSize)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  logs,
		"total": total,
	})
}

func AdminGetDeviceLogDetail(c *gin.Context) {
	id := c.Param("id")
	log, err := logic.GetDeviceLogByID(id)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "获取成功", log)
}

func AdminDeleteDeviceLog(c *gin.Context) {
	id := c.Param("id")
	if err := logic.DeleteDeviceLog(id); err != nil {
		response.Fail(c, 404, err.Error())
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteDeviceLog(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要删除的日志")
		return
	}
	if err := logic.BatchDeleteDeviceLog(req.IDs); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "批量删除成功", nil)
}

func AdminGetDashboardStats(c *gin.Context) {
	roomCount, err := logic.GetRoomCount()
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	deviceCount, err := logic.GetDeviceCount()
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	groupCount, err := logic.GetGroupCount()
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	orderCount, err := logic.GetOrderCount()
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"roomCount":   roomCount,
		"deviceCount": deviceCount,
		"groupCount":  groupCount,
		"orderCount":  orderCount,
	})
}
