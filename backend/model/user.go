package model

import "time"

type User struct {
	ID        uint32     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`                    // 主键
	Name      string     `json:"name" gorm:"not null;comment:用户名"`                                 // 用户名
	Account   string     `json:"account" gorm:"not null;comment:账号"`                               // 账号
	Password  string     `json:"password" gorm:"not null;comment:密码"`                              // 密码
	Email     *string    `json:"email" gorm:"size:255;comment:邮箱"`                                 // 邮箱
	Phone     *string    `json:"phone" gorm:"size:255;comment:手机号"`                                // 手机号
	MerchsID  int32      `json:"merchsId" gorm:"column:merchs_id;not null;default:0;comment:商家外键"` // 商家外键
	RoomsID   *int32     `json:"roomsId" gorm:"column:rooms_id;index;comment:房间外键"`                // 房间外键
	OrdersID  *int32     `json:"ordersId" gorm:"column:orders_id;index;comment:订单外键"`              // 订单外键
	Status    *string    `json:"status" gorm:"size:191;default:'0';comment:状态 0白名单 1黑名单"`          // 0白名单 1黑名单
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`      // 更新时间
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`      // 创建时间
}
