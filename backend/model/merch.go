package model

import "time"

type Merch struct {
	ID        uint32     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`               // 主键，int(11) unsigned
	Account   string     `json:"account" gorm:"size:50;not null;comment:账号"`                  // 账号
	Password  string     `json:"password" gorm:"size:255;not null;comment:密码"`                // 密码
	Email     *string    `json:"email" gorm:"size:100;comment:邮箱"`                            // 邮箱
	Phone     *string    `json:"phone" gorm:"size:20;comment:手机号"`                            // 手机号
	Role      *string    `json:"role" gorm:"size:20;default:'0';comment:角色 0商家 1管理者 2代理商"`    // 角色 0商家 1管理者 2代理商
	Status    *string    `json:"status" gorm:"size:20;default:'0';comment:状态 0白名单 1黑名单"`      // 0白名单 1黑名单
	LogAt     *time.Time `json:"logAt" gorm:"column:log_at;precision:3;comment:上次登录时间"`       // 上次登录时间
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"` // 创建时间
}
