package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"fmt"
	"math/rand"
	"net"
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
		Name      string `json:"name"`
		Type      string `json:"type"`
		Code      string `json:"code"`
		LockCount *int32 `json:"lockCount"`
		Status    string `json:"status"`
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

	// 更新锁定数量
	if req.LockCount != nil {
		device.LockCount = req.LockCount
	}

	// 更新设备状态
	if req.Status != "" {
		if req.Status != "在线" && req.Status != "离线" && req.Status != "维修" {
			response.Fail(c, http.StatusBadRequest, "设备状态只能是在线、离线或维修", nil)
			return
		}
		device.Status = &req.Status
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

// AddDeviceControl 添加设备控制记录
func AddDeviceControl(c *gin.Context) {
	var req struct {
		Code     string `json:"code" binding:"required"`
		Recharge string `json:"recharge" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error(), nil)
		return
	}

	// 验证设备编码
	if req.Code == "" {
		response.Fail(c, http.StatusBadRequest, "设备编码不能为空", nil)
		return
	}

	// 验证充值编码
	if req.Recharge == "" {
		response.Fail(c, http.StatusBadRequest, "充值编码不能为空", nil)
		return
	}

	// 查询设备是否存在
	var device model.Device
	if err := db.DB.Where("code = ?", req.Code).First(&device).Error; err == nil {
		// 设备已存在，返回错误
		response.Fail(c, http.StatusConflict, "设备已存在", nil)
		return
	}

	deviceType := "集控"
	device.Type = &deviceType
	device.Name = req.Code

	if err := db.DB.Create(&device).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建设备失败: "+err.Error(), nil)
		return
	}

	// 记录设备控制日志（这里可以扩展为存储到日志表）
	controlLog := gin.H{
		"deviceId":     device.ID,
		"deviceCode":   req.Code,
		"rechargeCode": req.Recharge,
	}

	// 使用封装的 EleControlCommon 发送控制指令
	// conn, err := EleControlCommon(req.Code, req.Recharge, "", 0)
	// if err != nil {
	// 	controlLog["error"] = err.Error()
	// 	response.Fail(c, http.StatusInternalServerError, "控制失败: "+err.Error(), controlLog)
	// 	return
	// }
	// defer conn.Close()
	// controlLog["status"] = "指令已发送"

	response.SuccessWithMsg(c, "控制成功", controlLog)
}

// EleControlCommon 控制电控设备公共方法
func EleControlCommon(code string, command string, realmName string) (string, error) {
	// 验证必填参数
	if code == "" {
		return "", fmt.Errorf("设备编码不能为空")
	}
	if command == "" {
		return "", fmt.Errorf("控制指令不能为空")
	}

	// 设置默认服务器地址
	if realmName == "" {
		realmName = "centraliztcp.bsldtech.cn" // 服务器域名地址
	}

	// 将设备编码 转换为 HEX 格式
	var output string
	for _, char := range code {
		hex := fmt.Sprintf("%02X", char)
		output += hex
	}

	// 创建 TCP 连接
	addr := fmt.Sprintf("%s:%d", realmName, 12341)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return "", fmt.Errorf("解析 TCP 地址失败：%w", err)
	}

	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		return "", fmt.Errorf("连接服务端失败：%w", err)
	}

	// 发送控制指令
	_, err = conn.Write([]byte("c"))
	if err != nil {
		conn.Close()
		return "", fmt.Errorf("发送指令失败：%w", err)
	}

	_, err = conn.Write([]byte(output))
	if err != nil {
		conn.Close()
		return "", fmt.Errorf("发送设备 ID 失败：%w", err)
	}

	_, err = conn.Write([]byte(command))
	if err != nil {
		conn.Close()
		return "", fmt.Errorf("发送命令失败：%w", err)
	}

	// 设置读取超时
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// 读取服务器返回的数据
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		conn.Close()
		return "", fmt.Errorf("读取响应失败：%w", err)
	}

	// 关闭连接
	conn.Close()

	// 返回服务器响应数据
	return string(buffer[:n]), nil
}

// DeviceCommon 控制设备接口
func DeviceCommon(c *gin.Context) {
	var req struct {
		Code      string `json:"code" binding:"required"`
		Command   string `json:"command" binding:"required"`
		RealmName string `json:"realmName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error(), nil)
		return
	}

	// 调用 EleControlCommon 发送控制指令
	responseData, err := EleControlCommon(req.Code, req.Command, req.RealmName)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "控制失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMsg(c, "控制成功", gin.H{
		"code":      req.Code,
		"command":   req.Command,
		"realmName": req.RealmName,
		"status":    "指令已发送",
		"response":  responseData,
	})
}

// BindDeviceToGroup 绑定设备到分组
// 根据设备编码查询设备，判断是否被其他分组或商家绑定，若无则更新设备数据
func BindDeviceToGroup(c *gin.Context) {
	var req struct {
		Code     string `json:"code" binding:"required"`
		GroupID  int64  `json:"groupId" binding:"required"`
		MerchsID int32  `json:"merchsId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error(), nil)
		return
	}

	// 根据设备编码查询设备
	var device model.Device
	if err := db.DB.Where("code = ?", req.Code).First(&device).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "设备不存在", nil)
		return
	}

	// 判断设备是否已被其他商家绑定
	if device.MerchsID != 0 && device.MerchsID != req.MerchsID {
		response.Fail(c, http.StatusConflict, "设备已被其他商家绑定", nil)
		return
	}

	// 判断设备是否已被其他分组绑定
	if device.GroupsID != nil && *device.GroupsID != 0 && int64(*device.GroupsID) != req.GroupID {
		response.Fail(c, http.StatusConflict, "设备已被其他分组绑定", nil)
		return
	}

	// 更新设备数据，绑定到当前分组
	device.MerchsID = req.MerchsID
	groupID := int32(req.GroupID)
	device.GroupsID = &groupID

	if err := db.DB.Save(&device).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "绑定失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMsg(c, "绑定成功", gin.H{
		"id":       device.ID,
		"code":     device.Code,
		"name":     device.Name,
		"groupsId": device.GroupsID,
	})
}

// UnbindDeviceFromGroup 解除设备与分组的绑定
// 根据设备ID，将设备的分组外键更新为0
func UnbindDeviceFromGroup(c *gin.Context) {
	var req struct {
		DeviceID int64 `json:"deviceId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error(), nil)
		return
	}

	// 查询设备是否存在
	var device model.Device
	if err := db.DB.First(&device, req.DeviceID).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "设备不存在", nil)
		return
	}

	// 将分组外键更新为0（解除绑定）
	zero := int32(0)
	device.GroupsID = &zero

	if err := db.DB.Save(&device).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "解除绑定失败: "+err.Error(), nil)
		return
	}

	response.SuccessWithMsg(c, "解除绑定成功", gin.H{
		"id":       device.ID,
		"code":     device.Code,
		"groupsId": device.GroupsID,
	})
}
