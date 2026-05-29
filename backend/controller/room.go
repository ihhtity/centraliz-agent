package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRoomList 获取房间列表
func GetRoomList(c *gin.Context) {
	merchsID := c.Query("merchs_id")
	groupsID := c.Query("groups_id")

	if merchsID == "" {
		response.Fail(c, 400, "商家ID不能为空")
		return
	}

	var rooms []model.Room
	query := db.DB.Model(&model.Room{}).Where("merchs_id = ?", merchsID)

	if groupsID != "" && groupsID != "0" {
		query = query.Where("groups_id = ?", groupsID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, 500, "获取房间总数失败: "+err.Error())
		return
	}

	if err := query.Order("id ASC").Find(&rooms).Error; err != nil {
		response.Fail(c, 500, "获取房间列表失败: "+err.Error())
		return
	}

	roomList := make([]gin.H, len(rooms))
	for i, room := range rooms {
		// 查询房间绑定的设备数量
		var deviceCount int64
		db.DB.Model(&model.Device{}).Where("rooms_id = ?", room.ID).Count(&deviceCount)

		roomList[i] = gin.H{
			"id":          room.ID,
			"name":        room.Name,
			"merchsId":    room.MerchsID,
			"groupsId":    room.GroupsID,
			"rulesId":     room.RulesID,
			"tag":         room.Tag,
			"status":      room.Status,
			"deviceCount": deviceCount,
			"createdAt":   room.CreatedAt,
			"updatedAt":   room.UpdatedAt,
		}
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":  roomList,
		"total": total,
	})
}

// GetRoomDetail 获取房间详情（包含绑定的设备）
func GetRoomDetail(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		response.Fail(c, 400, "房间ID不能为空")
		return
	}

	var room model.Room
	if err := db.DB.Where("id = ?", roomID).First(&room).Error; err != nil {
		response.Fail(c, 404, "房间不存在")
		return
	}

	// 获取分组名称
	var groupName string
	if room.GroupsID != nil {
		var group model.Group
		if err := db.DB.Where("id = ?", *room.GroupsID).First(&group).Error; err == nil && group.Name != nil {
			groupName = *group.Name
		}
	}

	// 获取绑定的设备列表
	var devices []model.Device
	var networkStatus = "离线"
	var signalStrength = 0
	var powerStatus = "关电"
	var lockStatus = "闭锁"

	if err := db.DB.Where("rooms_id = ?", room.ID).Find(&devices).Error; err == nil && len(devices) > 0 {
		// 如果有绑定设备，取第一个设备的状态作为房间状态
		device := devices[0]
		if device.Status != nil && *device.Status == "在线" {
			networkStatus = "在线"
			signalStrength = 90
		}
	}

	// 构建设备列表响应
	deviceList := make([]gin.H, 0, len(devices))
	for _, device := range devices {
		status := "离线"
		if device.Status != nil {
			status = *device.Status
		}
		deviceType := "集控"
		if device.Type != nil {
			deviceType = *device.Type
		}
		deviceList = append(deviceList, gin.H{
			"id":        device.ID,
			"name":      device.Name,
			"code":      device.Code,
			"type":      deviceType,
			"status":    status,
			"createdAt": device.CreatedAt,
		})
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"id":             room.ID,
		"name":           room.Name,
		"tag":            room.Tag,
		"status":         room.Status,
		"groupsId":       room.GroupsID,
		"groupName":      groupName,
		"merchsId":       room.MerchsID,
		"rulesId":        room.RulesID,
		"devices":        deviceList,
		"deviceCount":    len(devices),
		"networkStatus":  networkStatus,
		"signalStrength": signalStrength,
		"powerStatus":    powerStatus,
		"lockStatus":     lockStatus,
		"createdAt":      room.CreatedAt,
		"updatedAt":      room.UpdatedAt,
	})
}

