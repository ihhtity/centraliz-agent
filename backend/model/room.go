package model

import "time"

type Room struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	Name      *string   `json:"name" gorm:"size:255;comment:房间名称"`
	MerchsID  int32     `json:"merchsId" gorm:"column:merchs_id;not null;default:0;comment:商家外键"`
	GroupsID  *int32    `json:"groupsId" gorm:"column:groups_id;default:0;comment:分组外键"`
	RulesID   *int32    `json:"rulesId" gorm:"column:rules_id;default:0;comment:规则外键"`
	Tag       *string   `json:"tag" gorm:"size:255;comment:房间标签"`
	Status    *string   `json:"status" gorm:"size:191;default:'0';comment:0空闲 1租用 2维修"`
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`
}