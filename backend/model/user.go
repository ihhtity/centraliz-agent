package model

import "time"

type User struct {
	ID        uint32     `json:"id" gorm:"size:10;primaryKey;autoIncrement;comment:主键"`           // 主键
	MerchsID  int32      `json:"merchsId" gorm:"size:10;column:merchs_id;default:0;comment:商家外键"` // 商家外键
	RoomsID   *int32     `json:"roomsId" gorm:"size:10;column:rooms_id;default:0;comment:房间外键"`   // 房间外键
	OrdersID  *int32     `json:"ordersId" gorm:"size:20;column:orders_id;default:0;comment:订单外键"` // 订单外键
	Name      string     `json:"name" gorm:"size:64;not null;comment:用户名"`                        // 用户名
	Account   string     `json:"account" gorm:"size:64;not null;comment:账号"`                      // 账号
	Password  string     `json:"password" gorm:"size:255;not null;comment:密码"`                    // 密码
	Email     *string    `json:"email" gorm:"size:64;comment:邮箱"`                                 // 邮箱
	Phone     *string    `json:"phone" gorm:"size:20;comment:手机号"`                                // 手机号
	Privacy   string     `json:"privacy" gorm:"size:5;default:'0';comment:隐私政策 0拒绝 1同意"`          // 隐私政策 0拒绝 1同意
	Status    *string    `json:"status" gorm:"size:5;default:'0';comment:状态 0白名单 1黑名单"`           // 0白名单 1黑名单
	AvatarURL *string    `json:"avatarURL" gorm:"size:255;comment:头像URL"`                         // 头像URL
	UnionID   string     `json:"unionId" gorm:"size:100;index;comment:微信UnionID(跨端唯一标识)"`         // 微信UnionID
	OpenID    string     `json:"openid" gorm:"size:100;index;comment:微信OpenID(用户唯一标识)"`           // 微信小程序OpenID(用户唯一标识)
	GOpenID   string     `json:"gopenid" gorm:"size:100;index;comment:微信GOpenID(用户唯一标识)"`         // 微信公众号(用户唯一标识)
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`     // 更新时间
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`     // 创建时间
}
