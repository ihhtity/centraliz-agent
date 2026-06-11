package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetGroupList 获取分组列表
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

	if err := query.Order("id DESC").Find(&groups).Error; err != nil {
		response.Fail(c, 500, "获取分组列表失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"list":  groups,
		"total": total,
	})
}

// GetGroupDetail 获取分组详情
func GetGroupDetail(c *gin.Context) {
	id := c.Param("id")
	merchsID := c.Query("merchsId")
	rulesID := c.Query("rulesId")

	var group model.Group
	if err := db.DB.Where("id = ?", id).First(&group).Error; err != nil {
		response.Fail(c, 404, "分组不存在", nil)
		return
	}

	// 查询分组下的房间数量
	var roomCount int64
	db.DB.Model(&model.Room{}).Where("groups_id = ?", group.ID).Count(&roomCount)
	// 查询规则名称
	var rulename string
	db.DB.Model(&model.Rule{}).Where("id = ?", rulesID).Pluck("name", &rulename)

	// 查询规则列表
	var rules []model.Rule
	query := db.DB.Model(&model.Rule{})
	if merchsID != "" {
		query = query.Where("merchs_id = ?", merchsID)
	}
	if err := query.Order("id ASC").Find(&rules).Error; err != nil {
		response.Fail(c, 500, "获取规则列表失败: "+err.Error())
		return
	}

	// 查询分组下的所有设备
	var devices []model.Device
	if err := db.DB.Where("groups_id = ?", group.ID).Find(&devices).Error; err != nil {
		response.Fail(c, 500, "获取设备列表失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "获取成功", gin.H{
		"group": gin.H{
			"id":          group.ID,
			"name":        group.Name,
			"merchsId":    group.MerchsID,
			"rulesId":     group.RulesID,
			"rulename":    rulename,
			"phone":       group.Phone,
			"type":        group.Type,
			"bindNumber":  group.BindNumber,
			"consumePush": group.ConsumePush,
			"location":    group.Location,
			"count":       roomCount,
			"createdAt":   group.CreatedAt,
		},
		"rules":   rules,
		"devices": devices,
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

	// 验证手机号格式
	if req.Phone != "" {
		if len(req.Phone) > 20 {
			response.Fail(c, 400, "手机号不能超过20个字符")
			return
		}
	}

	// 验证位置信息
	if req.Location != "" && len(req.Location) > 255 {
		response.Fail(c, 400, "位置信息不能超过255个字符")
		return
	}

	// 检查同一商家下是否存在同名分组
	var exists bool
	if err := db.DB.Model(&model.Group{}).
		Select("count(*) > 0").
		Where("merchs_id = ? AND name = ?", req.MerchsID, req.Name).
		Find(&exists).Error; err != nil {
		response.Fail(c, 500, "检查分组名称失败: "+err.Error())
		return
	}

	if exists {
		response.Fail(c, 400, "已存在同名分组")
		return
	}

	// 创建分组对象
	group := model.Group{
		Name:     &req.Name,
		MerchsID: req.MerchsID,
		Type:     &req.Type,
		Count:    new(uint32),
	}

	if req.RulesID != nil && *req.RulesID > 0 {
		group.RulesID = req.RulesID
	}

	if req.Phone != "" {
		group.Phone = &req.Phone
	}

	if req.Location != "" {
		group.Location = &req.Location
	}

	// 创建分组
	if err := db.DB.Create(&group).Error; err != nil {
		response.Fail(c, 500, "创建分组失败: "+err.Error())
		return
	}

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
		Name        string `json:"name"`
		RulesID     *int32 `json:"rulesId"`
		Phone       string `json:"phone"`
		Type        string `json:"type"`
		Location    string `json:"location"`
		Rulename    string `json:"rulename"`
		BindNumber  string `json:"bind_number"`
		ConsumePush string `json:"consume_push"`
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

	// 解析分组ID
	id, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, 400, "无效的分组ID")
		return
	}

	var group model.Group
	if err := db.DB.First(&group, id).Error; err != nil {
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

	// 更新手机号
	if req.Phone != "" {
		if len(req.Phone) > 20 {
			response.Fail(c, 400, "手机号不能超过20个字符")
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

	// 更新分组类型
	if req.Type != "" {
		if req.Type != "存柜" && req.Type != "零售" {
			response.Fail(c, 400, "分组类型只能是存柜或零售")
			return
		}
		group.Type = &req.Type
	}

	// 更新绑定号码设置
	if req.BindNumber != "" {
		if req.BindNumber != "关闭" && req.BindNumber != "手动" && req.BindNumber != "自动" {
			response.Fail(c, 400, "绑定号码设置只能是 关闭、手动 或 自动")
			return
		}
		group.BindNumber = &req.BindNumber
	}

	// 更新消费推送设置
	if req.ConsumePush != "" {
		if req.ConsumePush != "关闭" && req.ConsumePush != "开启" {
			response.Fail(c, 400, "消费推送设置只能是 关闭 或 开启")
			return
		}
		group.ConsumePush = &req.ConsumePush
	}

	// 更新规则ID
	if req.RulesID != nil && *req.RulesID > 0 {
		group.RulesID = req.RulesID
		group.RuleName = &req.Rulename
	}

	// 更新分组绑定的房间数量
	var roomCount int64
	if err := db.DB.Model(&model.Room{}).Where("groups_id = ?", group.ID).Count(&roomCount).Error; err != nil {
		response.Fail(c, 500, "查询房间数量失败: "+err.Error())
		return
	}
	group.Count = new(uint32)
	*group.Count = uint32(roomCount)

	// 更新分组
	if err := db.DB.Save(&group).Error; err != nil {
		response.Fail(c, 500, "更新分组失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "更新成功", gin.H{
		"id":          group.ID,
		"name":        group.Name,
		"merchsId":    group.MerchsID,
		"rulesId":     group.RulesID,
		"phone":       group.Phone,
		"count":       group.Count,
		"type":        group.Type,
		"location":    group.Location,
		"bindNumber":  group.BindNumber,
		"consumePush": group.ConsumePush,
		"updatedAt":   group.UpdatedAt,
	})
}

// DeleteGroup 删除分组
// 删除指定ID的分组，需要验证商家密码，删除分组时同时删除其下的所有房间
func DeleteGroup(c *gin.Context) {
	groupID := c.Param("id")

	if groupID == "" {
		response.Fail(c, 400, "分组ID不能为空")
		return
	}

	// 解析分组ID
	id, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, 400, "无效的分组ID")
		return
	}

	// 获取请求中的密码
	type DeleteRequest struct {
		Password string `json:"password"`
	}
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	if req.Password == "" {
		response.Fail(c, 400, "请输入商家密码")
		return
	}

	// 查询分组信息
	var group model.Group
	if err := db.DB.First(&group, id).Error; err != nil {
		response.Fail(c, 404, "分组不存在")
		return
	}

	// 查询商家信息验证密码
	var merch model.Merch
	if err := db.DB.First(&merch, group.MerchsID).Error; err != nil {
		response.Fail(c, 404, "商家不存在")
		return
	}

	// 验证密码（使用bcrypt）
	if err := bcrypt.CompareHashAndPassword([]byte(merch.Password), []byte(req.Password)); err != nil {
		response.Fail(c, 400, "密码验证失败")
		return
	}

	// 使用事务删除分组及关联的房间
	tx := db.DB.Begin()

	// 删除分组下的所有房间
	if err := tx.Where("groups_id = ?", group.ID).Delete(&model.Room{}).Error; err != nil {
		tx.Rollback()
		response.Fail(c, 500, "删除房间失败: "+err.Error())
		return
	}

	// 删除分组
	if err := tx.Delete(&group).Error; err != nil {
		tx.Rollback()
		response.Fail(c, 500, "删除分组失败: "+err.Error())
		return
	}

	tx.Commit()

	response.SuccessWithMsg(c, "删除成功", gin.H{"success": true})
}
