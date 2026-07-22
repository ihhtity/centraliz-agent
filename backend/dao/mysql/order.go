package mysql

import "centraliz-backend/model"

func GetOrderCount() (int64, error) {
	var count int64
	err := DB.Model(&model.Order{}).Count(&count).Error
	return count, err
}

func GetOrderListFiltered(merchsID int32, usersID int32, roomsID int32, status string, orderCode string, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64
	db := DB.Model(&model.Order{}).Order("id DESC")

	if merchsID > 0 {
		db = db.Where("merchs_id = ?", merchsID)
	}
	if usersID > 0 {
		db = db.Where("users_id = ?", usersID)
	}
	if roomsID > 0 {
		db = db.Where("rooms_id = ?", roomsID)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if orderCode != "" {
		db = db.Where("code LIKE ?", "%"+orderCode+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		db = DB.Order("id DESC")
		if merchsID > 0 {
			db = db.Where("merchs_id = ?", merchsID)
		}
		if usersID > 0 {
			db = db.Where("users_id = ?", usersID)
		}
		if roomsID > 0 {
			db = db.Where("rooms_id = ?", roomsID)
		}
		if status != "" {
			db = db.Where("status = ?", status)
		}
		if orderCode != "" {
			db = db.Where("code LIKE ?", "%"+orderCode+"%")
		}
		db = db.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	err := db.Find(&orders).Error
	return orders, total, err
}

func GetOrderByID(id uint64) (*model.Order, error) {
	var order model.Order
	err := DB.Where("id = ?", id).First(&order).Error
	return &order, err
}

func CreateOrder(order *model.Order) error {
	return DB.Create(order).Error
}

func UpdateOrder(order *model.Order) error {
	return DB.Save(order).Error
}

func DeleteOrder(id uint64) error {
	return DB.Delete(&model.Order{}, id).Error
}

func BatchDeleteOrder(ids []uint64) error {
	return DB.Where("id IN (?)", ids).Delete(&model.Order{}).Error
}

func BatchUpdateOrder(reqs []struct {
	ID     uint64 `json:"id"`
	Status string `json:"status"`
	Remark string `json:"remark"`
}) error {
	for _, req := range reqs {
		order, err := GetOrderByID(req.ID)
		if err != nil {
			return err
		}
		if req.Status != "" {
			order.Status = req.Status
		}
		if req.Remark != "" {
			order.Remark = req.Remark
		}
		if err := UpdateOrder(order); err != nil {
			return err
		}
	}
	return nil
}