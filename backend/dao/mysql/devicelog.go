package mysql

import "centraliz-backend/model"

func GetDeviceLogListFiltered(merchsID, devicesID, roomID int32, code, deviceType, control, status string, page, pageSize int) ([]model.Devicelog, int64, error) {
	var logs []model.Devicelog
	var total int64
	db := DB.Model(&model.Devicelog{}).Order("id DESC")

	if merchsID > 0 {
		db = db.Where("merchs_id = ?", merchsID)
	}
	if devicesID > 0 {
		db = db.Where("devices_id = ?", devicesID)
	}
	if roomID > 0 {
		db = db.Where("room_id = ?", roomID)
	}
	if code != "" {
		db = db.Where("code LIKE ?", "%"+code+"%")
	}
	if deviceType != "" {
		db = db.Where("type = ?", deviceType)
	}
	if control != "" {
		db = db.Where("control = ?", control)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		db = DB.Order("id DESC")
		if merchsID > 0 {
			db = db.Where("merchs_id = ?", merchsID)
		}
		if devicesID > 0 {
			db = db.Where("devices_id = ?", devicesID)
		}
		if roomID > 0 {
			db = db.Where("room_id = ?", roomID)
		}
		if code != "" {
			db = db.Where("code LIKE ?", "%"+code+"%")
		}
		if deviceType != "" {
			db = db.Where("type = ?", deviceType)
		}
		if control != "" {
			db = db.Where("control = ?", control)
		}
		if status != "" {
			db = db.Where("status = ?", status)
		}
		db = db.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	err := db.Find(&logs).Error
	return logs, total, err
}

func GetDeviceLogByID(id uint64) (*model.Devicelog, error) {
	var log model.Devicelog
	err := DB.Where("id = ?", id).First(&log).Error
	return &log, err
}

func DeleteDeviceLog(id uint64) error {
	return DB.Delete(&model.Devicelog{}, id).Error
}

func BatchDeleteDeviceLog(ids []uint64) error {
	return DB.Where("id IN (?)", ids).Delete(&model.Devicelog{}).Error
}
