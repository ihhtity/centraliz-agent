package errno

import "fmt"

type Errno struct {
	Code    int
	Message string
}

func (e Errno) Error() string {
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

func New(e Errno) error {
	return e
}

var (
	Success       = Errno{Code: 0, Message: "success"}
	InternalError = Errno{Code: 500, Message: "内部服务器错误"}
	BadRequest    = Errno{Code: 400, Message: "请求参数错误"}
	Unauthorized  = Errno{Code: 401, Message: "未授权"}
	Forbidden     = Errno{Code: 403, Message: "禁止访问"}
	NotFound      = Errno{Code: 404, Message: "资源不存在"}

	RoomNotFound      = Errno{Code: 40001, Message: "房间不存在"}
	DeviceNotFound    = Errno{Code: 40002, Message: "设备不存在"}
	GroupNotFound     = Errno{Code: 40003, Message: "分组不存在"}
	RuleNotFound      = Errno{Code: 40004, Message: "规则不存在"}
	OrderNotFound     = Errno{Code: 40005, Message: "订单不存在"}
	MerchNotFound     = Errno{Code: 40006, Message: "商家不存在"}
	UserNotFound      = Errno{Code: 40007, Message: "用户不存在"}
	DeviceLogNotFound = Errno{Code: 40008, Message: "设备日志不存在"}

	RoomNameExists   = Errno{Code: 40011, Message: "房间名称已存在"}
	DeviceCodeExists = Errno{Code: 40012, Message: "设备编码已存在"}
	GroupNameExists  = Errno{Code: 40013, Message: "分组名称已存在"}
	RuleNameExists   = Errno{Code: 40014, Message: "规则名称已存在"}

	LockNoExists    = Errno{Code: 40021, Message: "锁号已被使用"}
	LockCountExceed = Errno{Code: 40022, Message: "锁数量不足"}
)
