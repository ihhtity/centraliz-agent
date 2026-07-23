package controller

import (
	"centraliz-backend/model"
	"centraliz-backend/pkg/db"
	"centraliz-backend/pkg/response"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func parseExcelFile(c *gin.Context) (*excelize.File, error) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, 400, "请选择要导入的文件")
		return nil, err
	}

	if file.Size > 10*1024*1024 {
		response.Fail(c, 400, "文件大小不能超过10MB")
		return nil, fmt.Errorf("file too large")
	}

	ext := file.Filename[len(file.Filename)-4:]
	if ext != ".xlsx" && ext != ".xls" {
		response.Fail(c, 400, "只支持xlsx、xls格式的Excel文件")
		return nil, fmt.Errorf("invalid file format")
	}

	openedFile, err := file.Open()
	if err != nil {
		response.Fail(c, 500, "读取文件失败: "+err.Error())
		return nil, err
	}
	defer openedFile.Close()

	xlsx, err := excelize.OpenReader(openedFile)
	if err != nil {
		response.Fail(c, 500, "解析Excel文件失败: "+err.Error())
		return nil, err
	}

	return xlsx, nil
}

func AdminImportRoom(c *gin.Context) {
	xlsx, err := parseExcelFile(c)
	if err != nil {
		return
	}

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		response.Fail(c, 500, "读取Excel行失败: "+err.Error())
		return
	}

	if len(rows) <= 1 {
		response.Fail(c, 400, "Excel文件为空或只有表头")
		return
	}

	successCount := 0
	skipCount := 0
	failCount := 0
	failMessages := make([]string, 0)

	tx := db.DB.Begin()

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}

		name := row[0]
		if name == "" {
			continue
		}

		merchsID := int32(0)
		if len(row) > 1 && row[1] != "" {
			if v, err := strconv.Atoi(row[1]); err == nil {
				merchsID = int32(v)
			}
		}

		if merchsID == 0 {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 商家ID不能为空", i+1))
			continue
		}

		groupsID := int32(0)
		if len(row) > 2 && row[2] != "" {
			if v, err := strconv.Atoi(row[2]); err == nil {
				groupsID = int32(v)
			}
		}

		tag := "普通柜"
		if len(row) > 3 && row[3] != "" {
			tag = row[3]
		}

		var existCount int64
		if err := tx.Model(&model.Room{}).Where("merchs_id = ? AND name = ?", merchsID, name).Count(&existCount).Error; err != nil {
			tx.Rollback()
			response.Fail(c, 500, "检查房间名称失败: "+err.Error())
			return
		}

		if existCount > 0 {
			skipCount++
			continue
		}

		room := &model.Room{
			Name:      name,
			MerchsID:  merchsID,
			GroupsID:  groupsID,
			Tag:       tag,
			Status:    "空闲",
			FreeTime:  time.Now(),
		}

		if len(row) > 4 && row[4] != "" {
			room.BoardNo = row[4]
		}

		if len(row) > 5 && row[5] != "" {
			room.LockNo = row[5]
		}

		if len(row) > 6 && row[6] != "" {
			if price, err := strconv.ParseFloat(row[6], 32); err == nil {
				room.Price = float32(price)
			}
		}

		if len(row) > 7 && row[7] != "" {
			room.Image = row[7]
		}

		if err := tx.Create(room).Error; err != nil {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 创建失败 - %s", i+1, err.Error()))
			continue
		}

		successCount++
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, 500, "提交事务失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "导入完成", gin.H{
		"success": successCount,
		"skip":    skipCount,
		"fail":    failCount,
		"errors":  failMessages,
	})
}

