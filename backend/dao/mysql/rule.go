package mysql

import "centraliz-backend/model"

func GetRuleListFiltered(merchsID int32, name string, ruleType string, page, pageSize int) ([]model.Rule, int64, error) {
	var rules []model.Rule
	var total int64
	db := DB.Model(&model.Rule{}).Order("id ASC")

	if merchsID > 0 {
		db = db.Where("merchs_id = ?", merchsID)
	}
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if ruleType != "" {
		db = db.Where("type = ?", ruleType)
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
		if ruleType != "" {
			db = db.Where("type = ?", ruleType)
		}
		db = db.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	err := db.Find(&rules).Error
	return rules, total, err
}

func GetRuleByID(id uint32) (*model.Rule, error) {
	var rule model.Rule
	err := DB.Where("id = ?", id).First(&rule).Error
	return &rule, err
}

func CreateRule(rule *model.Rule) error {
	return DB.Create(rule).Error
}

func UpdateRule(rule *model.Rule) error {
	return DB.Save(rule).Error
}

func DeleteRule(id uint32) error {
	return DB.Delete(&model.Rule{}, id).Error
}

func BatchDeleteRule(ids []uint32) error {
	return DB.Where("id IN (?)", ids).Delete(&model.Rule{}).Error
}

func BatchUpdateRule(reqs []struct {
	ID          uint32  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Mode        string  `json:"mode"`
	Price       float32 `json:"price"`
	Deposit     float32 `json:"deposit"`
	Description string  `json:"description"`
}) error {
	for _, req := range reqs {
		rule, err := GetRuleByID(req.ID)
		if err != nil {
			return err
		}
		if req.Name != "" {
			rule.Name = req.Name
		}
		if req.Type != "" {
			rule.Type = req.Type
		}
		if req.Mode != "" {
			rule.Mode = req.Mode
		}
		if req.Price >= 0 {
			rule.Price = req.Price
		}
		if req.Deposit >= 0 {
			rule.Deposit = req.Deposit
		}
		if req.Description != "" {
			rule.Description = req.Description
		}
		if err := UpdateRule(rule); err != nil {
			return err
		}
	}
	return nil
}