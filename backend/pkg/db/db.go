package db

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"centraliz-backend/model"
	"centraliz-backend/pkg/config"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// 使用Viper配置的MySQL连接字符串
	DB, err = gorm.Open(mysql.Open(config.GetMySQLDSN()), &gorm.Config{})
	if err != nil {
		panic("连接MySQL数据库失败: " + err.Error())
	}

	// 配置数据库连接池
	sqlDB, err := DB.DB()
	if err != nil {
		panic("获取数据库连接池失败: " + err.Error())
	}

	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(config.AppConfig.Database.MaxOpenConns)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(config.AppConfig.Database.MaxIdleConns)
	// 设置连接最大存活时间
	sqlDB.SetConnMaxLifetime(time.Duration(config.AppConfig.Database.ConnMaxLifetime) * time.Second)

	// 自动迁移表结构
	DB.AutoMigrate(
		&model.User{},
		&model.Device{},
		&model.Room{},
		&model.Order{},
	)
}