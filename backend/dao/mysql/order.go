package mysql

import "centraliz-backend/model"

func GetOrderCount() (int64, error) {
	var count int64
	err := GetDB().Model(&model.Order{}).Count(&count).Error
	return count, err
}

func GetOrderListFiltered(merchsID int32, usersID int32, roomsID int32, status string, orderCode string, orderNo string, userPhone string, payType string, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64
	db := GetDB().Model(&model.Order{}).Order("id DESC")

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
	if orderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+orderNo+"%")
	}
	if userPhone != "" {
		db = db.Where("user_phone LIKE ?", "%"+userPhone+"%")
	}
	if payType != "" {
		db = db.Where("pay_type = ?", payType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	err := db.Find(&orders).Error
	return orders, total, err
}

func GetOrderByID(id uint64) (*model.Order, error) {
	var order model.Order
	err := GetDB().Where("id = ?", id).First(&order).Error
	return &order, err
}

func CreateOrder(order *model.Order) error {
	return GetDB().Create(order).Error
}

func UpdateOrder(order *model.Order) error {
	return GetDB().Save(order).Error
}

func DeleteOrder(id uint64) error {
	return GetDB().Delete(&model.Order{}, id).Error
}

func BatchDeleteOrder(ids []uint64) error {
	return GetDB().Where("id IN (?)", ids).Delete(&model.Order{}).Error
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

func BatchUpdateOrderByIDs(ids []uint64, data map[string]interface{}) error {
	return GetDB().Model(&model.Order{}).Where("id IN (?)", ids).Updates(data).Error
}

func GetTodayOrderCount() (int64, error) {
	var count int64
	err := GetDB().Model(&model.Order{}).Where("DATE(created_at) = CURDATE()").Count(&count).Error
	return count, err
}

func GetTodayRevenue() (float64, error) {
	var revenue float64
	err := GetDB().Model(&model.Order{}).Where("DATE(created_at) = CURDATE() AND status = 'paid'").Select("COALESCE(SUM(amount), 0)").Scan(&revenue).Error
	return revenue, err
}

func GetTotalRevenue() (float64, error) {
	var revenue float64
	err := GetDB().Model(&model.Order{}).Where("status = 'paid'").Select("COALESCE(SUM(amount), 0)").Scan(&revenue).Error
	return revenue, err
}

func GetOrderTrend(days int) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := GetDB().Model(&model.Order{}).
		Select("DATE(created_at) as date, COUNT(*) as value").
		Where("created_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)", days).
		Group("DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&result).Error
	return result, err
}

func GetRevenueTrend(days int) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := GetDB().Model(&model.Order{}).
		Select("DATE(created_at) as date, COALESCE(SUM(amount), 0) as value").
		Where("created_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY) AND status = 'paid'").
		Group("DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&result).Error
	return result, err
}
