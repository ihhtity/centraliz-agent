package model

import "time"

type HuifuAccount struct {
	ID        int32      `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`          // 主键
	MerchID   *int32     `json:"merchId" gorm:"index;comment:商家id"`                      // 商家id
	Code      *string    `json:"code" gorm:"size:255;comment:汇付编码"`                      // 汇付编码
	Sharing   string     `json:"sharing" gorm:"not null;comment:多方分账"`                   // 多方分账，text类型
	Account   *string    `json:"account" gorm:"size:255;comment:账号"`                     // 账号
	Phone     *string    `json:"phone" gorm:"size:255;comment:手机号"`                      // 手机号
	Name      *string    `json:"name" gorm:"size:255;comment:姓名"`                        // 姓名
	Identity  *string    `json:"identity" gorm:"size:255;comment:身份证"`                   // 身份证
	Card      *string    `json:"card" gorm:"size:255;comment:银行卡"`                       // 银行卡
	Encrypt   *string    `json:"encrypt" gorm:"size:255;comment:营业执照编码"`                 // 营业执照编码
	Storename *string    `json:"storename" gorm:"size:255;comment:店名"`                   // 店名
	Area      string     `json:"area" gorm:"not null;comment:经营地址"`                      // 经营地址
	Picture   string     `json:"picture" gorm:"not null;comment:店铺图片"`                   // 店铺图片
	Remarks   string     `json:"remarks" gorm:"not null;comment:使用场景描述"`                 // 使用场景描述，text类型
	Type      *string    `json:"type" gorm:"size:255;comment:账号类型"`                      // 账号类型
	Choose    *string    `json:"choose" gorm:"default:'0';comment:0未选择 1已选择'"`           // 0未选择 1已选择
	Share     string     `json:"share" gorm:"not null;default:'0';comment:0关闭分账 1开启分账'"` // 0关闭分账 1开启分账
	Rate      *float64   `json:"rate" gorm:"default:0.00;comment:分账比率"`                  // 分账比率
	CreatedAt *time.Time `json:"createdAt" gorm:"comment:创建时间"`                          // 创建时间
}
