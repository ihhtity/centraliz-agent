package model

import "time"

type Device struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	Name      string    `json:"name" gorm:"size:100;not null;comment:设备名称"`
	MerchsID  int32     `json:"merchsId" gorm:"column:merchs_id;not null;comment:商家外键"`
	RoomsID   *int32    `json:"roomsId" gorm:"column:rooms_id;comment:房间外键"`
	GroupsID  *int32    `json:"groupsId" gorm:"column:groups_id;comment:分组外键"`
	OrdersID  *int32    `json:"ordersId" gorm:"column:orders_id;comment:订单外键"`
	Status    *string   `json:"status" gorm:"size:20;default:'0';comment:设备状态：0空闲 1租用 2维修"`
	Type      *string   `json:"type" gorm:"size:10;comment:设备类型"`
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`
}