package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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

	// 将房间转换为map以便追加额外字段
	roomMap := make(map[string]interface{})
	roomJSON, _ := json.Marshal(room)
	json.Unmarshal(roomJSON, &roomMap)

	roomMap["groupName"] = groupName
	roomMap["groupType"] = groupType
	roomMap["device"] = device

	response.SuccessWithMsg(c, "获取成功", roomMap)
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
		DeviceID int32  `json:"device_id"`
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
	if req.DeviceID > 0 {
		if err := db.DB.Where("id = ?", req.DeviceID).First(&device).Error; err == nil {
			if device.LockCount != 0 {
				deviceLockCount = device.LockCount
			}
		}
	}

	// 获取该设备下已使用的锁号列表
	var usedLockNos []string
	if req.DeviceID > 0 {
		if err := db.DB.Model(&model.Room{}).
			Where("devices_id = ? AND lock_no IS NOT NULL", req.DeviceID).
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
			Name:      roomName,
			MerchsID:  req.MerchsID,
			DevicesID: req.DeviceID,
			Tag:       req.Tag,
			BoardNo:   boardNoStr,
			LockNo:    lockNoStr,
			FreeTime:  time.Now(),
		}

		if req.GroupsID != nil && *req.GroupsID != 0 {
			room.GroupsID = *req.GroupsID
		}

		if req.RulesID != nil && *req.RulesID != 0 {
			room.RulesID = *req.RulesID
		}

		if req.DeviceID > 0 {
			room.DevicesID = req.DeviceID
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
		Price  string `json:"price"`
		Image  string `json:"image"`
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
		if room.DevicesID > 0 {
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
	if req.Price != "" {
		price, err := strconv.ParseFloat(req.Price, 64)
		if err != nil {
			response.Fail(c, 400, "价格格式错误")
			return
		}

		room.Price = float32(price)
	}
	if req.Image != "" {
		room.Image = req.Image
	}
	// 更新免费时间
	room.FreeTime = time.Now()

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

	// 将设备信息添加到房间数据中
	lockers := make([]gin.H, 0)
	lockers = append(lockers, gin.H{
		"id":        room.ID,
		"merchsId":  room.MerchsID,
		"usersid":   room.UsersID,
		"groupsId":  room.GroupsID,
		"devicesId": room.DevicesID,
		"ordersId":  room.OrdersID,
		"name":      room.Name,
		"tag":       room.Tag,
		"status":    room.Status,
		"boardNo":   room.BoardNo,
		"lockNo":    room.LockNo,
		"freeTime":  room.FreeTime,
		"combo":     room.Combo,
		"device":    device,
	})

	// 获取分组信息
	var group model.Group
	var rules *model.Rule
	if room.GroupsID > 0 {
		if err := db.DB.Where("id = ?", room.GroupsID).Order("id ASC").Find(&group).Error; err == nil {
			// 获取分组关联的规则
			if group.RulesID > 0 {
				var rule model.Rule
				if err := db.DB.Where("id = ?", group.RulesID).Order("id ASC").Find(&rule).Error; err == nil {
					rules = &rule
				}
			}
		}
	}

	response.SuccessWithMsg(c, "获取成功", []gin.H{{
		"lockers": lockers,
		"groups":  group,
		"rules":   rules,
	}})
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
			"usersid":   room.UsersID,
			"groupsId":  room.GroupsID,
			"devicesId": room.DevicesID,
			"ordersId":  room.OrdersID,
			"name":      room.Name,
			"tag":       room.Tag,
			"status":    room.Status,
			"boardNo":   room.BoardNo,
			"lockNo":    room.LockNo,
			"freeTime":  room.FreeTime,
			"combo":     room.Combo,
			"device":    device,
		})
	}

	response.SuccessWithMsg(c, "获取成功", []gin.H{{
		"lockers": lockers,
		"groups":  group,
		"rules":   rules,
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
				"usersid":   room.UsersID,
				"groupsId":  room.GroupsID,
				"devicesId": room.DevicesID,
				"ordersId":  room.OrdersID,
				"name":      room.Name,
				"tag":       room.Tag,
				"status":    room.Status,
				"boardNo":   room.BoardNo,
				"lockNo":    room.LockNo,
				"freeTime":  room.FreeTime,
				"combo":     room.Combo,
				"device":    device,
			})
		}

		if len(lockers) > 0 {
			result = append(result, gin.H{
				"lockers": lockers,
				"groups":  group,
				"rules":   rules,
			})
		}
	}

	response.SuccessWithMsg(c, "获取成功", result)
}

// ==================== 房间标签相关接口 ====================

// GetRoomTagList 获取房间标签列表
func GetRoomTagList(c *gin.Context) {
	var tags []model.RoomTag
	merchID := c.Param("merchId")
	if merchID == "" || merchID == "0" {
		response.Fail(c, 400, "商家ID不能为空")
		return
	}

	if err := db.DB.Where("merchs_id = ?", merchID).Order("id ASC").Find(&tags).Error; err != nil {
		response.Fail(c, 500, "获取标签列表失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":  tags,
		"total": len(tags),
	})
}