// CreateRoom 创建房间（支持批量创建）
func CreateRoom(c *gin.Context) {
	type CreateRoomRequest struct {
		Name     string `json:"name" binding:"required"`
		Count    int    `json:"count"` // 创建数量，默认1，范围1-50
		MerchsID int32  `json:"merchs_id" binding:"required"`
		GroupsID *int32 `json:"groups_id"`
		RulesID  *int32 `json:"rules_id"`
		Tag      string `json:"tag"`
	}

	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	if req.MerchsID <= 0 {
		response.Fail(c, 400, "无效的商家ID")
		return
	}

	if len(req.Name) > 50 {
		response.Fail(c, 400, "房间名称不能超过50个字符")
		return
	}

	if len(req.Name) == 0 {
		response.Fail(c, 400, "房间名称不能为空")
		return
	}

	// 设置默认数量为1
	if req.Count <= 0 {
		req.Count = 1
	}
	if req.Count > 50 {
		req.Count = 50
	}

	if req.Tag != "" && len(req.Tag) > 30 {
		response.Fail(c, 400, "房间标签不能超过30个字符")
		return
	}

	if req.Tag == "" {
		req.Tag = "普通柜"
	}

	if req.GroupsID != nil && *req.GroupsID <= 0 {
		response.Fail(c, 400, "无效的分组ID")
		return
	}

	if req.RulesID != nil && *req.RulesID <= 0 {
		response.Fail(c, 400, "无效的规则ID")
		return
	}

	// 批量创建房间
	createdRooms := make([]gin.H, 0, req.Count)
	tx := db.DB.Begin()

	for i := 0; i < req.Count; i++ {
		var roomName string
		if req.Count == 1 {
			// 只创建一个房间时，直接使用原始名称
			roomName = req.Name
		} else {
			// 创建多个房间时，从第二个开始加数字
			if i == 0 {
				roomName = req.Name
			} else {
				prefix, startNum := parseRoomName(req.Name)
				roomName = generateRoomName(prefix, startNum+i)
			}
		}

		// 检查房间名称是否已存在
		var existCount int64
		if err := tx.Model(&model.Room{}).
			Where("merchs_id = ? AND name = ?", req.MerchsID, roomName).
			Count(&existCount).Error; err != nil {
			tx.Rollback()
			response.Fail(c, 500, "检查房间名称失败: "+err.Error())
			return
		}

		if existCount > 0 {
			tx.Rollback()
			response.Fail(c, 400, "房间名称 "+roomName+" 已存在")
			return
		}

		room := model.Room{
			Name:     &roomName,
			MerchsID: req.MerchsID,
			Tag:      &req.Tag,
		}

		if req.GroupsID != nil {
			room.GroupsID = req.GroupsID
		}

		if req.RulesID != nil {
			room.RulesID = req.RulesID
		}

		if err := tx.Create(&room).Error; err != nil {
			tx.Rollback()
			response.Fail(c, 500, "创建房间失败: "+err.Error())
			return
		}

		createdRooms = append(createdRooms, gin.H{
			"id":        room.ID,
			"name":      room.Name,
			"merchsId":  room.MerchsID,
			"groupsId":  room.GroupsID,
			"rulesId":   room.RulesID,
			"tag":       room.Tag,
			"createdAt": room.CreatedAt,
		})
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, 500, "提交事务失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", gin.H{
		"count": len(createdRooms),
		"list":  createdRooms,
	})
}

// parseRoomName 解析房间名称，提取前缀和数字
// 例如："A01" -> ("A", 1), "B123" -> ("B", 123), "Room5" -> ("Room", 5)
func parseRoomName(name string) (string, int) {
	if name == "" {
		return "", 1
	}

	// 从后往前找数字
	i := len(name) - 1
	for i >= 0 && name[i] >= '0' && name[i] <= '9' {
		i--
	}

	// 如果没有数字，直接返回原名称
	if i == len(name)-1 {
		return name, 1
	}

	prefix := name[:i+1]
	numStr := name[i+1:]

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return name, 1
	}

	return prefix, num
}

// generateRoomName 生成房间名称，保持数字位数一致
// 例如：("A", 1) -> "A01", ("A", 12) -> "A12"
func generateRoomName(prefix string, num int) string {
	if prefix == "" {
		return strconv.Itoa(num)
	}
	return prefix + strconv.Itoa(num)
}

