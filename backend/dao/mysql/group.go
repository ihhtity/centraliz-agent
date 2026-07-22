package mysql

import "centraliz-backend/model"

func GetGroupCount() (int64, error) {
	var count int64
	err := DB.Model(&model.Group{}).Count(&count).Error
	return count, err
}

func GetGroupListFiltered(merchsID int32, name string, groupType string, page, pageSize int) ([]model.Group, int64, error) {
	var groups []model.Group
	var total int64
	db := DB.Model(&model.Group{}).Order("id ASC")

	if merchsID > 0 {
		db = db.Where("merchs_id = ?", merchsID)
	}
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if groupType != "" {
		db = db.Where("type = ?", groupType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		db = DB.Order("id ASC")
		if merchsID > 0 {
			db = db.Where("merchs_id = ?", merchsID)
		}
		if name != "" {
			db = db.Where("name LIKE ?", "%"+name+"%")
		}
		if groupType != "" {
			db = db.Where("type = ?", groupType)
		}
		db = db.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	err := db.Find(&groups).Error
	return groups, total, err
}

func GetGroupByID(id uint64) (*model.Group, error) {
	var group model.Group
	err := DB.Where("id = ?", id).First(&group).Error
	return &group, err
}

func CreateGroup(group *model.Group) error {
	return DB.Create(group).Error
}

func UpdateGroup(group *model.Group) error {
	return DB.Save(group).Error
}

func DeleteGroup(id uint64) error {
	return DB.Delete(&model.Group{}, id).Error
}

func BatchDeleteGroup(ids []uint64) error {
	return DB.Where("id IN (?)", ids).Delete(&model.Group{}).Error
}

func BatchUpdateGroup(reqs []struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Location    string `json:"location"`
	Phone       string `json:"phone"`
	BindNumber  string `json:"bind_number"`
	ConsumePush string `json:"consume_push"`
}) error {
	for _, req := range reqs {
		group, err := GetGroupByID(req.ID)
		if err != nil {
			return err
		}
		if req.Name != "" {
			group.Name = req.Name
		}
		if req.Type != "" {
			group.Type = req.Type
		}
		if req.Location != "" {
			group.Location = req.Location
		}
		if req.Phone != "" {
			group.Phone = req.Phone
		}
		if req.BindNumber != "" {
			group.BindNumber = req.BindNumber
		}
		if req.ConsumePush != "" {
			group.ConsumePush = req.ConsumePush
		}
		if err := UpdateGroup(group); err != nil {
			return err
		}
	}
	return nil
}