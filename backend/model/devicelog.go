package model

import "time"

type Devicelog struct {
	ID         uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	MerchsID   int32     `json:"merchsId" gorm:"column:merchs_id;default:0;comment:商家外键"`
	UsersID    int32     `json:"usersId" gorm:"column:user_id;default:0;comment:用户外键"`
	DevicesID  int32     `json:"devicesId" gorm:"column:devices_id;default:0;comment:设备外键"`
	RoomID     int32     `json:"roomId" gorm:"column:room_id;default:0;comment:房间外键"`
	Code       string    `json:"code" gorm:"size:100;default:'';comment:设备编码"`
	DeviceName string    `json:"deviceName" gorm:"size:64;default:'';comment:使用设备名称"`
	RoomName   string    `json:"roomName" gorm:"size:100;default:'';comment:房间名称"`
	Type       string    `json:"type" gorm:"size:64;default:'手机';comment:使用设备类型：手机 电脑"`
	Control    string    `json:"control" gorm:"size:64;default:'开锁';comment:控制类型：开锁 关锁"`
	Status     string    `json:"status" gorm:"size:64;default:'成功';comment:状态：成功 失败"`
	Occupant   string    `json:"occupant" gorm:"size:64;default:'用户';comment:使用人：用户 商家"`
	Phone      string    `json:"phone" gorm:"size:64;default:'';comment:用户手机号"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;precision:3;comment:创建时间"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;precision:3;comment:更新时间"`
}
