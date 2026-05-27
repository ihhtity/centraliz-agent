package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

// GetGroupList 获取分组列表
// 支持通过merchs_id参数过滤特定商家的分组
func GetGroupList(c *gin.Context) {
	merchsID := c.Query("merchs_id")

	var groups []model.Group
	query := db.DB.Model(&model.Group{})

	if merchsID != "" {
		query = query.Where("merchs_id = ?", merchsID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.Fail(c, 500, "获取分组总数失败: "+err.Error())
		return
	}

	if err := query.Order("created_at DESC").Find(&groups).Error; err != nil {
		response.Fail(c, 500, "获取分组列表失败: "+err.Error())
		return
	}

	groupList := make([]gin.H, len(groups))
	for i, group := range groups {
		// 查询分组下的房间数量
		var roomCount int64
		db.DB.Model(&model.Room{}).Where("groups_id = ?", group.ID).Count(&roomCount)

		// 查询分组下的设备数量（通过房间关联）
		var deviceCount int64
		db.DB.Model(&model.Device{}).
			Joins("LEFT JOIN rooms ON devices.rooms_id = rooms.id").
			Where("rooms.groups_id = ?", group.ID).
			Count(&deviceCount)

		groupList[i] = gin.H{
			"id":          group.ID,
			"name":        group.Name,
			"merchsId":    group.MerchsID,
			"rulesId":     group.RulesID,
			"phone":       group.Phone,
			"count":       roomCount,
			"deviceCount": deviceCount,
			"type":        group.Type,
			"location":    group.Location,
			"createdAt":   group.CreatedAt,
			"updatedAt":   group.UpdatedAt,
		}
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":  groupList,
		"total": total,
	})
}

// GetGroupDetail 获取分组详情
// 根据分组ID获取单个分组的详细信息，包含关联的房间列表和设备统计
func GetGroupDetail(c *gin.Context) {
	id := c.Param("id")

	var group model.Group
	if err := db.DB.Where("id = ?", id).First(&group).Error; err != nil {
		response.Fail(c, 404, "分组不存在", nil)
		return
	}

	// 查询分组下的房间列表
	var rooms []model.Room
	db.DB.Where("groups_id = ?", id).Order("created_at DESC").Find(&rooms)

	// 查询分组下的设备数量
	var deviceCount int64
	db.DB.Model(&model.Device{}).
		Joins("LEFT JOIN rooms ON devices.rooms_id = rooms.id").
		Where("rooms.groups_id = ?", id).
		Count(&deviceCount)

	// 构建房间列表响应
	roomList := make([]gin.H, len(rooms))
	for i, room := range rooms {
		// 查询房间绑定的设备数量
		var roomDeviceCount int64
		db.DB.Model(&model.Device{}).Where("rooms_id = ?", room.ID).Count(&roomDeviceCount)

		roomList[i] = gin.H{
			"id":          room.ID,
			"name":        room.Name,
			"status":      room.Status,
			"tag":         room.Tag,
			"deviceCount": roomDeviceCount,
			"createdAt":   room.CreatedAt,
		}
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"id":          group.ID,
		"name":        group.Name,
		"merchsId":    group.MerchsID,
		"rulesId":     group.RulesID,
		"phone":       group.Phone,
		"count":       len(rooms),
		"deviceCount": deviceCount,
		"type":        group.Type,
		"location":    group.Location,
		"rooms":       roomList,
		"createdAt":   group.CreatedAt,
		"updatedAt":   group.UpdatedAt,
	})
}

// CreateGroup 创建分组
// 创建新的分组，包含完整的参数验证和业务逻辑检查
func CreateGroup(c *gin.Context) {
	type CreateGroupRequest struct {
		Name     string `json:"name" binding:"required"`
		MerchsID int32  `json:"merchs_id" binding:"required"`
		RulesID  *int32 `json:"rules_id"`
		Phone    string `json:"phone"`
		Type     string `json:"type"`
		Location string `json:"location"`
	}

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 验证商家ID
	if req.MerchsID <= 0 {
		response.Fail(c, 400, "无效的商家ID")
		return
	}

	// 验证分组名称
	if len(req.Name) == 0 {
		response.Fail(c, 400, "分组名称不能为空")
		return
	}
	if len(req.Name) > 100 {
		response.Fail(c, 400, "分组名称不能超过100个字符")
		return
	}

	// 设置默认类型
	if req.Type == "" {
		req.Type = "0"
	}

	// 验证分组类型
	validType := map[string]bool{"0": true, "1": true}
	if !validType[req.Type] {
		response.Fail(c, 400, "无效的分组类型，必须是0(仓储)或1(零售)")
		return
	}

	// 验证手机号格式
	if req.Phone != "" {
		if len(req.Phone) > 20 {
			response.Fail(c, 400, "手机号不能超过20个字符")
			return
		}
		// 简单的手机号格式验证
		phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
		if !phoneRegex.MatchString(req.Phone) {
			response.Fail(c, 400, "手机号格式不正确")
			return
		}
	}

	// 验证位置信息
	if req.Location != "" && len(req.Location) > 255 {
		response.Fail(c, 400, "位置信息不能超过255个字符")
		return
	}

	// 检查同一商家下是否存在同名分组
	var count int64
	if err := db.DB.Model(&model.Group{}).
		Where("merchs_id = ? AND name = ?", req.MerchsID, req.Name).
		Count(&count).Error; err != nil {
		response.Fail(c, 500, "检查分组名称失败: "+err.Error())
		return
	}

	if count > 0 {
		response.Fail(c, 400, "该商家下已存在同名分组")
		return
	}

	// 创建分组对象
	group := model.Group{
		Name:     &req.Name,
		MerchsID: req.MerchsID,
		Type:     &req.Type,
		Count:    new(uint32),
	}

	if req.RulesID != nil {
		group.RulesID = req.RulesID
	}

	if req.Phone != "" {
		group.Phone = &req.Phone
	}

	if req.Location != "" {
		group.Location = &req.Location
	}

	// 使用事务创建分组
	tx := db.DB.Begin()
	if err := tx.Create(&group).Error; err != nil {
		tx.Rollback()
		response.Fail(c, 500, "创建分组失败: "+err.Error())
		return
	}
	tx.Commit()

	response.SuccessWithMsg(c, "创建成功", gin.H{
		"id":        group.ID,
		"name":      group.Name,
		"merchsId":  group.MerchsID,
		"rulesId":   group.RulesID,
		"phone":     group.Phone,
		"count":     group.Count,
		"type":      group.Type,
		"location":  group.Location,
		"createdAt": group.CreatedAt,
	})
}

