package model

import (
	"fmt"
	"time"
)

// Rule 规则模型
type Rule struct {
	ID           uint32     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	MerchsID     int32      `json:"merchsId" gorm:"column:merchs_id;not null;comment:商家外键"`
	Name         string     `json:"name" gorm:"size:100;not null;comment:规则名称"`
	Type         string     `json:"type" gorm:"size:20;not null;comment:规则类型 free:免费模式 charge:收费模式"`
	Mode         string     `json:"mode" gorm:"size:50;default:single;comment:模式类型 single:单次开锁;deposit:一存一取;pay_single:单次付费;pay_deposit:先存后取;pay_hourly:按时付费;pay_time:预付费"`
	Price        *float64   `json:"price" gorm:"type:decimal(10,2);default:0.00;comment:单价（元）"`
	Deposit      *float64   `json:"deposit" gorm:"type:decimal(10,2);default:0.00;comment:押金（元）"`
	DurationUnit string     `json:"durationUnit" gorm:"size:20;default:hour;comment:时长单位 hour:小时 day:天 month:月 minute:分钟"`
	AutoEndTime  *int64     `json:"autoEndTime" gorm:"default:0;comment:自动结束时间(分钟),0表示不自动结束"`
	FreeTime     *int64     `json:"freeTime" gorm:"default:0;comment:免费时间(分钟),在此时间内可以临时开锁和结束订单不收费"`
	AutoRefund   *bool      `json:"autoRefund" gorm:"default:false;comment:是否开启自动退款(预付费模式)"`
	ManualRenew  *bool      `json:"manualRenew" gorm:"default:false;comment:是否开启手动续费(预付费模式)"`
	Description  *string    `json:"description" gorm:"size:500;comment:规则描述"`
	TimeOptions  *string    `json:"timeOptions" gorm:"type:text;comment:预付费的时间选项JSON"`
	CreatedAt    *time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`
	UpdatedAt    *time.Time `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`
}

// TimeOption 时间选项
type TimeOption struct {
	Title        string  `json:"title"`        // 选项标题，如"1小时"、"2小时"
	Duration     int64   `json:"duration"`     // 时长
	DurationUnit string  `json:"durationUnit"` // 时长单位
	Price        float64 `json:"price"`        // 价格
	Discount     float64 `json:"discount"`     // 优惠金额
}

// formatDuration 格式化时长显示
func formatDuration(duration int64, unit string) string {
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
