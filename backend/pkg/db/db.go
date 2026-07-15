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

	// 使用Viper配置的MySQL连接字符串
	DB, err = gorm.Open(mysql.Open(config.GetMySQLDSN()), &gorm.Config{})
	if err != nil {
		return err
	}

	// 配置数据库连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(config.AppConfig.Database.MaxOpenConns)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(config.AppConfig.Database.MaxIdleConns)
	// 设置连接最大存活时间
	sqlDB.SetConnMaxLifetime(time.Duration(config.AppConfig.Database.ConnMaxLifetime) * time.Second)

	// 设置表选项，强制使用InnoDB引擎
	// DB = DB.Set("gorm:table_options", "ENGINE=InnoDB")

	// 自动迁移表结构
	err = DB.AutoMigrate(
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
	)
	if err != nil {
		return err
	}

	return nil
}