// UpdateGroup 更新分组
// 更新分组信息，支持部分字段更新
func UpdateGroup(c *gin.Context) {
	type UpdateGroupRequest struct {
		Name      string     `json:"name"`
		RulesID   *int32     `json:"rules_id"`
		Phone     string     `json:"phone"`
		Type      string     `json:"type"`
		Location  string     `json:"location"`
		UpdatedAt *time.Time `json:"updated_at"`
	}

	var req UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	groupID := c.Param("id")
	if groupID == "" {
		response.Fail(c, 400, "分组ID不能为空")
		return
	}

	var group model.Group
	if err := db.DB.Where("id = ?", groupID).First(&group).Error; err != nil {
		response.Fail(c, 404, "分组不存在")
		return
	}

	// 更新分组名称
	if req.Name != "" {
		if len(req.Name) > 100 {
			response.Fail(c, 400, "分组名称不能超过100个字符")
			return
		}
		group.Name = &req.Name
	}

	// 更新分组类型
	if req.Type != "" {
		validType := map[string]bool{"0": true, "1": true}
		if !validType[req.Type] {
			response.Fail(c, 400, "无效的分组类型，必须是0(仓储)或1(零售)")
			return
		}
		group.Type = &req.Type
	}

	// 更新手机号
	if req.Phone != "" {
		if len(req.Phone) > 20 {
			response.Fail(c, 400, "手机号不能超过20个字符")
			return
		}
		phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
		if !phoneRegex.MatchString(req.Phone) {
			response.Fail(c, 400, "手机号格式不正确")
			return
		}
		group.Phone = &req.Phone
	}

	// 更新位置信息
	if req.Location != "" {
		if len(req.Location) > 255 {
			response.Fail(c, 400, "位置信息不能超过255个字符")
			return
		}
		group.Location = &req.Location
	}

	// 更新规则ID
	if req.RulesID != nil {
		if *req.RulesID <= 0 {
			response.Fail(c, 400, "无效的规则ID")
			return
		}
		group.RulesID = req.RulesID
	}

	// 使用事务更新
	tx := db.DB.Begin()
	if err := tx.Save(&group).Error; err != nil {
		tx.Rollback()
		response.Fail(c, 500, "更新分组失败: "+err.Error())
		return
	}
	tx.Commit()

	response.SuccessWithMsg(c, "更新成功", gin.H{
		"id":        group.ID,
		"name":      group.Name,
		"merchsId":  group.MerchsID,
		"rulesId":   group.RulesID,
		"phone":     group.Phone,
		"count":     group.Count,
		"type":      group.Type,
		"location":  group.Location,
		"updatedAt": group.UpdatedAt,
	})
}

// DeleteGroup 删除分组
// 删除指定ID的分组，删除前检查是否有关联的房间
func DeleteGroup(c *gin.Context) {
	groupID := c.Param("id")

	if groupID == "" {
		response.Fail(c, 400, "分组ID不能为空")
		return
	}

	var group model.Group
	if err := db.DB.Where("id = ?", groupID).First(&group).Error; err != nil {
		response.Fail(c, 404, "分组不存在")
		return
	}

	// 检查是否有关联的房间
	var roomCount int64
	if err := db.DB.Model(&model.Room{}).
		Where("groups_id = ?", group.ID).
		Count(&roomCount).Error; err != nil {
		response.Fail(c, 500, "检查关联房间失败: "+err.Error())
		return
	}

	if roomCount > 0 {
		response.Fail(c, 400, "该分组下存在关联的房间，请先删除或转移关联房间")
		return
	}

	// 使用事务删除
	tx := db.DB.Begin()
	if err := tx.Delete(&group).Error; err != nil {
		tx.Rollback()
		response.Fail(c, 500, "删除分组失败: "+err.Error())
		return
	}
	tx.Commit()

	response.SuccessWithMsg(c, "删除成功", gin.H{"success": true})
}
