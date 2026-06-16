package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRoomList 获取房间列表
func GetRoomList(c *gin.Context) {
	groupsID := c.Query("groups_id")

	if groupsID == "" || groupsID == "0" {
		response.Fail(c, 400, "分组ID不能为空")
		return
	}

	// 获取分组下的房间列表
	var rooms []model.Room
	query := db.DB.Model(&model.Room{}).Where("groups_id = ?", groupsID)
	if err := query.Order("id ASC").Find(&rooms).Error; err != nil {
		response.Fail(c, 500, "获取房间列表失败: "+err.Error())
		return
	}

	// 查询设备列表
	var devices []model.Device
	devicequery := db.DB.Model(&model.Device{}).Where("groups_id = ? AND type = ?", groupsID, "集控")
	if err := devicequery.Order("id ASC").Find(&devices).Error; err != nil {
		response.Fail(c, 500, "获取设备列表失败: "+err.Error())
		return
	}

	// 计算分组设备锁定数量总和
	var totalLockCount int64
	if err := db.DB.Model(&model.Device{}).
		Where("groups_id = ?", groupsID).
		Select("COALESCE(SUM(lock_count), 0)").
		Scan(&totalLockCount).Error; err != nil {
		response.Fail(c, 500, "获取设备锁定数量失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":      rooms,
		"devices":   devices,
		"total":     len(rooms),
		"lockCount": totalLockCount,
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
	var groupType string
	if room.GroupsID > 0 {
		var group model.Group
		if err := db.DB.Where("id = ?", room.GroupsID).First(&group).Error; err == nil && group.Name != "" {
			groupName = group.Name
			groupType = group.Type
		}
	}

	// 获取绑定的设备信息
	var device model.Device
	if err := db.DB.Where("id = ?", room.DevicesID).First(&device).Error; err != nil {
		device = model.Device{}
	}

	// 获取设备状态
	deviceStatus := device.Status
	if deviceStatus == "" {
		deviceStatus = "在线"
	}

	// 获取设备类型
	deviceType := device.Type
	if deviceType == "" {
		deviceType = "集控"
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"id":        room.ID,
		"name":      room.Name,
		"tag":       room.Tag,
		"status":    room.Status,
		"groupsId":  room.GroupsID,
		"groupName": groupName,
		"groupType": groupType,
		"merchsId":  room.MerchsID,
		"rulesId":   room.RulesID,
		"devicesId": room.DevicesID,
		"device": gin.H{
			"id":        device.ID,
			"name":      device.Name,
			"code":      device.Code,
			"boardNo":   device.BoardNo,
			"status":    deviceStatus,
			"type":      deviceType,
			"createdAt": device.CreatedAt,
		},
		"boardNo": room.BoardNo,
		"lockNo":  room.LockNo,
	})
}

// CreateRoom 创建房间（支持批量创建）
func CreateRoom(c *gin.Context) {
	type CreateRoomRequest struct {
		Name     string `json:"name" binding:"required"`
		Count    int    `json:"count"`
		BoardNo  int    `json:"board_no"`
		LockNo   int    `json:"lock_no"`
		MerchsID int32  `json:"merchs_id" binding:"required"`
		GroupsID *int32 `json:"groups_id"`
		RulesID  *int32 `json:"rules_id"`
		Tag      string `json:"tag"`
		DeviceID *int64 `json:"device_id"`
	}

	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Count <= 0 {
		req.Count = 1
	}

	if req.Tag == "" {
		req.Tag = "普通柜"
	}

	// 获取设备信息，用于自动分配锁号
	var device model.Device
	var deviceLockCount int32 = 500 // 默认锁数量上限
	if req.DeviceID != nil {
		if err := db.DB.Where("id = ?", *req.DeviceID).First(&device).Error; err == nil {
			if device.LockCount != 0 {
				deviceLockCount = device.LockCount
			}
		}
	}

	// 获取该设备下已使用的锁号列表
	var usedLockNos []string
	if req.DeviceID != nil {
		if err := db.DB.Model(&model.Room{}).
			Where("devices_id = ? AND lock_no IS NOT NULL", strconv.FormatInt(*req.DeviceID, 10)).
			Pluck("lock_no", &usedLockNos).Error; err != nil {
			response.Fail(c, 500, "获取已使用锁号失败: "+err.Error())
			return
		}
	}

	// 转换为整数集合方便查找
	usedLockNoSet := make(map[int]bool)
	for _, lockNoStr := range usedLockNos {
		if lockNoStr != "" {
			if lockNo, err := strconv.Atoi(lockNoStr); err == nil {
				usedLockNoSet[lockNo] = true
			}
		}
	}

	// 确定起始锁号：优先使用用户传入的，否则自动分配最小可用锁号
	startLockNo := req.LockNo
	if startLockNo <= 0 {
		// 查找最小可用锁号
		for i := 1; i <= int(deviceLockCount); i++ {
			if !usedLockNoSet[i] {
				startLockNo = i
				break
			}
		}
	}

	// 验证锁号是否可用
	if startLockNo <= 0 || startLockNo > int(deviceLockCount) {
		response.Fail(c, 400, "无法分配有效的锁号，设备锁数量已达上限")
		return
	}

	// 批量创建房间
	skipCount := 0
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
			// 房间名称已存在，跳过处理，继续添加剩余房间
			skipCount++
			continue
		}

		// 计算当前房间的锁号
		currentLockNo := startLockNo + i
		if currentLockNo > int(deviceLockCount) {
			tx.Rollback()
			response.Fail(c, 400, "设备锁数量不足，无法创建所有房间")
			return
		}

		// 转换板号和锁号为固定格式（板号2位，锁号2-3位）
		boardNoStr := formatBoardNo(req.BoardNo) // 2位数，如 "01", "11", "99"
		lockNoStr := formatLockNo(currentLockNo) // 2-3位数，如 "01", "58", "158"

		room := model.Room{
			Name:     roomName,
			MerchsID: req.MerchsID,
			Tag:      req.Tag,
			BoardNo:  boardNoStr,
			LockNo:   lockNoStr,
		}

		if req.GroupsID != nil && *req.GroupsID != 0 {
			room.GroupsID = *req.GroupsID
		}

		if req.RulesID != nil && *req.RulesID != 0 {
			room.RulesID = *req.RulesID
		}

		if req.DeviceID != nil {
			room.DevicesID = strconv.FormatInt(*req.DeviceID, 10)
		}

		if err := tx.Create(&room).Error; err != nil {
			tx.Rollback()
			response.Fail(c, 500, "创建房间失败: "+err.Error())
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, 500, "提交事务失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", gin.H{
		"count":     req.Count - skipCount,
		"skipCount": skipCount,
	})
}