func AdminImportDevice(c *gin.Context) {
	xlsx, err := parseExcelFile(c)
	if err != nil {
		return
	}

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		response.Fail(c, 500, "读取Excel行失败: "+err.Error())
		return
	}

	if len(rows) <= 1 {
		response.Fail(c, 400, "Excel文件为空或只有表头")
		return
	}

	successCount := 0
	skipCount := 0
	failCount := 0
	failMessages := make([]string, 0)

	tx := db.DB.Begin()

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}

		name := row[0]
		code := row[1]

		if name == "" || code == "" {
			continue
		}

		merchsID := int32(0)
		if len(row) > 2 && row[2] != "" {
			if v, err := strconv.Atoi(row[2]); err == nil {
				merchsID = int32(v)
			}
		}

		if merchsID == 0 {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 商家ID不能为空", i+1))
			continue
		}

		var existCount int64
		if err := tx.Model(&model.Device{}).Where("code = ?", code).Count(&existCount).Error; err != nil {
			tx.Rollback()
			response.Fail(c, 500, "检查设备编码失败: "+err.Error())
			return
		}

		if existCount > 0 {
			skipCount++
			continue
		}

		deviceType := "集控"
		if len(row) > 3 && row[3] != "" {
			deviceType = row[3]
		}

		groupsID := int32(0)
		if len(row) > 4 && row[4] != "" {
			if v, err := strconv.Atoi(row[4]); err == nil {
				groupsID = int32(v)
			}
		}

		lockCount := int32(500)
		if len(row) > 5 && row[5] != "" {
			if v, err := strconv.Atoi(row[5]); err == nil {
				lockCount = int32(v)
			}
		}

		device := &model.Device{
			Name:      name,
			Code:      code,
			MerchsID:  merchsID,
			Type:      deviceType,
			GroupsID:  groupsID,
			LockCount: lockCount,
			Status:    "在线",
		}

		if len(row) > 6 && row[6] != "" {
			device.BoardNo = row[6]
		}

		if err := tx.Create(device).Error; err != nil {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 创建失败 - %s", i+1, err.Error()))
			continue
		}

		successCount++
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, 500, "提交事务失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "导入完成", gin.H{
		"success": successCount,
		"skip":    skipCount,
		"fail":    failCount,
		"errors":  failMessages,
	})
}

func AdminImportMerch(c *gin.Context) {
	xlsx, err := parseExcelFile(c)
	if err != nil {
		return
	}

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		response.Fail(c, 500, "读取Excel行失败: "+err.Error())
		return
	}

	if len(rows) <= 1 {
		response.Fail(c, 400, "Excel文件为空或只有表头")
		return
	}

	successCount := 0
	skipCount := 0
	failCount := 0
	failMessages := make([]string, 0)

	tx := db.DB.Begin()

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}

		account := row[0]
		password := row[1]

		if account == "" || password == "" {
			continue
		}

		var existCount int64
		if err := tx.Model(&model.Merch{}).Where("account = ?", account).Count(&existCount).Error; err != nil {
			tx.Rollback()
			response.Fail(c, 500, "检查商家账号失败: "+err.Error())
			return
		}

		if existCount > 0 {
			skipCount++
			continue
		}

		role := "商家"
		if len(row) > 2 && row[2] != "" {
			role = row[2]
		}

		status := "0"
		if len(row) > 3 && row[3] != "" {
			status = row[3]
		}

		phone := ""
		if len(row) > 4 && row[4] != "" {
			phone = row[4]
		}

		email := ""
		if len(row) > 5 && row[5] != "" {
			email = row[5]
		}

		merch := &model.Merch{
			Account:  account,
			Password: password,
			Role:     role,
			Status:   status,
			Phone:    phone,
			Email:    email,
		}

		if err := tx.Create(merch).Error; err != nil {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 创建失败 - %s", i+1, err.Error()))
			continue
		}

		successCount++
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, 500, "提交事务失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "导入完成", gin.H{
		"success": successCount,
		"skip":    skipCount,
		"fail":    failCount,
		"errors":  failMessages,
	})
}

