package model

import "time"

type Group struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	Name        string    `json:"name" gorm:"size:100;comment:房间名称"`
	MerchsID    int32     `json:"merchsId" gorm:"column:merchs_id;default:0;comment:商家外键"`
	RulesID     int32     `json:"rulesId" gorm:"column:rules_id;default:0;comment:规则外键"`
	DevicesID   int32     `json:"devicesId" gorm:"column:devices_id;default:0;comment:设备外键"`
	Phone       string    `json:"phone" gorm:"size:20;comment:客服手机号"`
	Count       int32     `json:"count" gorm:"index;default:0;comment:房间数量"`
	Type        string    `json:"type" gorm:"size:5;default:'存柜';comment:存柜 零售"`
	Location    string    `json:"location" gorm:"size:255;comment:房间位置"`
	RuleName    string    `json:"ruleName" gorm:"size:100;comment:规则名称"`
	BindNumber  string    `json:"bindNumber" gorm:"size:10;default:'关闭';comment:绑定号码 关闭 手动 自动"`
	ConsumePush string    `json:"consumePush" gorm:"size:10;default:'关闭';comment:消费推送 关闭 开启"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`
}