// UpdateRoom 更新房间
func UpdateRoom(c *gin.Context) {
	type UpdateRoomRequest struct {
		Name   string `json:"name"`
		Tag    string `json:"tag"`
		LockNo string `json:"lock_no"`
		Status string `json:"status"`
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

	// 检查锁号是否被其他房间使用
	if req.LockNo != "" {
		var count int64
		query := db.DB.Model(&model.Room{}).Where("lock_no = ? AND id != ?", req.LockNo, roomID)
		if room.DevicesID != "" {
			query = query.Where("devices_id = ?", room.DevicesID)
		}
		if err := query.Count(&count).Error; err != nil {
			response.Fail(c, 500, "检查锁号失败: "+err.Error())
			return
		}
		if count > 0 {
			response.Fail(c, 400, "锁号已被其他房间使用，请更换锁号")
			return
		}
	}

	// 更新房间数据
	if req.Name != "" {
		room.Name = req.Name
	}
	if req.Tag != "" {
		room.Tag = req.Tag
	}
	if req.LockNo != "" {
		// 格式化锁号
		room.LockNo = formatLockNoString(req.LockNo)
	}
	if req.Status != "" {
		room.Status = req.Status
	}

	if err := db.DB.Save(&room).Error; err != nil {
		response.Fail(c, 500, "更新房间失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", gin.H{"success": true})
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

// GenerateQRCode 生成二维码
func GenerateQRCode(c *gin.Context) {
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

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"roomId":        room.ID,
		"roomName":      room.Name,
		"h5QrcodeUrl":   "https://centraliz.bsldtech.cn/#/pages/room/detail?id=" + roomID,
		"miniQrcodeUrl": "https://centraliz.bsldtech.cn/api/qrcode/miniprogram?roomId=" + roomID,
	})
}

// formatBoardNo 格式化板号为2位数，如 1 -> "01", 11 -> "11", 99 -> "99"
func formatBoardNo(boardNo int) string {
	if boardNo < 0 {
		boardNo = 0
	}
	if boardNo > 99 {
		boardNo = 99
	}
	return fmt.Sprintf("%02d", boardNo)
}

// formatBoardNoFixed 格式化板号为固定2位数
func formatBoardNoFixed(boardNo int) string {
	return fmt.Sprintf("%02d", boardNo)
}

// formatLockNo 格式化锁号为2-3位数，如 1 -> "01", 58 -> "58", 158 -> "158"
func formatLockNo(lockNo int) string {
	if lockNo < 0 {
		lockNo = 0
	}
	if lockNo > 999 {
		lockNo = 999
	}
	// 锁号：0-9 -> 2位(00-09), 10-99 -> 2位(10-99), 100-999 -> 3位(100-999)
	if lockNo < 10 {
		return fmt.Sprintf("%02d", lockNo)
	}
	return strconv.Itoa(lockNo)
}

// formatLockNoString 格式化锁号字符串为2-3位数
func formatLockNoString(lockNoStr string) string {
	lockNo, err := strconv.Atoi(lockNoStr)
	if err != nil {
		// 如果转换失败，尝试提取数字
		for _, c := range lockNoStr {
			if c >= '0' && c <= '9' {
				lockNo = lockNo*10 + int(c-'0')
			}
		}
	}
	return formatLockNo(lockNo)
}

// ==================== 用户端柜子相关接口 ====================

// GetUserRoomDetail 获取单个柜子详情（用户端）
func GetUserRoomDetail(c *gin.Context) {
	roomID := c.Param("id")
	if roomID == "" {
		response.Fail(c, 400, "柜子ID不能为空")
		return
	}

	// 获取柜子详情
	var room model.Room
	if err := db.DB.Where("id = ?", roomID).Order("id ASC").Find(&room).Error; err != nil {
		response.Fail(c, 404, "柜子不存在")
		return
	}

	// 获取柜子绑定的设备详情
	var device model.Device
	if err := db.DB.Where("id = ?", room.DevicesID).First(&device).Error; err != nil {
		response.Fail(c, 404, "设备不存在")
		return
	}

	// 获取分组信息
	var group model.Group
	var groupName, groupType string
	var rules *model.Rule
	if room.GroupsID > 0 {
		if err := db.DB.Where("id = ?", room.GroupsID).Order("id ASC").Find(&group).Error; err == nil {
			groupName = group.Name
			groupType = group.Type

			// 获取分组关联的规则
			if group.RulesID > 0 {
				var rule model.Rule
				if err := db.DB.Where("id = ?", group.RulesID).Order("id ASC").Find(&rule).Error; err == nil {
					rules = &rule
				}
			}
		}
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"id":        room.ID,
		"merchsId":  room.MerchsID,
		"devicesId": room.DevicesID,
		"usersid":   room.UsersID,
		"name":      room.Name,
		"tag":       room.Tag,
		"price":     5,
		"status":    room.Status,
		"boardNo":   room.BoardNo,
		"lockNo":    room.LockNo,
		"groupsId":  room.GroupsID,
		"groupName": groupName,
		"groupType": groupType,
		"rules":     rules,
		"device":    device,
	})
}