// UpdateRoom 更新房间
func UpdateRoom(c *gin.Context) {
	type UpdateRoomRequest struct {
		Name     string `json:"name"`
		GroupsID *int32 `json:"groups_id"`
		RulesID  *int32 `json:"rules_id"`
		Tag      string `json:"tag"`
		Status   string `json:"status"`
	}

	var req UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	roomID := c.Param("id")
	if roomID == "" {
		response.Fail(c, 400, "房间ID不能为空")
		return
	}

	var room model.Room
	if err := db.DB.Where("id = ?", roomID).First(&room).Error; err != nil {
		response.Fail(c, 404, "房间不存在")
		return
	}

	if req.Name != "" {
		if len(req.Name) > 50 {
			response.Fail(c, 400, "房间名称不能超过50个字符")
			return
		}
		room.Name = &req.Name
	}

	if req.Tag != "" {
		if len(req.Tag) > 30 {
			response.Fail(c, 400, "房间标签不能超过30个字符")
			return
		}
		room.Tag = &req.Tag
	}

	if req.GroupsID != nil {
		if *req.GroupsID <= 0 {
			response.Fail(c, 400, "无效的分组ID")
			return
		}
		room.GroupsID = req.GroupsID
	}

	if req.RulesID != nil {
		if *req.RulesID <= 0 {
			response.Fail(c, 400, "无效的规则ID")
			return
		}
		room.RulesID = req.RulesID
	}

	if err := db.DB.Save(&room).Error; err != nil {
		response.Fail(c, 500, "更新房间失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", gin.H{
		"id":        room.ID,
		"name":      room.Name,
		"merchsId":  room.MerchsID,
		"groupsId":  room.GroupsID,
		"rulesId":   room.RulesID,
		"tag":       room.Tag,
		"status":    room.Status,
		"updatedAt": room.UpdatedAt,
	})
}

// DeleteRoom 删除房间
func DeleteRoom(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		response.Fail(c, 400, "房间ID不能为空")
		return
	}

	var room model.Room
	if err := db.DB.Where("id = ?", roomID).First(&room).Error; err != nil {
		response.Fail(c, 404, "房间不存在")
		return
	}

	if err := db.DB.Delete(&room).Error; err != nil {
		response.Fail(c, 500, "删除房间失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", gin.H{"success": true})
}

// BindDevice 绑定设备到房间
func BindDevice(c *gin.Context) {
	type BindDeviceRequest struct {
		RoomID   uint64 `json:"roomId" binding:"required"`
		DeviceID uint64 `json:"deviceId" binding:"required"`
	}

	var req BindDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 检查房间是否存在
	var room model.Room
	if err := db.DB.Where("id = ?", req.RoomID).First(&room).Error; err != nil {
		response.Fail(c, 404, "房间不存在")
		return
	}

	// 检查设备是否存在
	var device model.Device
	if err := db.DB.Where("id = ?", req.DeviceID).First(&device).Error; err != nil {
		response.Fail(c, 404, "设备不存在")
		return
	}

	// 检查设备是否已绑定到其他房间
	if device.RoomsID != nil && *device.RoomsID != 0 && *device.RoomsID != int32(req.RoomID) {
		response.Fail(c, 400, "设备已绑定到其他房间")
		return
	}

	// 绑定设备到房间
	roomID := int32(req.RoomID)
	device.RoomsID = &roomID

	if err := db.DB.Save(&device).Error; err != nil {
		response.Fail(c, 500, "绑定设备失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "绑定成功", gin.H{
		"roomId":   req.RoomID,
		"deviceId": req.DeviceID,
	})
}

// UnbindDevice 解绑房间的设备
func UnbindDevice(c *gin.Context) {
	type UnbindDeviceRequest struct {
		RoomID   uint64 `json:"roomId" binding:"required"`
		DeviceID uint64 `json:"deviceId" binding:"required"`
	}

	var req UnbindDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 检查房间是否存在
	var room model.Room
	if err := db.DB.Where("id = ?", req.RoomID).First(&room).Error; err != nil {
		response.Fail(c, 404, "房间不存在")
		return
	}

	// 检查设备是否存在
	var device model.Device
	if err := db.DB.Where("id = ?", req.DeviceID).First(&device).Error; err != nil {
		response.Fail(c, 404, "设备不存在")
		return
	}

	// 检查设备是否绑定到此房间
	if device.RoomsID == nil || *device.RoomsID != int32(req.RoomID) {
		response.Fail(c, 400, "设备未绑定到此房间")
		return
	}

	// 解绑设备
	device.RoomsID = nil

	if err := db.DB.Save(&device).Error; err != nil {
		response.Fail(c, 500, "解绑设备失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "解绑成功", gin.H{
		"roomId":   req.RoomID,
		"deviceId": req.DeviceID,
	})
}
