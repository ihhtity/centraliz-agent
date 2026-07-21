package model

import "time"

type RoomImage struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	Name      string    `json:"name" gorm:"size:255;comment:房间图片名称"`
	Image     string    `json:"image" gorm:"type:text;comment:房间图片地址"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`
}
