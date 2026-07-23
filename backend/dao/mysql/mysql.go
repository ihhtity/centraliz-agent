package mysql

import (
	"centraliz-backend/pkg/db"

	"gorm.io/gorm"
)

func GetDB() *gorm.DB {
	return db.DB
}