func AdminImportGroup(c *gin.Context) {
	xlsx, err := parseExcelFile(c)
	if err != nil {
		return
	}

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		response.Fail(c, 500, "读取Excel行失败: "+err.Error())
		return
	}

	if len(rows) <= 1 {
		response.Fail(c, 400, "Excel文件为空或只有表头")
		return
	}

	successCount := 0
	skipCount := 0
	failCount := 0
	failMessages := make([]string, 0)

	tx := db.DB.Begin()

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}

		name := row[0]
		if name == "" {
			continue
		}

		merchsID := int32(0)
		if len(row) > 1 && row[1] != "" {
			if v, err := strconv.Atoi(row[1]); err == nil {
				merchsID = int32(v)
			}
		}

		if merchsID == 0 {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 商家ID不能为空", i+1))
			continue
		}

		var existCount int64
		if err := tx.Model(&model.Group{}).Where("merchs_id = ? AND name = ?", merchsID, name).Count(&existCount).Error; err != nil {
			tx.Rollback()
			response.Fail(c, 500, "检查分组名称失败: "+err.Error())
			return
		}

		if existCount > 0 {
			skipCount++
			continue
		}

		groupType := "存柜"
		if len(row) > 2 && row[2] != "" {
			groupType = row[2]
		}

		rulesID := int32(0)
		if len(row) > 3 && row[3] != "" {
			if v, err := strconv.Atoi(row[3]); err == nil {
				rulesID = int32(v)
			}
		}

		phone := ""
		if len(row) > 4 && row[4] != "" {
			phone = row[4]
		}

		location := ""
		if len(row) > 5 && row[5] != "" {
			location = row[5]
		}

		bindNumber := "关闭"
		if len(row) > 6 && row[6] != "" {
			bindNumber = row[6]
		}

		consumePush := "关闭"
		if len(row) > 7 && row[7] != "" {
			consumePush = row[7]
		}

		group := &model.Group{
			Name:        name,
			MerchsID:    merchsID,
			Type:        groupType,
			RulesID:     rulesID,
			Phone:       phone,
			Location:    location,
			BindNumber:  bindNumber,
			ConsumePush: consumePush,
		}

		if err := tx.Create(group).Error; err != nil {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 创建失败 - %s", i+1, err.Error()))
			continue
		}

		successCount++
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, 500, "提交事务失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "导入完成", gin.H{
		"success": successCount,
		"skip":    skipCount,
		"fail":    failCount,
		"errors":  failMessages,
	})
}

func AdminImportRule(c *gin.Context) {
	xlsx, err := parseExcelFile(c)
	if err != nil {
		return
	}

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		response.Fail(c, 500, "读取Excel行失败: "+err.Error())
		return
	}

	if len(rows) <= 1 {
		response.Fail(c, 400, "Excel文件为空或只有表头")
		return
	}

	successCount := 0
	skipCount := 0
	failCount := 0
	failMessages := make([]string, 0)

	tx := db.DB.Begin()

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}

		name := row[0]
		if name == "" {
			continue
		}

		merchsID := int32(0)
		if len(row) > 1 && row[1] != "" {
			if v, err := strconv.Atoi(row[1]); err == nil {
				merchsID = int32(v)
			}
		}

		if merchsID == 0 {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 商家ID不能为空", i+1))
			continue
		}

		var existCount int64
		if err := tx.Model(&model.Rule{}).Where("merchs_id = ? AND name = ?", merchsID, name).Count(&existCount).Error; err != nil {
			tx.Rollback()
			response.Fail(c, 500, "检查规则名称失败: "+err.Error())
			return
		}

		if existCount > 0 {
			skipCount++
			continue
		}

		ruleType := "free"
		if len(row) > 2 && row[2] != "" {
			ruleType = row[2]
		}

		mode := "single"
		if len(row) > 3 && row[3] != "" {
			mode = row[3]
		}

		price := float32(0)
		if len(row) > 4 && row[4] != "" {
			if v, err := strconv.ParseFloat(row[4], 32); err == nil {
				price = float32(v)
			}
		}

		deposit := float32(0)
		if len(row) > 5 && row[5] != "" {
			if v, err := strconv.ParseFloat(row[5], 32); err == nil {
				deposit = float32(v)
			}
		}

		durationUnit := "hour"
		if len(row) > 6 && row[6] != "" {
			durationUnit = row[6]
		}

		rule := &model.Rule{
			Name:         name,
			MerchsID:     merchsID,
			Type:         ruleType,
			Mode:         mode,
			Price:        price,
			Deposit:      deposit,
			DurationUnit: durationUnit,
		}

		if len(row) > 7 && row[7] != "" {
			if v, err := strconv.Atoi(row[7]); err == nil {
				rule.AutoEndTime = int32(v)
			}
		}

		if len(row) > 8 && row[8] != "" {
			rule.Description = row[8]
		}

		if len(row) > 9 && row[9] != "" {
			if v, err := strconv.Atoi(row[9]); err == nil {
				rule.FreeTime = int32(v)
			}
		}

		if err := tx.Create(rule).Error; err != nil {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 创建失败 - %s", i+1, err.Error()))
			continue
		}

		successCount++
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, 500, "提交事务失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "导入完成", gin.H{
		"success": successCount,
		"skip":    skipCount,
		"fail":    failCount,
		"errors":  failMessages,
	})
}

