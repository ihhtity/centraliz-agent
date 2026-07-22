package model

import (
	"fmt"
	"time"
)

// Rule 规则模型
type Rule struct {
	ID           uint32    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	MerchsID     int32     `json:"merchsId" gorm:"column:merchs_id;default:0;comment:商家外键"`
	Name         string    `json:"name" gorm:"size:100;default:'';comment:规则名称"`
	Type         string    `json:"type" gorm:"size:20;default:'free';comment:规则类型 free:免费模式 charge:收费模式"`
	Mode         string    `json:"mode" gorm:"size:50;default:single;comment:模式类型 single:单次开锁;deposit:一存一取;pay_single:单次付费;pay_deposit:先存后取;pay_hourly:按时付费;pay_time:预付费"`
	Price        float32   `json:"price" gorm:"type:decimal(10,2);default:0.00;comment:单价（元）"`
	Deposit      float32   `json:"deposit" gorm:"type:decimal(10,2);default:0.00;comment:押金（元）"`
	Rate         float32   `json:"rate" gorm:"type:decimal(5,2);default:0;comment:费率"`
	DurationUnit string    `json:"durationUnit" gorm:"size:20;default:hour;comment:时长单位 hour:小时 day:天 month:月 minute:分钟"`
	AutoEndTime  int32     `json:"autoEndTime" gorm:"default:0;comment:自动结束时间(分钟),0表示不自动结束"`
	Description  string    `json:"description" gorm:"type:text;comment:规则描述"`
	FreeTime     int32     `json:"freeTime" gorm:"default:0;comment:免费时间(分钟),在此时间内可以临时开锁和结束订单不收费"`
	AutoRefund   bool      `json:"autoRefund" gorm:"default:false;comment:是否开启自动退款(预付费模式)"`
	ManualRenew  bool      `json:"manualRenew" gorm:"default:false;comment:是否开启手动续费(预付费模式)"`
	TimeOptions  string    `json:"timeOptions" gorm:"type:text;comment:预付费的时间选项JSON"`
	Sort         int32     `json:"sort" gorm:"default:0;comment:排序"`
	Tag          string    `json:"tag" gorm:"size:255;comment:标签"`
	Duration     int32     `json:"duration" gorm:"default:0;comment:时长"`
	Status       int32     `json:"status" gorm:"default:1;comment:状态 0禁用 1启用"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`
}

// TimeOption 时间选项模型
type TimeOption struct {
	Title        string  `json:"title"`        // 选项标题，如"1小时"、"2小时"
	Duration     int32   `json:"duration"`     // 时长
	DurationUnit string  `json:"durationUnit"` // 时长单位
	Price        float32 `json:"price"`        // 价格
	Discount     float32 `json:"discount"`     // 优惠金额
}

// formatDuration 格式化时长显示
func formatDuration(duration int32, unit string) string {
	if duration <= 0 {
		return "不限制"
	}
	switch unit {
	case "hour":
		return fmt.Sprintf("%d小时", duration)
	case "day":
		return fmt.Sprintf("%d天", duration)
	case "month":
		return fmt.Sprintf("%d月", duration)
	default: // minute
		if duration < 60 {
			return fmt.Sprintf("%d分钟", duration)
		}
		hours := duration / 60
		mins := duration % 60
		if mins == 0 {
			return fmt.Sprintf("%d小时", hours)
		}
		return fmt.Sprintf("%d小时%d分钟", hours, mins)
	}
}
