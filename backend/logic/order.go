package logic

import (
	"centraliz-backend/dao/mysql"
	"centraliz-backend/model"
	"centraliz-backend/pkg/errno"
	"strconv"

	"gorm.io/gorm"
)

func GetOrderCount() (int64, error) {
	count, err := mysql.GetOrderCount()
	if err != nil {
		return 0, errno.New(errno.InternalError)
	}
	return count, nil
}

func GetOrderListFiltered(merchsID int32, usersID int32, roomsID int32, status string, orderCode string, orderNo string, userPhone string, payType string, page, pageSize int) ([]model.Order, int64, error) {
	orders, total, err := mysql.GetOrderListFiltered(merchsID, usersID, roomsID, status, orderCode, orderNo, userPhone, payType, page, pageSize)
	if err != nil {
		return nil, 0, errno.New(errno.InternalError)
	}
	return orders, total, nil
}

func GetOrderByID(id string) (*model.Order, error) {
	orderID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errno.New(errno.BadRequest)
	}
	order, err := mysql.GetOrderByID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.New(errno.OrderNotFound)
		}
		return nil, errno.New(errno.InternalError)
	}
	return order, nil
}

func CreateOrder(order *model.Order) error {
	if err := mysql.CreateOrder(order); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func UpdateOrder(order *model.Order) error {
	if err := mysql.UpdateOrder(order); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func DeleteOrder(id string) error {
	orderID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errno.New(errno.BadRequest)
	}
	if err := mysql.DeleteOrder(orderID); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchDeleteOrder(ids []string) error {
	var orderIDs []uint64
	for _, id := range ids {
		orderID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return errno.New(errno.BadRequest)
		}
		orderIDs = append(orderIDs, orderID)
	}
	if err := mysql.BatchDeleteOrder(orderIDs); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchUpdateOrder(reqs []struct {
	ID     uint64 `json:"id"`
	Status string `json:"status"`
	Remark string `json:"remark"`
}) error {
	if err := mysql.BatchUpdateOrder(reqs); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchUpdateOrderByIDs(ids []string, data map[string]interface{}) error {
	var orderIDs []uint64
	for _, id := range ids {
		orderID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return errno.New(errno.BadRequest)
		}
		orderIDs = append(orderIDs, orderID)
	}
	if err := mysql.BatchUpdateOrderByIDs(orderIDs, data); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func GetTodayOrderCount() (int64, error) {
	count, err := mysql.GetTodayOrderCount()
	if err != nil {
		return 0, errno.New(errno.InternalError)
	}
	return count, nil
}

func GetTodayRevenue() (float64, error) {
	revenue, err := mysql.GetTodayRevenue()
	if err != nil {
		return 0, errno.New(errno.InternalError)
	}
	return revenue, nil
}

func GetTotalRevenue() (float64, error) {
	revenue, err := mysql.GetTotalRevenue()
	if err != nil {
		return 0, errno.New(errno.InternalError)
	}
	return revenue, nil
}

func GetOrderTrend(days int) ([]map[string]interface{}, error) {
	data, err := mysql.GetOrderTrend(days)
	if err != nil {
		return nil, errno.New(errno.InternalError)
	}
	return data, nil
}

func GetRevenueTrend(days int) ([]map[string]interface{}, error) {
	data, err := mysql.GetRevenueTrend(days)
	if err != nil {
		return nil, errno.New(errno.InternalError)
	}
	return data, nil
}