func AdminImportSubMerch(c *gin.Context) {
	xlsx, err := parseExcelFile(c)
	if err != nil {
		return
	}

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		response.Fail(c, 500, "读取Excel行失败: "+err.Error())
		return
	}

	if len(rows) <= 1 {
		response.Fail(c, 400, "Excel文件为空或只有表头")
		return
	}

	successCount := 0
	skipCount := 0
	failCount := 0
	failMessages := make([]string, 0)

	tx := db.DB.Begin()

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}

		account := row[0]
		password := row[1]

		if account == "" || password == "" {
			continue
		}

		merchsID := int32(0)
		if len(row) > 2 && row[2] != "" {
			if v, err := strconv.Atoi(row[2]); err == nil {
				merchsID = int32(v)
			}
		}

		if merchsID == 0 {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 商家ID不能为空", i+1))
			continue
		}

		var existCount int64
		if err := tx.Model(&model.SubMerch{}).Where("merchs_id = ? AND account = ?", merchsID, account).Count(&existCount).Error; err != nil {
			tx.Rollback()
			response.Fail(c, 500, "检查子商户账号失败: "+err.Error())
			return
		}

		if existCount > 0 {
			skipCount++
			continue
		}

		role := "0"
		if len(row) > 3 && row[3] != "" {
			role = row[3]
		}

		status := "0"
		if len(row) > 4 && row[4] != "" {
			status = row[4]
		}

		phone := ""
		if len(row) > 5 && row[5] != "" {
			phone = row[5]
		}

		email := ""
		if len(row) > 6 && row[6] != "" {
			email = row[6]
		}

		submerch := &model.SubMerch{
			MerchsID: merchsID,
			Account:  account,
			Password: password,
			Role:     role,
			Status:   status,
			Phone:    phone,
			Email:    email,
		}

		if err := tx.Create(submerch).Error; err != nil {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 创建失败 - %s", i+1, err.Error()))
			continue
		}

		successCount++
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, 500, "提交事务失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "导入完成", gin.H{
		"success": successCount,
		"skip":    skipCount,
		"fail":    failCount,
		"errors":  failMessages,
	})
}

func AdminImportRoomTag(c *gin.Context) {
	xlsx, err := parseExcelFile(c)
	if err != nil {
		return
	}

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		response.Fail(c, 500, "读取Excel行失败: "+err.Error())
		return
	}

	if len(rows) <= 1 {
		response.Fail(c, 400, "Excel文件为空或只有表头")
		return
	}

	successCount := 0
	skipCount := 0
	failCount := 0
	failMessages := make([]string, 0)

	tx := db.DB.Begin()

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}

		name := row[0]
		if name == "" {
			continue
		}

		merchsID := int32(0)
		if len(row) > 1 && row[1] != "" {
			if v, err := strconv.Atoi(row[1]); err == nil {
				merchsID = int32(v)
			}
		}

		if merchsID == 0 {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 商家ID不能为空", i+1))
			continue
		}

		var existCount int64
		if err := tx.Model(&model.RoomTag{}).Where("merchs_id = ? AND name = ?", merchsID, name).Count(&existCount).Error; err != nil {
			tx.Rollback()
			response.Fail(c, 500, "检查标签名称失败: "+err.Error())
			return
		}

		if existCount > 0 {
			skipCount++
			continue
		}

		tag := &model.RoomTag{
			MerchsID: merchsID,
			Name:     name,
		}

		if err := tx.Create(tag).Error; err != nil {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 创建失败 - %s", i+1, err.Error()))
			continue
		}

		successCount++
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, 500, "提交事务失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "导入完成", gin.H{
		"success": successCount,
		"skip":    skipCount,
		"fail":    failCount,
		"errors":  failMessages,
	})
}

