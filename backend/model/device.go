package model

import "time"

type Device struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	MerchsID    int32     `json:"merchsId" gorm:"default:0;column:merchs_id;comment:商家外键"`
	RoomsID     int32     `json:"roomsId" gorm:"default:0;column:rooms_id;comment:房间外键"`
	GroupsID    int32     `json:"groupsId" gorm:"default:0;column:groups_id;comment:分组外键"`
	OrdersID    int32     `json:"ordersId" gorm:"default:0;column:orders_id;comment:订单外键"`
	Name        string    `json:"name" gorm:"size:100;not null;comment:设备名称"`
	Code        string    `json:"code" gorm:"size:100;comment:设备编码"`
	Recharge    string    `json:"recharge" gorm:"default:'';comment:充值编码"`
	LockCount   int32     `json:"lockCount" gorm:"default:0;comment:锁定数量"`
	BoardNo     string    `json:"boardNo" gorm:"size:100;default:'01';comment:板号"`
	Version     string    `json:"version" gorm:"default:'';comment:设备版本"`
	Status      string    `json:"status" gorm:"size:64;default:'在线';comment:设备状态：在线 离线 维修"`
	Type        string    `json:"type" gorm:"size:64;default:'集控';comment:设备类型：集控 摄像头"`
	Signal      string    `json:"signal" gorm:"size:64;default:'100';comment:信号强度"`
	Heat        string    `json:"heat" gorm:"size:64;default:'25';comment:当前温度"`
	ProtectHeat string    `json:"protectHeat" gorm:"size:64;default:'85';comment:保护温度"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`
}
