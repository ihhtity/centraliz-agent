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
// 支持按设备编码、商家ID、控制类型、状态等条件过滤，支持分页
func GetDeviceLogList(c *gin.Context) {
	deviceCode := c.Query("code")
	merchsID := c.Query("merchsId")
	controlType := c.Query("control")
	status := c.Query("status")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")

	// 设置默认分页参数
	page := 1
	pageSize := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
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
		if st, err := strconv.Atoi(status); err == nil {
			query = query.Where("status = ?", st)
		}
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

	logList := make([]gin.H, len(logs))
	for i, log := range logs {
		control := "开锁"
		if log.Control != nil && *log.Control != "" {
			control = *log.Control
		}

		occupant := "用户"
		if log.Occupant != nil && *log.Occupant != "" {
			occupant = *log.Occupant
		}

		logList[i] = gin.H{
			"id":        log.ID,
			"code":      log.Code,
			"merchsId":  log.MerchsID,
			"devicesId": log.DevicesID,
			"name":      log.Name,
			"type":      log.Type,
			"control":   control,
			"status":    log.Status,
			"occupant":  occupant,
			"model":     log.Model,
			"createdAt": log.CreatedAt,
			"updatedAt": log.UpdatedAt,
		}
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":      logList,
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
	if log.Control != nil && *log.Control != "" {
		control = *log.Control
	}

	occupant := "用户"
	if log.Occupant != nil && *log.Occupant != "" {
		occupant = *log.Occupant
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"id":        log.ID,
		"code":      log.Code,
		"merchsId":  log.MerchsID,
		"devicesId": log.DevicesID,
		"name":      log.Name,
		"type":      log.Type,
		"control":   control,
		"status":    log.Status,
		"occupant":  occupant,
		"model":     log.Model,
		"createdAt": log.CreatedAt,
		"updatedAt": log.UpdatedAt,
	})
}

// CreateDeviceLog 创建设备日志
func CreateDeviceLog(c *gin.Context) {
	var req struct {
		Code      string `json:"code" binding:"required"`
		MerchsID  int32  `json:"merchsId" binding:"required"`
		DevicesID *int32 `json:"devicesId"`
		Name      string `json:"name" binding:"required"`
		Type      string `json:"type"`
		Control   string `json:"control"`
		Status    int8   `json:"status"`
		Occupant  string `json:"occupant"`
		Model     string `json:"model"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error(), nil)
		return
	}

	// 验证必填字段
	if req.Code == "" {
		response.Fail(c, http.StatusBadRequest, "设备编码不能为空", nil)
		return
	}

	if req.Name == "" {
		response.Fail(c, http.StatusBadRequest, "名称不能为空", nil)
		return
	}

	if req.MerchsID <= 0 {
		response.Fail(c, http.StatusBadRequest, "无效的商家ID", nil)
		return
	}

	// 创建设备日志对象
	log := model.Devicelog{
		Code:     req.Code,
		MerchsID: req.MerchsID,
		Name:     req.Name,
		Status:   req.Status,
	}

	if req.DevicesID != nil && *req.DevicesID > 0 {
		log.DevicesID = req.DevicesID
	}

	if req.Type != "" {
		log.Type = &req.Type
	}

	if req.Control != "" {
		log.Control = &req.Control
	}

	if req.Occupant != "" {
		log.Occupant = &req.Occupant
	}

	if req.Model != "" {
		log.Model = &req.Model
	}

	if err := db.DB.Create(&log).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建日志失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMsg(c, "创建成功", gin.H{
		"id":        log.ID,
		"code":      log.Code,
		"name":      log.Name,
		"createdAt": log.CreatedAt,
	})
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
		Code      string `json:"code"`
		DevicesID *int32 `json:"devicesId"`
		Name      string `json:"name"`
		Type      string `json:"type"`
		Control   string `json:"control"`
		Status    int8   `json:"status"`
		Occupant  string `json:"occupant"`
		Model     string `json:"model"`
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

	if req.Name != "" {
		log.Name = req.Name
	}

	if req.DevicesID != nil && *req.DevicesID > 0 {
		log.DevicesID = req.DevicesID
	}

	if req.Type != "" {
		log.Type = &req.Type
	}

	if req.Control != "" {
		log.Control = &req.Control
	}

	if req.Occupant != "" {
		log.Occupant = &req.Occupant
	}

	if req.Model != "" {
		log.Model = &req.Model
	}

	// 更新状态（允许设置为0）
	if req.Status >= 0 {
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