func AdminImportRoomImage(c *gin.Context) {
	xlsx, err := parseExcelFile(c)
	if err != nil {
		return
	}

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		response.Fail(c, 500, "读取Excel行失败: "+err.Error())
		return
	}

	if len(rows) <= 1 {
		response.Fail(c, 400, "Excel文件为空或只有表头")
		return
	}

	successCount := 0
	skipCount := 0
	failCount := 0
	failMessages := make([]string, 0)

	tx := db.DB.Begin()

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}

		name := row[0]
		image := row[1]

		if name == "" || image == "" {
			continue
		}

		var existCount int64
		if err := tx.Model(&model.RoomImage{}).Where("name = ?", name).Count(&existCount).Error; err != nil {
			tx.Rollback()
			response.Fail(c, 500, "检查图片名称失败: "+err.Error())
			return
		}

		if existCount > 0 {
			skipCount++
			continue
		}

		roomImage := &model.RoomImage{
			Name:  name,
			Image: image,
		}

		if err := tx.Create(roomImage).Error; err != nil {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 创建失败 - %s", i+1, err.Error()))
			continue
		}

		successCount++
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, 500, "提交事务失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "导入完成", gin.H{
		"success": successCount,
		"skip":    skipCount,
		"fail":    failCount,
		"errors":  failMessages,
	})
}

func AdminImportHuifuAccount(c *gin.Context) {
	xlsx, err := parseExcelFile(c)
	if err != nil {
		return
	}

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		response.Fail(c, 500, "读取Excel行失败: "+err.Error())
		return
	}

	if len(rows) <= 1 {
		response.Fail(c, 400, "Excel文件为空或只有表头")
		return
	}

	successCount := 0
	skipCount := 0
	failCount := 0
	failMessages := make([]string, 0)

	tx := db.DB.Begin()

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 3 {
			continue
		}

		merchsID := int32(0)
		if row[0] != "" {
			if v, err := strconv.Atoi(row[0]); err == nil {
				merchsID = int32(v)
			}
		}

		if merchsID == 0 {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 商家ID不能为空", i+1))
			continue
		}

		accountType := row[1]
		account := row[2]

		if accountType == "" || account == "" {
			continue
		}

		name := ""
		if len(row) > 3 && row[3] != "" {
			name = row[3]
		}

		phone := ""
		if len(row) > 4 && row[4] != "" {
			phone = row[4]
		}

		huifuAccount := &model.HuifuAccount{
			MerchsID: merchsID,
			Type:     accountType,
			Account:  account,
			Name:     name,
			Phone:    phone,
			Choose:   "0",
		}

		if len(row) > 5 && row[5] != "" {
			huifuAccount.Identity = row[5]
		}

		if len(row) > 6 && row[6] != "" {
			huifuAccount.Card = row[6]
		}

		if len(row) > 7 && row[7] != "" {
			huifuAccount.Storename = row[7]
		}

		if len(row) > 8 && row[8] != "" {
			huifuAccount.Sharing = row[8]
		}

		if len(row) > 9 && row[9] != "" {
			huifuAccount.Share = row[9]
		}

		if len(row) > 10 && row[10] != "" {
			if v, err := strconv.ParseFloat(row[10], 64); err == nil {
				huifuAccount.Rate = v
			}
		}

		if err := tx.Create(huifuAccount).Error; err != nil {
			failCount++
			failMessages = append(failMessages, fmt.Sprintf("第%d行: 创建失败 - %s", i+1, err.Error()))
			continue
		}

		successCount++
	}

	if err := tx.Commit().Error; err != nil {
		response.Fail(c, 500, "提交事务失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "导入完成", gin.H{
		"success": successCount,
		"skip":    skipCount,
		"fail":    failCount,
		"errors":  failMessages,
	})
}