// GetUserRoomListByGroup 获取分组下的柜子列表（用户端）
func GetUserRoomListByGroup(c *gin.Context) {
	groupID := c.Param("groupId")
	if groupID == "" || groupID == "0" {
		response.Fail(c, 400, "分组ID不能为空")
		return
	}

	// 获取分组信息
	var group model.Group
	if err := db.DB.Where("id = ?", groupID).Order("id ASC").Find(&group).Error; err != nil {
		response.Fail(c, 404, "分组不存在")
		return
	}

	// 获取分组关联的规则
	var rules *model.Rule
	if group.RulesID > 0 {
		var rule model.Rule
		if err := db.DB.Where("id = ?", group.RulesID).Order("id ASC").Find(&rule).Error; err == nil {
			rules = &rule
		}
	}

	// 获取该分组下的所有房间（柜子）
	var rooms []model.Room
	if err := db.DB.Where("groups_id = ?", groupID).Order("id ASC").Find(&rooms).Error; err != nil {
		response.Fail(c, 500, "获取柜子列表失败: "+err.Error())
		return
	}

	// 转换房间数据
	lockers := make([]gin.H, 0)
	for _, room := range rooms {
		// 获取柜子绑定的设备详情
		var device model.Device
		if err := db.DB.Where("id = ?", room.DevicesID).First(&device).Error; err != nil {
			continue
		}

		lockers = append(lockers, gin.H{
			"id":        room.ID,
			"merchsId":  room.MerchsID,
			"devicesId": room.DevicesID,
			"usersid":   room.UsersID,
			"name":      room.Name,
			"tag":       room.Tag,
			"status":    room.Status,
			"boardNo":   room.BoardNo,
			"lockNo":    room.LockNo,
			"device":    device,
		})
	}

	response.SuccessWithMsg(c, "获取成功", []gin.H{{
		"groupId":   group.ID,
		"groupName": group.Name,
		"groupType": group.Type,
		"rules":     rules,
		"lockers":   lockers,
	}})
}

