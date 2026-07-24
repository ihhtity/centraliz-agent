package db

import (
	"time"

	"centraliz-backend/model"
	"centraliz-backend/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error

	DB, err = gorm.Open(mysql.Open(config.GetMySQLDSN()), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(config.AppConfig.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.AppConfig.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.AppConfig.Database.ConnMaxLifetime) * time.Second)

	err = DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").AutoMigrate(
		&model.User{},
		&model.Device{},
		&model.Room{},
		&model.Order{},
		&model.Merch{},
		&model.SubMerch{},
		&model.Group{},
		&model.HuifuAccount{},
		&model.MerchPay{},
		&model.WechatUser{},
		&model.Rule{},
		&model.Devicelog{},
		&model.RoomImage{},
		&model.RoomTag{},
		&model.AssistantSession{},
		&model.AssistantMessage{},
	)
	if err != nil {
		return err
	}

	return nil
}
