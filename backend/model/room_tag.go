package model

import "time"

type RoomTag struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	MerchsID  int32     `json:"merchsId" gorm:"column:merchs_id;default:0;comment:商家外键"`
	Name      string    `json:"name" gorm:"type:text;comment:标签名称"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`
}
