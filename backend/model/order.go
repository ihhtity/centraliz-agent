package model

import "time"

type Order struct {
	ID          uint64     `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	OrderNo     string     `json:"orderNo" gorm:"size:64;comment:订单编号"`
	MerchsID    int32      `json:"merchsId" gorm:"column:merchs_id;default:0;comment:商家外键"`
	UsersID     int32      `json:"usersId" gorm:"column:users_id;default:0;comment:用户外键"`
	RoomsID     int32      `json:"roomsId" gorm:"column:rooms_id;default:0;comment:房间外键"`
	GroupsID    int32      `json:"groupsId" gorm:"column:groups_id;default:0;comment:分组外键"`
	RulesID     int32      `json:"rulesId" gorm:"column:rules_id;default:0;comment:规则外键"`
	Name        string     `json:"name" gorm:"size:255;comment:订单名称，可为空"`
	Code        string     `json:"code" gorm:"size:255;comment:订单编号，可为空"`
	PayCode     string     `json:"payCode" gorm:"size:255;comment:支付单号"`
	Type        string     `json:"type" gorm:"size:20;comment:订单类型"`
	Mode        string     `json:"mode" gorm:"size:20;comment:模式"`
	Status      string     `json:"status" gorm:"size:20;default:'进行中';comment:订单状态：未完成 进行中 已完成 申请退款 已退款 拒绝退款"`
	Tag         string     `json:"tag" gorm:"size:20;default:'未完成';comment:订单标签：未完成 已完成 押金"`
	Amount      int32      `json:"amount" gorm:"default:1;comment:商品数量"`
	Duration    int32      `json:"duration" gorm:"default:0;comment:使用时长"`
	Price       float64    `json:"price" gorm:"type:decimal(10,2);comment:支付金额"`
	Deposit     float64    `json:"deposit" gorm:"type:decimal(10,2);comment:押金"`
	PayPrice    float64    `json:"payPrice" gorm:"type:decimal(10,2);comment:实际支付金额"`
	PayType     string     `json:"payType" gorm:"size:20;comment:支付方式"`
	RefundPrice float64    `json:"refundPrice" gorm:"type:decimal(10,2);comment:退款金额"`
	Remark      string     `json:"remark" gorm:"type:text;comment:订单备注"`
	UserPhone   string     `json:"userPhone" gorm:"size:20;comment:用户手机号"`
	MerchPhone  string     `json:"merchPhone" gorm:"size:20;comment:商家手机号"`
	ReqSeqID    string     `json:"reqSeqId" gorm:"size:64;comment:汇付请求流水号"`
	ReqDate     string     `json:"reqDate" gorm:"size:20;comment:汇付支付日期"`
	FreeTime    time.Time  `json:"freeTime" gorm:"column:free_time;precision:3;comment:免费时间"`
	StartTime   time.Time  `json:"startTime" gorm:"column:start_time;precision:3;comment:订单开始时间"`
	EndTime     *time.Time `json:"endTime" gorm:"column:end_time;precision:3;comment:订单结束时间"`
	PayTime     *time.Time `json:"payTime" gorm:"column:pay_time;precision:3;comment:支付时间"`
	RefundTime  *time.Time `json:"refundTime" gorm:"column:refund_time;precision:3;comment:退款时间"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`
}
