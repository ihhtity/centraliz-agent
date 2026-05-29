package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetDeviceList 获取设备列表
// 支持通过merchs_id参数过滤特定商家的设备，支持分页
func GetDeviceList(c *gin.Context) {
	merchsID := c.Query("merchs_id")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	// 设置默认分页参数
	page := 1
	pageSize := 50

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

	var devices []model.Device
	var total int64
	query := db.DB.Model(&model.Device{})

	if merchsID != "" {
		query = query.Where("merchs_id = ?", merchsID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, 400, "获取设备总数失败", nil)
		return
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&devices).Error; err != nil {
		response.Fail(c, 400, "获取设备列表失败", nil)
		return
	}

	deviceList := make([]gin.H, len(devices))
	for i, device := range devices {
		// 获取状态（默认为"在线"）
		status := "在线"
		if device.Status != nil && *device.Status != "" {
			status = *device.Status
		}

		// 获取类型（默认为"集控"）
		deviceType := "集控"
		if device.Type != nil && *device.Type != "" {
			deviceType = *device.Type
		}

		deviceList[i] = gin.H{
			"id":        device.ID,
			"name":      device.Name,
			"code":      device.Code,
			"merchsId":  device.MerchsID,
			"roomsId":   device.RoomsID,
			"groupsId":  device.GroupsID,
			"status":    status,
			"type":      deviceType,
			"createdAt": device.CreatedAt,
			"updatedAt": device.UpdatedAt,
		}
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":      deviceList,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
		"totalPage": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetDeviceDetail 获取设备详情
func GetDeviceDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的设备ID", nil)
		return
	}

	var device model.Device
	if err := db.DB.First(&device, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "设备不存在", nil)
		return
	}

	status := "在线"
	if device.Status != nil && *device.Status != "" {
		status = *device.Status
	}

	deviceType := "集控"
	if device.Type != nil && *device.Type != "" {
		deviceType = *device.Type
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"id":        device.ID,
		"name":      device.Name,
		"code":      device.Code,
		"merchsId":  device.MerchsID,
		"roomsId":   device.RoomsID,
		"groupsId":  device.GroupsID,
		"status":    status,
		"type":      deviceType,
		"createdAt": device.CreatedAt,
		"updatedAt": device.UpdatedAt,
	})
}

// CreateDevice 创建设备
// 根据设备类型自动生成设备名称（类型+3位随机数）
func CreateDevice(c *gin.Context) {
	var req struct {
		Code     string `json:"code" binding:"required"`
		Type     string `json:"type"`
		MerchsID int32  `json:"merchsId" binding:"required"`
		RoomsID  int32  `json:"roomsId"`
		GroupsID int32  `json:"groupsId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", nil)
		return
	}

	// 检查设备编码是否已存在
	var existingDevice model.Device
	if err := db.DB.Where("code = ?", req.Code).First(&existingDevice).Error; err == nil {
		response.Fail(c, http.StatusBadRequest, "设备编码已存在", nil)
		return
	}

	// 设置设备类型，默认为"集控"
	deviceType := req.Type
	if deviceType == "" {
		deviceType = "集控"
	}

	// 生成设备名称：设备类型 + 3位随机数
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(900) + 100 // 生成100-999的随机数
	deviceName := deviceType + strconv.Itoa(randomNum)

	// 创建设备
	device := model.Device{
		Name:     deviceName,
		Code:     req.Code,
		MerchsID: req.MerchsID,
	}

	if req.RoomsID != 0 {
		device.RoomsID = &req.RoomsID
	}
	if req.GroupsID != 0 {
		device.GroupsID = &req.GroupsID
	}

	// 设置设备类型
	device.Type = &deviceType

	if err := db.DB.Create(&device).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建设备失败", nil)
		return
	}

	response.SuccessWithMsg(c, "创建成功", gin.H{
		"id":        device.ID,
		"name":      device.Name,
		"code":      device.Code,
		"type":      deviceType,
		"createdAt": device.CreatedAt,
	})
}

// UpdateDevice 更新设备
func UpdateDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的设备ID", nil)
		return
	}

	var req struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", nil)
		return
	}

	var device model.Device
	if err := db.DB.First(&device, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "设备不存在", nil)
		return
	}

	if req.Name != "" {
		device.Name = req.Name
	}

	if req.Type != "" {
		device.Type = &req.Type
	}

	if req.Code != "" {
		// 检查新编码是否被其他设备使用
		var existingDevice model.Device
		if err := db.DB.Where("code = ? AND id != ?", req.Code, id).First(&existingDevice).Error; err == nil {
			response.Fail(c, http.StatusBadRequest, "设备编码已存在", nil)
			return
		}
		device.Code = req.Code
	}

	if err := db.DB.Save(&device).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新设备失败", nil)
		return
	}

	response.SuccessWithMsg(c, "更新成功", gin.H{
		"id":        device.ID,
		"name":      device.Name,
		"updatedAt": device.UpdatedAt,
	})
}

// DeleteDevice 删除设备
func DeleteDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的设备ID", nil)
		return
	}

	var device model.Device
	if err := db.DB.First(&device, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "设备不存在", nil)
		return
	}

	if err := db.DB.Delete(&device).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除设备失败", nil)
		return
	}

	response.SuccessWithMsg(c, "删除成功", gin.H{"success": true})
}