// GetUserRoomListByMerchant 获取商家全部分组下的柜子列表（用户端）
func GetUserRoomListByMerchant(c *gin.Context) {
	merchID := c.Param("merchId")
	if merchID == "" || merchID == "0" {
		response.Fail(c, 400, "商家ID不能为空")
		return
	}

	// 获取商家的所有分组
	var groups []model.Group
	if err := db.DB.Where("merchs_id = ?", merchID).Order("id ASC").Find(&groups).Error; err != nil {
		response.Fail(c, 500, "获取分组列表失败: "+err.Error())
		return
	}

	if len(groups) == 0 {
		response.SuccessWithMsg(c, "获取成功", []gin.H{})
		return
	}

	// 按分组获取柜子数据
	result := make([]gin.H, 0)
	for _, group := range groups {
		// 获取分组关联的规则
		var rules *model.Rule
		if group.RulesID > 0 {
			var rule model.Rule
			if err := db.DB.Where("id = ?", group.RulesID).Order("id ASC").Find(&rule).Error; err == nil {
				rules = &rule
			}
		}

		// 获取该分组下的所有房间（柜子）
		var rooms []model.Room
		if err := db.DB.Where("groups_id = ?", group.ID).Order("id ASC").Find(&rooms).Error; err != nil {
			continue
		}

		// 转换房间数据
		lockers := make([]gin.H, 0)
		for _, room := range rooms {
			// 获取柜子绑定的设备详情
			var device model.Device
			if err := db.DB.Where("id = ?", room.DevicesID).First(&device).Error; err != nil {
				continue
			}

			lockers = append(lockers, gin.H{
				"id":        room.ID,
				"merchsId":  room.MerchsID,
				"devicesId": room.DevicesID,
				"usersid":   room.UsersID,
				"name":      room.Name,
				"tag":       room.Tag,
				"price":     5,
				"status":    room.Status,
				"boardNo":   room.BoardNo,
				"lockNo":    room.LockNo,
				"device":    device,
			})
		}

		if len(lockers) > 0 {
			result = append(result, gin.H{
				"groupId":   group.ID,
				"groupName": group.Name,
				"groupType": group.Type,
				"rules":     rules,
				"lockers":   lockers,
			})
		}
	}

	response.SuccessWithMsg(c, "获取成功", result)
}
