package model

import "time"

type Code struct {
	ID        uint32    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`                                      // 主键
	Phone     *string   `json:"phone" gorm:"size:20;index;comment:手机号"`                                             // 手机号
	Email     *string   `json:"email" gorm:"size:100;index;comment:邮箱"`                                             // 邮箱
	Code      string    `json:"code" gorm:"size:10;not null;index;comment:验证码"`                                     // 验证码
	Type      int8      `json:"type" gorm:"not null;default:1;comment:1-注册, 2-找回密码, 3-绑定手机, 4-绑定邮箱"`                // 1-注册, 2-找回密码, 3-绑定手机, 4-绑定邮箱
	Status    int8      `json:"status" gorm:"not null;default:0;comment:状态: 0-未使用，1-已使用，2-已过期"`                     // 状态：0-未使用，1-已使用，2-已过期
	ExpireAt  time.Time `json:"expireAt" gorm:"column:expire_at;not null;index;comment:过期时间"`                       // 过期时间
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;comment:更新时间"` // 更新时间
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;comment:创建时间"` // 创建时间
}
