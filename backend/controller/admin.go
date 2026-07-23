package controller

import (
	"centraliz-backend/logic"
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
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
	boardNo := c.Query("board_no")
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
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
			pageSize = v
		}
	}

	rooms, total, err := logic.GetRoomListFiltered(merchsID, groupsID, name, status, boardNo, page, pageSize)
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
	var req struct {
		IDs  []string               `json:"ids"`
		Data map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要更新的房间")
		return
	}
	if err := logic.BatchUpdateRoomByIDs(req.IDs, req.Data); err != nil {
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
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
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
	var req struct {
		IDs  []string               `json:"ids"`
		Data map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要更新的设备")
		return
	}
	if err := logic.BatchUpdateDeviceByIDs(req.IDs, req.Data); err != nil {
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
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
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
	var req struct {
		IDs  []string               `json:"ids"`
		Data map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要更新的分组")
		return
	}
	if err := logic.BatchUpdateGroupByIDs(req.IDs, req.Data); err != nil {
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
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
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
	var req struct {
		IDs  []string               `json:"ids"`
		Data map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要更新的规则")
		return
	}
	if err := logic.BatchUpdateRuleByIDs(req.IDs, req.Data); err != nil {
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
	orderNo := c.Query("order_no")
	userPhone := c.Query("user_phone")
	payType := c.Query("pay_type")
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
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
			pageSize = v
		}
	}

	orders, total, err := logic.GetOrderListFiltered(merchsID, usersID, roomsID, status, orderCode, orderNo, userPhone, payType, page, pageSize)
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
	var req struct {
		IDs  []string               `json:"ids"`
		Data map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要更新的订单")
		return
	}
	if err := logic.BatchUpdateOrderByIDs(req.IDs, req.Data); err != nil {
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
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
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
	var req struct {
		IDs  []string               `json:"ids"`
		Data map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要更新的商家")
		return
	}
	if err := logic.BatchUpdateMerchByIDs(req.IDs, req.Data); err != nil {
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
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
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

// ==================== 汇付账号管理 ====================

func AdminGetHuifuAccountList(c *gin.Context) {
	merchsIDStr := c.Query("merchs_id")
	code := c.Query("code")
	account := c.Query("account")
	phone := c.Query("phone")
	name := c.Query("name")
	typeStr := c.Query("type")
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
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
			pageSize = v
		}
	}

	query := db.DB.Model(&model.HuifuAccount{})
	if merchsID > 0 {
		query = query.Where("merchs_id = ?", merchsID)
	}
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if account != "" {
		query = query.Where("account LIKE ?", "%"+account+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if typeStr != "" {
		query = query.Where("type = ?", typeStr)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	var accounts []model.HuifuAccount
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&accounts).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  accounts,
		"total": total,
	})
}

func AdminGetHuifuAccountDetail(c *gin.Context) {
	id := c.Param("id")
	var account model.HuifuAccount
	if err := db.DB.Where("id = ?", id).First(&account).Error; err != nil {
		response.Fail(c, 404, "汇付账号不存在")
		return
	}
	response.SuccessWithMsg(c, "获取成功", account)
}

func AdminCreateHuifuAccount(c *gin.Context) {
	var req struct {
		MerchsID  int32   `json:"merchsId" binding:"required"`
		Type      string  `json:"type" binding:"required"`
		Account   string  `json:"account" binding:"required"`
		Name      string  `json:"name" binding:"required"`
		Phone     string  `json:"phone" binding:"required"`
		Identity  string  `json:"identity" binding:"required"`
		Card      string  `json:"card" binding:"required"`
		Storename string  `json:"storename"`
		Encrypt   string  `json:"encrypt"`
		Area      string  `json:"area"`
		Picture   string  `json:"picture"`
		Remarks   string  `json:"remarks"`
		Sharing   string  `json:"sharing"`
		Share     string  `json:"share"`
		Rate      float64 `json:"rate"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	account := &model.HuifuAccount{
		MerchsID:  req.MerchsID,
		Type:      req.Type,
		Account:   req.Account,
		Name:      req.Name,
		Phone:     req.Phone,
		Identity:  req.Identity,
		Card:      req.Card,
		Storename: req.Storename,
		Encrypt:   req.Encrypt,
		Area:      req.Area,
		Picture:   req.Picture,
		Remarks:   req.Remarks,
		Sharing:   req.Sharing,
		Share:     req.Share,
		Rate:      req.Rate,
		Choose:    "0",
	}

	if err := db.DB.Create(account).Error; err != nil {
		response.Fail(c, 500, "创建失败")
		return
	}

	response.SuccessWithMsg(c, "创建成功", account)
}

func AdminUpdateHuifuAccount(c *gin.Context) {
	id := c.Param("id")
	var account model.HuifuAccount
	if err := db.DB.Where("id = ?", id).First(&account).Error; err != nil {
		response.Fail(c, 404, "汇付账号不存在")
		return
	}

	var req struct {
		Type      string  `json:"type"`
		Account   string  `json:"account"`
		Name      string  `json:"name"`
		Phone     string  `json:"phone"`
		Identity  string  `json:"identity"`
		Card      string  `json:"card"`
		Storename string  `json:"storename"`
		Encrypt   string  `json:"encrypt"`
		Area      string  `json:"area"`
		Picture   string  `json:"picture"`
		Remarks   string  `json:"remarks"`
		Sharing   string  `json:"sharing"`
		Share     string  `json:"share"`
		Rate      float64 `json:"rate"`
		Choose    string  `json:"choose"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	if req.Type != "" {
		account.Type = req.Type
	}
	if req.Account != "" {
		account.Account = req.Account
	}
	if req.Name != "" {
		account.Name = req.Name
	}
	if req.Phone != "" {
		account.Phone = req.Phone
	}
	if req.Identity != "" {
		account.Identity = req.Identity
	}
	if req.Card != "" {
		account.Card = req.Card
	}
	if req.Storename != "" {
		account.Storename = req.Storename
	}
	if req.Encrypt != "" {
		account.Encrypt = req.Encrypt
	}
	if req.Area != "" {
		account.Area = req.Area
	}
	if req.Picture != "" {
		account.Picture = req.Picture
	}
	if req.Remarks != "" {
		account.Remarks = req.Remarks
	}
	if req.Sharing != "" {
		account.Sharing = req.Sharing
	}
	if req.Share != "" {
		account.Share = req.Share
	}
	if req.Rate >= 0 {
		account.Rate = req.Rate
	}
	if req.Choose != "" {
		account.Choose = req.Choose
	}

	if err := db.DB.Save(&account).Error; err != nil {
		response.Fail(c, 500, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", account)
}

func AdminDeleteHuifuAccount(c *gin.Context) {
	id := c.Param("id")
	var account model.HuifuAccount
	if err := db.DB.Where("id = ?", id).First(&account).Error; err != nil {
		response.Fail(c, 404, "汇付账号不存在")
		return
	}
	if err := db.DB.Delete(&account).Error; err != nil {
		response.Fail(c, 500, "删除失败")
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteHuifuAccount(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要删除的汇付账号")
		return
	}
	if err := db.DB.Delete(&model.HuifuAccount{}, "id IN (?)", req.IDs).Error; err != nil {
		response.Fail(c, 500, "批量删除失败")
		return
	}
	response.SuccessWithMsg(c, "批量删除成功", nil)
}

func AdminBatchUpdateHuifuAccount(c *gin.Context) {
	var req struct {
		IDs  []string               `json:"ids"`
		Data map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要更新的汇付账号")
		return
	}
	if err := db.DB.Model(&model.HuifuAccount{}).Where("id IN (?)", req.IDs).Updates(req.Data).Error; err != nil {
		response.Fail(c, 500, "批量更新失败")
		return
	}
	response.SuccessWithMsg(c, "批量更新成功", nil)
}

// ==================== 商户支付管理 ====================

func AdminGetMerchPayList(c *gin.Context) {
	merchsIDStr := c.Query("merchs_id")
	code := c.Query("code")
	name := c.Query("name")
	typeStr := c.Query("type")
	status := c.Query("status")
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
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
			pageSize = v
		}
	}

	query := db.DB.Model(&model.MerchPay{})
	if merchsID > 0 {
		query = query.Where("merchs_id = ?", merchsID)
	}
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if typeStr != "" {
		query = query.Where("type = ?", typeStr)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	var orders []model.MerchPay
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  orders,
		"total": total,
	})
}

func AdminGetMerchPayDetail(c *gin.Context) {
	id := c.Param("id")
	var order model.MerchPay
	if err := db.DB.Where("id = ?", id).First(&order).Error; err != nil {
		response.Fail(c, 404, "订单不存在")
		return
	}
	response.SuccessWithMsg(c, "获取成功", order)
}

func AdminDeleteMerchPay(c *gin.Context) {
	id := c.Param("id")
	var order model.MerchPay
	if err := db.DB.Where("id = ?", id).First(&order).Error; err != nil {
		response.Fail(c, 404, "订单不存在")
		return
	}
	if err := db.DB.Delete(&order).Error; err != nil {
		response.Fail(c, 500, "删除失败")
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteMerchPay(c *gin.Context) {
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
	if err := db.DB.Delete(&model.MerchPay{}, "id IN (?)", req.IDs).Error; err != nil {
		response.Fail(c, 500, "批量删除失败")
		return
	}
	response.SuccessWithMsg(c, "批量删除成功", nil)
}

// ==================== 房间图片管理 ====================

func AdminGetRoomImageList(c *gin.Context) {
	name := c.Query("name")
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
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
			pageSize = v
		}
	}

	query := db.DB.Model(&model.RoomImage{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	var images []model.RoomImage
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&images).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  images,
		"total": total,
	})
}

func AdminGetRoomImageDetail(c *gin.Context) {
	id := c.Param("id")
	var image model.RoomImage
	if err := db.DB.Where("id = ?", id).First(&image).Error; err != nil {
		response.Fail(c, 404, "房间图片不存在")
		return
	}
	response.SuccessWithMsg(c, "获取成功", image)
}

func AdminCreateRoomImage(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Image string `json:"image" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	image := &model.RoomImage{
		Name:  req.Name,
		Image: req.Image,
	}

	if err := db.DB.Create(image).Error; err != nil {
		response.Fail(c, 500, "创建失败")
		return
	}

	response.SuccessWithMsg(c, "创建成功", image)
}

func AdminUpdateRoomImage(c *gin.Context) {
	id := c.Param("id")
	var image model.RoomImage
	if err := db.DB.Where("id = ?", id).First(&image).Error; err != nil {
		response.Fail(c, 404, "房间图片不存在")
		return
	}

	var req struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	if req.Name != "" {
		image.Name = req.Name
	}
	if req.Image != "" {
		image.Image = req.Image
	}

	if err := db.DB.Save(&image).Error; err != nil {
		response.Fail(c, 500, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", image)
}

func AdminDeleteRoomImage(c *gin.Context) {
	id := c.Param("id")
	var image model.RoomImage
	if err := db.DB.Where("id = ?", id).First(&image).Error; err != nil {
		response.Fail(c, 404, "房间图片不存在")
		return
	}
	if err := db.DB.Delete(&image).Error; err != nil {
		response.Fail(c, 500, "删除失败")
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteRoomImage(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要删除的房间图片")
		return
	}
	if err := db.DB.Delete(&model.RoomImage{}, "id IN (?)", req.IDs).Error; err != nil {
		response.Fail(c, 500, "批量删除失败")
		return
	}
	response.SuccessWithMsg(c, "批量删除成功", nil)
}

// ==================== 房间标签管理 ====================

func AdminGetRoomTagList(c *gin.Context) {
	merchsIDStr := c.Query("merchs_id")
	name := c.Query("name")
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
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
			pageSize = v
		}
	}

	query := db.DB.Model(&model.RoomTag{})
	if merchsID > 0 {
		query = query.Where("merchs_id = ?", merchsID)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	var tags []model.RoomTag
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&tags).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  tags,
		"total": total,
	})
}

func AdminGetRoomTagDetail(c *gin.Context) {
	id := c.Param("id")
	var tag model.RoomTag
	if err := db.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		response.Fail(c, 404, "房间标签不存在")
		return
	}
	response.SuccessWithMsg(c, "获取成功", tag)
}

func AdminCreateRoomTag(c *gin.Context) {
	var req struct {
		MerchsID int32  `json:"merchsId"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	tag := &model.RoomTag{
		MerchsID: req.MerchsID,
		Name:     req.Name,
	}

	if err := db.DB.Create(tag).Error; err != nil {
		response.Fail(c, 500, "创建失败")
		return
	}

	response.SuccessWithMsg(c, "创建成功", tag)
}

func AdminUpdateRoomTag(c *gin.Context) {
	id := c.Param("id")
	var tag model.RoomTag
	if err := db.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		response.Fail(c, 404, "房间标签不存在")
		return
	}

	var req struct {
		MerchsID int32  `json:"merchsId"`
		Name     string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	if req.MerchsID > 0 {
		tag.MerchsID = req.MerchsID
	}
	if req.Name != "" {
		tag.Name = req.Name
	}

	if err := db.DB.Save(&tag).Error; err != nil {
		response.Fail(c, 500, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", tag)
}

func AdminDeleteRoomTag(c *gin.Context) {
	id := c.Param("id")
	var tag model.RoomTag
	if err := db.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		response.Fail(c, 404, "房间标签不存在")
		return
	}
	if err := db.DB.Delete(&tag).Error; err != nil {
		response.Fail(c, 500, "删除失败")
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteRoomTag(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要删除的房间标签")
		return
	}
	if err := db.DB.Delete(&model.RoomTag{}, "id IN (?)", req.IDs).Error; err != nil {
		response.Fail(c, 500, "批量删除失败")
		return
	}
	response.SuccessWithMsg(c, "批量删除成功", nil)
}

// ==================== 子商户管理 ====================

func AdminGetSubMerchList(c *gin.Context) {
	merchsIDStr := c.Query("merchs_id")
	account := c.Query("account")
	phone := c.Query("phone")
	role := c.Query("role")
	status := c.Query("status")
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
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
			pageSize = v
		}
	}

	query := db.DB.Model(&model.SubMerch{})
	if merchsID > 0 {
		query = query.Where("merchs_id = ?", merchsID)
	}
	if account != "" {
		query = query.Where("account LIKE ?", "%"+account+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	var submerchs []model.SubMerch
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&submerchs).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  submerchs,
		"total": total,
	})
}

func AdminGetSubMerchDetail(c *gin.Context) {
	id := c.Param("id")
	var submerch model.SubMerch
	if err := db.DB.Where("id = ?", id).First(&submerch).Error; err != nil {
		response.Fail(c, 404, "子商户不存在")
		return
	}
	response.SuccessWithMsg(c, "获取成功", submerch)
}

func AdminCreateSubMerch(c *gin.Context) {
	var req struct {
		MerchsID int32  `json:"merchsId"`
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Role     string `json:"role"`
		Status   string `json:"status"`
		Rule     string `json:"rule"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	role := req.Role
	if role == "" {
		role = "0"
	}

	status := req.Status
	if status == "" {
		status = "0"
	}

	submerch := &model.SubMerch{
		MerchsID: req.MerchsID,
		Account:  req.Account,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     role,
		Status:   status,
		Rule:     req.Rule,
	}

	if err := db.DB.Create(submerch).Error; err != nil {
		response.Fail(c, 500, "创建失败")
		return
	}

	response.SuccessWithMsg(c, "创建成功", submerch)
}

func AdminUpdateSubMerch(c *gin.Context) {
	id := c.Param("id")
	var submerch model.SubMerch
	if err := db.DB.Where("id = ?", id).First(&submerch).Error; err != nil {
		response.Fail(c, 404, "子商户不存在")
		return
	}

	var req struct {
		MerchsID int32  `json:"merchsId"`
		Account  string `json:"account"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Role     string `json:"role"`
		Status   string `json:"status"`
		Rule     string `json:"rule"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	if req.MerchsID > 0 {
		submerch.MerchsID = req.MerchsID
	}
	if req.Account != "" {
		submerch.Account = req.Account
	}
	if req.Password != "" {
		submerch.Password = req.Password
	}
	if req.Email != "" {
		submerch.Email = req.Email
	}
	if req.Phone != "" {
		submerch.Phone = req.Phone
	}
	if req.Role != "" {
		submerch.Role = req.Role
	}
	if req.Status != "" {
		submerch.Status = req.Status
	}
	if req.Rule != "" {
		submerch.Rule = req.Rule
	}

	if err := db.DB.Save(&submerch).Error; err != nil {
		response.Fail(c, 500, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", submerch)
}

func AdminDeleteSubMerch(c *gin.Context) {
	id := c.Param("id")
	var submerch model.SubMerch
	if err := db.DB.Where("id = ?", id).First(&submerch).Error; err != nil {
		response.Fail(c, 404, "子商户不存在")
		return
	}
	if err := db.DB.Delete(&submerch).Error; err != nil {
		response.Fail(c, 500, "删除失败")
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteSubMerch(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要删除的子商户")
		return
	}
	if err := db.DB.Delete(&model.SubMerch{}, "id IN (?)", req.IDs).Error; err != nil {
		response.Fail(c, 500, "批量删除失败")
		return
	}
	response.SuccessWithMsg(c, "批量删除成功", nil)
}

func AdminBatchUpdateSubMerch(c *gin.Context) {
	var req struct {
		IDs  []string               `json:"ids"`
		Data map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要更新的子商户")
		return
	}
	if err := db.DB.Model(&model.SubMerch{}).Where("id IN (?)", req.IDs).Updates(req.Data).Error; err != nil {
		response.Fail(c, 500, "批量更新失败")
		return
	}
	response.SuccessWithMsg(c, "批量更新成功", nil)
}

// ==================== 用户管理 ====================

func AdminGetUserList(c *gin.Context) {
	merchsIDStr := c.Query("merchs_id")
	roomsIDStr := c.Query("rooms_id")
	name := c.Query("name")
	account := c.Query("account")
	phone := c.Query("phone")
	privacy := c.Query("privacy")
	status := c.Query("status")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	merchsID := int32(0)
	roomsID := int32(0)
	page := 1
	pageSize := 10

	if merchsIDStr != "" {
		if v, err := strconv.Atoi(merchsIDStr); err == nil {
			merchsID = int32(v)
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
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
			pageSize = v
		}
	}

	query := db.DB.Model(&model.User{})
	if merchsID > 0 {
		query = query.Where("merchs_id = ?", merchsID)
	}
	if roomsID > 0 {
		query = query.Where("rooms_id = ?", roomsID)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if account != "" {
		query = query.Where("account LIKE ?", "%"+account+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if privacy != "" {
		query = query.Where("privacy = ?", privacy)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	var users []model.User
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  users,
		"total": total,
	})
}

func AdminGetUserDetail(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		response.Fail(c, 404, "用户不存在")
		return
	}
	response.SuccessWithMsg(c, "获取成功", user)
}

func AdminUpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		response.Fail(c, 404, "用户不存在")
		return
	}

	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Privacy string `json:"privacy"`
		Status  string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if req.Privacy != "" {
		user.Privacy = req.Privacy
	}
	if req.Status != "" {
		user.Status = &req.Status
	}

	if err := db.DB.Save(&user).Error; err != nil {
		response.Fail(c, 500, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", user)
}

func AdminDeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		response.Fail(c, 404, "用户不存在")
		return
	}
	if err := db.DB.Delete(&user).Error; err != nil {
		response.Fail(c, 500, "删除失败")
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteUser(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要删除的用户")
		return
	}
	if err := db.DB.Delete(&model.User{}, "id IN (?)", req.IDs).Error; err != nil {
		response.Fail(c, 500, "批量删除失败")
		return
	}
	response.SuccessWithMsg(c, "批量删除成功", nil)
}

func AdminBatchUpdateUser(c *gin.Context) {
	var req struct {
		IDs  []string               `json:"ids"`
		Data map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要更新的用户")
		return
	}
	if err := db.DB.Model(&model.User{}).Where("id IN (?)", req.IDs).Updates(req.Data).Error; err != nil {
		response.Fail(c, 500, "批量更新失败")
		return
	}
	response.SuccessWithMsg(c, "批量更新成功", nil)
}

// ==================== 微信用户管理 ====================

func AdminGetWxUserList(c *gin.Context) {
	merchsIDStr := c.Query("merchs_id")
	usersIDStr := c.Query("users_id")
	nickname := c.Query("nickname")
	openid := c.Query("openid")
	platform := c.Query("platform")
	status := c.Query("status")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	merchsID := int32(0)
	usersID := int32(0)
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

	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			page = v
		}
	}

	if pageSizeStr != "" {
		if v, err := strconv.Atoi(pageSizeStr); err == nil && v > 0 && v <= 200 {
			pageSize = v
		}
	}

	query := db.DB.Model(&model.WechatUser{})
	if merchsID > 0 {
		query = query.Where("merchs_id = ?", merchsID)
	}
	if usersID > 0 {
		query = query.Where("users_id = ?", usersID)
	}
	if nickname != "" {
		query = query.Where("nickname LIKE ?", "%"+nickname+"%")
	}
	if openid != "" {
		query = query.Where("openid LIKE ?", "%"+openid+"%")
	}
	if platform != "" {
		query = query.Where("platform = ?", platform)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	var wxusers []model.WechatUser
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&wxusers).Error; err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data":  wxusers,
		"total": total,
	})
}

func AdminGetWxUserDetail(c *gin.Context) {
	id := c.Param("id")
	var wxuser model.WechatUser
	if err := db.DB.Where("id = ?", id).First(&wxuser).Error; err != nil {
		response.Fail(c, 404, "微信用户不存在")
		return
	}
	response.SuccessWithMsg(c, "获取成功", wxuser)
}

func AdminUpdateWxUser(c *gin.Context) {
	id := c.Param("id")
	var wxuser model.WechatUser
	if err := db.DB.Where("id = ?", id).First(&wxuser).Error; err != nil {
		response.Fail(c, 404, "微信用户不存在")
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Gender   int    `json:"gender"`
		Country  string `json:"country"`
		Province string `json:"province"`
		City     string `json:"city"`
		Language string `json:"language"`
		Status   string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	if req.Nickname != "" {
		wxuser.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		wxuser.Avatar = req.Avatar
	}
	if req.Gender >= 0 {
		wxuser.Gender = req.Gender
	}
	if req.Country != "" {
		wxuser.Country = req.Country
	}
	if req.Province != "" {
		wxuser.Province = req.Province
	}
	if req.City != "" {
		wxuser.City = req.City
	}
	if req.Language != "" {
		wxuser.Language = req.Language
	}
	if req.Status != "" {
		wxuser.Status = req.Status
	}

	if err := db.DB.Save(&wxuser).Error; err != nil {
		response.Fail(c, 500, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", wxuser)
}

func AdminDeleteWxUser(c *gin.Context) {
	id := c.Param("id")
	var wxuser model.WechatUser
	if err := db.DB.Where("id = ?", id).First(&wxuser).Error; err != nil {
		response.Fail(c, 404, "微信用户不存在")
		return
	}
	if err := db.DB.Delete(&wxuser).Error; err != nil {
		response.Fail(c, 500, "删除失败")
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

func AdminBatchDeleteWxUser(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		response.Fail(c, 400, "请选择要删除的微信用户")
		return
	}
	if err := db.DB.Delete(&model.WechatUser{}, "id IN (?)", req.IDs).Error; err != nil {
		response.Fail(c, 500, "批量删除失败")
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

	todayOrders, err := logic.GetTodayOrderCount()
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	todayRevenue, err := logic.GetTodayRevenue()
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	totalUsers, err := logic.GetTotalUserCount()
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	totalRevenue, err := logic.GetTotalRevenue()
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"roomCount":    roomCount,
		"deviceCount":  deviceCount,
		"groupCount":   groupCount,
		"orderCount":   orderCount,
		"todayOrders":  todayOrders,
		"todayRevenue": todayRevenue,
		"totalUsers":   totalUsers,
		"totalRevenue": totalRevenue,
	})
}

func AdminGetTrendStats(c *gin.Context) {
	statType := c.Query("type")
	daysStr := c.Query("days")

	days := 7
	if daysStr != "" {
		if v, err := strconv.Atoi(daysStr); err == nil && v > 0 {
			days = v
		}
	}

	var data []map[string]interface{}
	var err error

	switch statType {
	case "order":
		data, err = logic.GetOrderTrend(days)
	case "revenue":
		data, err = logic.GetRevenueTrend(days)
	default:
		data, err = logic.GetOrderTrend(days)
	}

	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"data": data,
	})
}
