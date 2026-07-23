package mysql

import "centraliz-backend/model"

func GetTotalUserCount() (int64, error) {
	var count int64
	err := GetDB().Model(&model.User{}).Count(&count).Error
	return count, err
}