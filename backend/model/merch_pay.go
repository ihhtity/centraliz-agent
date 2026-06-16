package model

import "time"

type MerchPay struct {
	ID            int32     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`                // 主键
	Code          string    `json:"code" gorm:"size:64;uniqueIndex;comment:订单号"`                  // 订单号
	MerchsID      int32     `json:"merchsId" gorm:"index;default:0;comment:商家外键"`                 // 商家外键
	Name          string    `json:"name" gorm:"size:64;index;comment:商品名称"`                       // 商品名称
	ReqDate       string    `json:"reqDate" gorm:"size:64;comment:汇付支付时间"`                        // 汇付支付时间
	HfSeqId       string    `json:"hfSeqId" gorm:"size:64;comment:汇付订单号"`                         // 汇付订单号
	OriginalPrice float64   `json:"originalPrice" gorm:"type:decimal(10,2);comment:订单原价"`         // 订单原价
	Price         float64   `json:"price" gorm:"type:decimal(10,2);comment:实际支付金额"`               // 实际支付金额
	Locktotal     int32     `json:"locktotal" gorm:"default:0;comment:锁总数"`                       // 锁总数
	Type          string    `json:"type" gorm:"size:10;default:'短信';comment:订单类型 短信"`             // 订单类型 短信
	Status        string    `json:"status" gorm:"size:10;default:'未完成';comment:订单状态 未完成 已完成 已关闭"` // 订单状态 未完成 已完成 已关闭
	Remarks       string    `json:"remarks" gorm:"type:text;comment:订单备注"`                        // 订单备注
	CreatedAt     time.Time `json:"createdAt" gorm:"comment:创建时间"`                                // 创建时间
}
