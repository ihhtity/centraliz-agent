package model

import "time"

type WechatUser struct {
	ID           uint32     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`                          // 主键
	OpenID       string     `json:"openId" gorm:"size:100;uniqueIndex;comment:微信小程序OpenID"`                 // 微信小程序OpenID
	GopenID      *string    `json:"gopenId" gorm:"size:100;uniqueIndex;comment:微信公众号GopenID"`               // 微信公众号GopenID
	UnionID      *string    `json:"unionId" gorm:"size:100;index;comment:微信UnionID"`                        // 微信UnionID
	SessionKey   *string    `json:"sessionKey" gorm:"size:100;comment:小程序会话密钥"`                             // 小程序会话密钥
	AccessToken  *string    `json:"accessToken" gorm:"size:500;comment:公众号访问令牌"`                            // 公众号访问令牌
	RefreshToken *string    `json:"refreshToken" gorm:"size:500;comment:刷新令牌"`                              // 刷新令牌
	Nickname     *string    `json:"nickname" gorm:"size:100;comment:微信昵称"`                                  // 微信昵称
	Avatar       *string    `json:"avatar" gorm:"size:500;comment:头像URL"`                                   // 头像URL
	Gender       *int       `json:"gender" gorm:"comment:性别 0未知 1男 2女"`                                     // 性别 0未知 1男 2女
	Country      *string    `json:"country" gorm:"size:50;comment:国家"`                                      // 国家
	Province     *string    `json:"province" gorm:"size:50;comment:省份"`                                     // 省份
	City         *string    `json:"city" gorm:"size:50;comment:城市"`                                         // 城市
	Language     *string    `json:"language" gorm:"size:20;comment:语言"`                                     // 语言
	Platform     string     `json:"platform" gorm:"size:20;not null;comment:平台类型 miniprogram/mp"`           // 平台类型 miniprogram小程序 mp公众号
	MerchsID     int32      `json:"merchsId" gorm:"column:merchs_id;index;not null;default:0;comment:商家外键"` // 商家外键
	UsersID      *int32     `json:"usersId" gorm:"column:users_id;index;comment:用户外键"`                      // 用户外键
	Status       *string    `json:"status" gorm:"size:20;default:'0';comment:状态 0正常 1禁用"`                   // 状态 0正常 1禁用
	ExpiredAt    *time.Time `json:"expiredAt" gorm:"column:expired_at;precision:3;comment:令牌过期时间"`          // 令牌过期时间
	CreatedAt    *time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`            // 创建时间
	UpdatedAt    *time.Time `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`            // 更新时间
}

func (WechatUser) TableName() string {
	return "wechat_users"
}
