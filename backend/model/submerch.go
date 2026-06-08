package model

import "time"

type SubMerch struct {
	ID        int32      `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`    // 主键
	Account   string     `json:"account" gorm:"not null;comment:账号"`               // 账号
	Password  string     `json:"password" gorm:"not null;comment:密码"`              // 密码
	MerchsID  int32      `json:"merchsId" gorm:"not null;comment:商家外键"`            // 商家外键
	Email     *string    `json:"email" gorm:"size:255;comment:邮箱"`                 // 邮箱
	Phone     *string    `json:"phone" gorm:"size:255;comment:手机号"`                // 手机号
	Role      *string    `json:"role" gorm:"default:'0';comment:角色 0商家 1管理者 2代理商"` // 角色 0商家 1管理者 2代理商
	Status    *string    `json:"status" gorm:"default:'0';comment:状态 0白名单 1黑名单"`   // 0白名单 1黑名单
	Rule      string     `json:"rule" gorm:"not null;comment:使用权限"`                // 使用权限，text类型
	LogAt     *time.Time `json:"logAt" gorm:"comment:上次登录时间"`                      // 上次登录时间
	CreatedAt *time.Time `json:"createdAt" gorm:"comment:创建时间"`                    // 创建时间
}
