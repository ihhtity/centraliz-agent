package mysql

import "centraliz-backend/model"

func GetMerchListFiltered(account, phone, role, status string, page, pageSize int) ([]model.Merch, int64, error) {
	var merchs []model.Merch
	var total int64
	db := DB.Model(&model.Merch{}).Order("id ASC")

	if account != "" {
		db = db.Where("account LIKE ?", "%"+account+"%")
	}
	if phone != "" {
		db = db.Where("phone LIKE ?", "%"+phone+"%")
	}
	if role != "" {
		db = db.Where("role = ?", role)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		db = DB.Order("id ASC")
		if account != "" {
			db = db.Where("account LIKE ?", "%"+account+"%")
		}
		if phone != "" {
			db = db.Where("phone LIKE ?", "%"+phone+"%")
		}
		if role != "" {
			db = db.Where("role = ?", role)
		}
		if status != "" {
			db = db.Where("status = ?", status)
		}
		db = db.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	err := db.Find(&merchs).Error
	return merchs, total, err
}

func GetMerchByID(id uint64) (*model.Merch, error) {
	var merch model.Merch
	err := DB.Where("id = ?", id).First(&merch).Error
	return &merch, err
}

func CreateMerch(merch *model.Merch) error {
	return DB.Create(merch).Error
}

func UpdateMerch(merch *model.Merch) error {
	return DB.Save(merch).Error
}

func DeleteMerch(id uint64) error {
	return DB.Delete(&model.Merch{}, id).Error
}

func BatchDeleteMerch(ids []uint64) error {
	return DB.Where("id IN (?)", ids).Delete(&model.Merch{}).Error
}

func BatchUpdateMerch(reqs []struct {
	ID       uint64 `json:"id"`
	Account  string `json:"account"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}) error {
	for _, req := range reqs {
		merch, err := GetMerchByID(req.ID)
		if err != nil {
			return err
		}
		if req.Account != "" {
			merch.Account = req.Account
		}
		if req.Email != "" {
			merch.Email = req.Email
		}
		if req.Phone != "" {
			merch.Phone = req.Phone
		}
		if req.Role != "" {
			merch.Role = req.Role
		}
		if req.Status != "" {
			merch.Status = req.Status
		}
		if err := UpdateMerch(merch); err != nil {
			return err
		}
	}
	return nil
}
