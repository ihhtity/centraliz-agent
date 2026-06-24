package model

import "time"

type Order struct {
	ID         uint64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	MerchsID   int32      `json:"merchsId" gorm:"column:merchs_id;default:0;comment:商家外键"`
	UsersID    int32      `json:"usersId" gorm:"column:users_id;default:0;comment:用户外键"`
	RoomsID    int32      `json:"roomsId" gorm:"column:rooms_id;default:0;comment:房间外键"`
	GroupsID   int32      `json:"groupsId" gorm:"column:groups_id;default:0;comment:分组外键"`
	Name       string     `json:"name" gorm:"size:255;comment:订单名称，可为空"`
	Code       string     `json:"code" gorm:"size:255;comment:订单编号，可为空"`
	Status     string     `json:"status" gorm:"size:20;default:'进行中';comment:订单状态：未完成 进行中 已完成 申请退款 已退款 拒绝退款"`
	Tag        string     `json:"tag" gorm:"size:20;default:'已完成';comment:订单标签：已完成 申请退款 已退款 拒绝退款"`
	Amount     int32      `json:"amount" gorm:"default:1;comment:商品数量"`
	Duration   int32      `json:"duration" gorm:"default:0;comment:使用时长"`
	Price      float64    `json:"price" gorm:"type:decimal(10,2);comment:支付金额"`
	Deposit    float64    `json:"deposit" gorm:"type:decimal(10,2);comment:押金"`
	UserPhone  string     `json:"userPhone" gorm:"size:20;comment:用户手机号"`
	MerchPhone string     `json:"merchPhone" gorm:"size:20;comment:商家手机号"`
	ReqDate    string     `json:"reqDate" gorm:"size:255;comment:支付时间日期"`
	FreeTime   time.Time  `json:"freeTime" gorm:"column:free_time;precision:3;comment:免费时间"`
	StartTime  time.Time  `json:"startTime" gorm:"column:start_time;precision:3;comment:订单开始时间"`
	EndTime    *time.Time `json:"endTime" gorm:"column:end_time;precision:3;comment:订单结束时间"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`
}