// CreateRoomTag 创建房间标签
func CreateRoomTag(c *gin.Context) {
	type CreateRoomTagRequest struct {
		MerchsID int32  `json:"merchsId" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	var req CreateRoomTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	if req.MerchsID == 0 {
		response.Fail(c, 400, "商家ID不能为空")
		return
	}

	if req.Name == "" {
		response.Fail(c, 400, "标签名称不能为空")
		return
	}

	tag := model.RoomTag{
		MerchsID: req.MerchsID,
		Name:     req.Name,
	}

	if err := db.DB.Create(&tag).Error; err != nil {
		response.Fail(c, 500, "创建标签失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", tag)
}

// UpdateRoomTag 更新房间标签
func UpdateRoomTag(c *gin.Context) {
	type UpdateRoomTagRequest struct {
		MerchsID int32  `json:"merchsId" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	var req UpdateRoomTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	if req.MerchsID == 0 {
		response.Fail(c, 400, "商家ID不能为空")
		return
	}

	tagID := c.Param("id")
	if tagID == "" {
		response.Fail(c, 400, "标签ID不能为空")
		return
	}

	var tag model.RoomTag
	if err := db.DB.Where("id = ?", tagID).First(&tag).Error; err != nil {
		response.Fail(c, 404, "标签不存在")
		return
	}

	if tag.MerchsID != req.MerchsID {
		response.Fail(c, 400, "标签所属商家与请求商家不一致")
		return
	}

	if req.Name == "" {
		response.Fail(c, 400, "标签名称不能为空")
		return
	}

	tag.Name = req.Name
	if err := db.DB.Save(&tag).Error; err != nil {
		response.Fail(c, 500, "更新标签失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", tag)
}

// DeleteRoomTag 删除房间标签
func DeleteRoomTag(c *gin.Context) {
	tagID := c.Param("id")
	if tagID == "" {
		response.Fail(c, 400, "标签ID不能为空")
		return
	}

	var tag model.RoomTag
	if err := db.DB.Where("id = ?", tagID).First(&tag).Error; err != nil {
		response.Fail(c, 404, "标签不存在")
		return
	}

	if err := db.DB.Delete(&tag).Error; err != nil {
		response.Fail(c, 500, "删除标签失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", gin.H{"success": true})
}

// ==================== 房间图片相关接口 ====================

// GetRoomImageList 获取房间图片列表
func GetRoomImageList(c *gin.Context) {
	var images []model.RoomImage
	if err := db.DB.Order("created_at DESC").Find(&images).Error; err != nil {
		response.Fail(c, 500, "获取图片列表失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":  images,
		"total": len(images),
	})
}

// CreateRoomImage 创建房间图片
func CreateRoomImage(c *gin.Context) {
	type CreateRoomImageRequest struct {
		Name  string `json:"name" binding:"required"`
		Image string `json:"image" binding:"required"`
	}

	var req CreateRoomImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	if req.Name == "" {
		response.Fail(c, 400, "图片名称不能为空")
		return
	}

	if req.Image == "" {
		response.Fail(c, 400, "图片地址不能为空")
		return
	}

	image := model.RoomImage{
		Name:  req.Name,
		Image: req.Image,
	}

	if err := db.DB.Create(&image).Error; err != nil {
		response.Fail(c, 500, "创建图片失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "创建成功", image)
}

// UpdateRoomImage 更新房间图片
func UpdateRoomImage(c *gin.Context) {
	type UpdateRoomImageRequest struct {
		Name string `json:"name"`
	}

	var req UpdateRoomImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	imageID := c.Param("id")
	if imageID == "" {
		response.Fail(c, 400, "图片ID不能为空")
		return
	}

	var image model.RoomImage
	if err := db.DB.Where("id = ?", imageID).First(&image).Error; err != nil {
		response.Fail(c, 404, "图片不存在")
		return
	}

	if req.Name != "" {
		image.Name = req.Name
	}

	if err := db.DB.Save(&image).Error; err != nil {
		response.Fail(c, 500, "更新图片失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", image)
}

// DeleteRoomImage 删除房间图片
func DeleteRoomImage(c *gin.Context) {
	imageID := c.Param("id")
	if imageID == "" {
		response.Fail(c, 400, "图片ID不能为空")
		return
	}

	var image model.RoomImage
	if err := db.DB.Where("id = ?", imageID).First(&image).Error; err != nil {
		response.Fail(c, 404, "图片不存在")
		return
	}

	if err := db.DB.Delete(&image).Error; err != nil {
		response.Fail(c, 500, "删除图片失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "删除成功", gin.H{"success": true})
}

// UploadRoomImage 上传房间图片
func UploadRoomImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		response.Fail(c, 400, "请选择要上传的图片")
		return
	}

	if file.Size > 10*1024*1024 {
		response.Fail(c, 400, "图片大小不能超过10MB")
		return
	}

	ext := file.Filename[len(file.Filename)-4:]
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" && ext != ".gif" {
		response.Fail(c, 400, "只支持jpg、png、jpeg、gif格式的图片")
		return
	}

	filePath := "./uploads/room_images/"
	err = c.SaveUploadedFile(file, filePath+file.Filename)
	if err != nil {
		response.Fail(c, 500, "上传图片失败: "+err.Error())
		return
	}

	imageURL := "/uploads/room_images/" + file.Filename

	response.SuccessWithMsg(c, "上传成功", gin.H{"url": imageURL})
}
