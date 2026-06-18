package model

import "time"

type Room struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	MerchsID  int32     `json:"merchsId" gorm:"column:merchs_id;default:0;comment:商家外键"`
	UsersID   int32     `json:"usersId" gorm:"column:users_id;default:0;comment:用户外键"`
	GroupsID  int32     `json:"groupsId" gorm:"column:groups_id;default:0;comment:分组外键"`
	RulesID   int32     `json:"rulesId" gorm:"column:rules_id;default:0;comment:规则外键"`
	DevicesID string    `json:"devicesId" gorm:"column:devices_id;default:0;comment:设备外键"`
	OrdersID  int32     `json:"ordersId" gorm:"column:orders_id;default:0;comment:订单外键"`
	Name      string    `json:"name" gorm:"size:255;comment:房间名称"`
	Tag       string    `json:"tag" gorm:"size:255;comment:房间标签"`
	Status    string    `json:"status" gorm:"size:191;default:'空闲';comment:空闲 租用 维修"`
	BoardNo   string    `json:"boardNo" gorm:"size:100;default:'01';comment:板号"`
	LockNo    string    `json:"lockNo" gorm:"size:100;default:'00';comment:锁号"`
	FreeTime  time.Time `json:"freeTime" gorm:"column:free_time;precision:3;comment:免费时间"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`
}
