package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetDeviceLogList 获取设备日志列表
// 支持按设备编码、商家 ID、控制类型、状态、房间 ID、日期范围等条件过滤，支持分页
func GetDeviceLogList(c *gin.Context) {
	deviceCode := c.Query("code")
	merchsID := c.Query("merchs_id")
	controlType := c.Query("control")
	status := c.Query("status")
	roomID := c.Query("roomId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")

	// 设置默认分页参数，分页最多50条
	page := 1
	pageSize := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			// 分页最多50条
			if ps > 50 {
				ps = 50
			}
			pageSize = ps
		}
	}

	var logs []model.Devicelog
	var total int64
	query := db.DB.Model(&model.Devicelog{})

	if deviceCode != "" {
		query = query.Where("code = ?", deviceCode)
	}

	if merchsID != "" {
		if mid, err := strconv.Atoi(merchsID); err == nil && mid > 0 {
			query = query.Where("merchs_id = ?", mid)
		}
	}

	if controlType != "" {
		query = query.Where("control = ?", controlType)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if roomID != "" {
		if rid, err := strconv.Atoi(roomID); err == nil {
			query = query.Where("room_id = ?", rid)
		}
	}

	// 日期范围查询
	if startDate != "" {
		query = query.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(created_at) <= ?", endDate)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取日志总数失败: "+err.Error(), nil)
		return
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取日志列表失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":      logs,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
		"totalPage": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetDeviceLogDetail 获取设备日志详情
func GetDeviceLogDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的日志ID", nil)
		return
	}

	var log model.Devicelog
	if err := db.DB.First(&log, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "日志不存在", nil)
		return
	}

	control := "开锁"
	if log.Control != "" {
		control = log.Control
	}

	occupant := "用户"
	if log.Occupant != "" {
		occupant = log.Occupant
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"id":         log.ID,
		"code":       log.Code,
		"merchsId":   log.MerchsID,
		"devicesId":  log.DevicesID,
		"deviceName": log.DeviceName,
		"type":       log.Type,
		"control":    control,
		"status":     log.Status,
		"occupant":   occupant,
		"createdAt":  log.CreatedAt,
		"updatedAt":  log.UpdatedAt,
	})
}

// CreateDeviceLog 创建设备日志
func CreateDeviceLog(c *gin.Context) {
	var req struct {
		Code       string `json:"code"`
		MerchsID   int32  `json:"merchsId"`
		UsersID    int32  `json:"usersId"`
		DevicesID  int32  `json:"devicesId"`
		RoomID     int32  `json:"roomId"`
		DeviceName string `json:"deviceName"`
		Type       string `json:"type"`
		Control    string `json:"control"`
		Phone      string `json:"phone"`
		Status     string `json:"status"`
		Occupant   string `json:"occupant"`
		RoomName   string `json:"roomName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error(), nil)
		return
	}

	// 创建设备日志对象
	log := model.Devicelog{
		Code:       req.Code,
		MerchsID:   req.MerchsID,
		DeviceName: req.DeviceName,
		Type:       req.Type,
		Control:    req.Control,
		Phone:      req.Phone,
		Status:     req.Status,
		Occupant:   req.Occupant,
		RoomName:   req.RoomName,
	}

	if req.UsersID != 0 {
		log.UsersID = req.UsersID
	}

	if req.DevicesID != 0 {
		log.DevicesID = req.DevicesID
	}

	if req.RoomID != 0 {
		log.RoomID = req.RoomID
	}

	if err := db.DB.Create(&log).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建日志失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMsg(c, "创建成功", gin.H{"success": true})
}

// UpdateDeviceLog 更新设备日志
func UpdateDeviceLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的日志ID", nil)
		return
	}

	var req struct {
		Code       string `json:"code"`
		UsersID    int32  `json:"usersId"`
		DevicesID  int32  `json:"devicesId"`
		DeviceName string `json:"deviceName"`
		Type       string `json:"type"`
		Control    string `json:"control"`
		Status     string `json:"status"`
		Occupant   string `json:"occupant"`
		RoomID     int32  `json:"roomId"`
		RoomName   string `json:"roomName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error(), nil)
		return
	}

	var log model.Devicelog
	if err := db.DB.First(&log, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "日志不存在", nil)
		return
	}

	if req.Code != "" {
		log.Code = req.Code
	}

	if req.DeviceName != "" {
		log.DeviceName = req.DeviceName
	}

	if req.UsersID != 0 {
		log.UsersID = req.UsersID
	}

	if req.DevicesID != 0 {
		log.DevicesID = req.DevicesID
	}

	if req.Type != "" {
		log.Type = req.Type
	}

	if req.Control != "" {
		log.Control = req.Control
	}

	if req.Occupant != "" {
		log.Occupant = req.Occupant
	}

	if req.RoomID != 0 {
		log.RoomID = req.RoomID
	}

	if req.RoomName != "" {
		log.RoomName = req.RoomName
	}

	// 更新状态（允许设置为0）
	if req.Status != "" {
		log.Status = req.Status
	}

	if err := db.DB.Save(&log).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新日志失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMsg(c, "更新成功", gin.H{
		"id":        log.ID,
		"updatedAt": log.UpdatedAt,
	})
}

// DeleteDeviceLog 删除设备日志
func DeleteDeviceLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的日志ID", nil)
		return
	}

	var log model.Devicelog
	if err := db.DB.First(&log, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "日志不存在", nil)
		return
	}

	if err := db.DB.Delete(&log).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除日志失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMsg(c, "删除成功", gin.H{"success": true})
}
