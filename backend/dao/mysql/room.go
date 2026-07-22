package mysql

import "centraliz-backend/model"

func GetRoomCount() (int64, error) {
	var count int64
	err := DB.Model(&model.Room{}).Count(&count).Error
	return count, err
}

func GetRoomListFiltered(merchsID int32, groupsID int32, name string, status string, page, pageSize int) ([]model.Room, int64, error) {
	var rooms []model.Room
	var total int64
	db := DB.Model(&model.Room{}).Order("id ASC")

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
		db = db.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	err := db.Find(&rooms).Error
	return rooms, total, err
}

func GetRoomByID(id uint64) (*model.Room, error) {
	var room model.Room
	err := DB.Where("id = ?", id).First(&room).Error
	return &room, err
}

func CreateRoom(room *model.Room) error {
	return DB.Create(room).Error
}

func UpdateRoom(room *model.Room) error {
	return DB.Save(room).Error
}

func DeleteRoom(id uint64) error {
	return DB.Delete(&model.Room{}, id).Error
}

func BatchDeleteRoom(ids []uint64) error {
	return DB.Where("id IN (?)", ids).Delete(&model.Room{}).Error
}

func BatchUpdateRoom(reqs []struct {
	ID       uint64  `json:"id"`
	Name     string  `json:"name"`
	Tag      string  `json:"tag"`
	Status   string  `json:"status"`
	LockNo   string  `json:"lock_no"`
	Price    float32 `json:"price"`
	Image    string  `json:"image"`
}) error {
	for _, req := range reqs {
		room, err := GetRoomByID(req.ID)
		if err != nil {
			return err
		}
		if req.Name != "" {
			room.Name = req.Name
		}
		if req.Tag != "" {
			room.Tag = req.Tag
		}
		if req.Status != "" {
			room.Status = req.Status
		}
		if req.LockNo != "" {
			room.LockNo = req.LockNo
		}
		if req.Price > 0 {
			room.Price = req.Price
		}
		if req.Image != "" {
			room.Image = req.Image
		}
		if err := UpdateRoom(room); err != nil {
			return err
		}
	}
	return nil
}