package mysql

import "centraliz-backend/model"

func GetDeviceCount() (int64, error) {
	var count int64
	err := DB.Model(&model.Device{}).Count(&count).Error
	return count, err
}

func GetDeviceListFiltered(merchsID int32, groupsID int32, name string, status string, deviceType string, page, pageSize int) ([]model.Device, int64, error) {
	var devices []model.Device
	var total int64
	db := DB.Model(&model.Device{}).Order("id ASC")

	if merchsID > 0 {
		db = db.Where("merchs_id = ?", merchsID)
	}
	if groupsID > 0 {
		db = db.Where("groups_id = ?", groupsID)
	}
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if deviceType != "" {
		db = db.Where("type = ?", deviceType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		db = DB.Order("id ASC")
		if merchsID > 0 {
			db = db.Where("merchs_id = ?", merchsID)
		}
		if groupsID > 0 {
			db = db.Where("groups_id = ?", groupsID)
		}
		if name != "" {
			db = db.Where("name LIKE ?", "%"+name+"%")
		}
		if status != "" {
			db = db.Where("status = ?", status)
		}
		if deviceType != "" {
			db = db.Where("type = ?", deviceType)
		}
		db = db.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	err := db.Find(&devices).Error
	return devices, total, err
}

func GetDeviceByID(id uint64) (*model.Device, error) {
	var device model.Device
	err := DB.Where("id = ?", id).First(&device).Error
	return &device, err
}

func CreateDevice(device *model.Device) error {
	return DB.Create(device).Error
}

func UpdateDevice(device *model.Device) error {
	return DB.Save(device).Error
}

func DeleteDevice(id uint64) error {
	return DB.Delete(&model.Device{}, id).Error
}

func BatchDeleteDevice(ids []uint64) error {
	return DB.Where("id IN (?)", ids).Delete(&model.Device{}).Error
}

func BatchUpdateDevice(reqs []struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	LockCount int32  `json:"lock_count"`
}) error {
	for _, req := range reqs {
		device, err := GetDeviceByID(req.ID)
		if err != nil {
			return err
		}
		if req.Name != "" {
			device.Name = req.Name
		}
		if req.Code != "" {
			device.Code = req.Code
		}
		if req.Status != "" {
			device.Status = req.Status
		}
		if req.Type != "" {
			device.Type = req.Type
		}
		if req.LockCount >= 0 {
			device.LockCount = req.LockCount
		}
		if err := UpdateDevice(device); err != nil {
			return err
		}
	}
	return nil
}