package logic

import (
	"centraliz-backend/dao/mysql"
	"centraliz-backend/model"
	"centraliz-backend/pkg/errno"
	"strconv"

	"gorm.io/gorm"
)

func GetMerchListFiltered(account, phone, role, status string, page, pageSize int) ([]model.Merch, int64, error) {
	merchs, total, err := mysql.GetMerchListFiltered(account, phone, role, status, page, pageSize)
	if err != nil {
		return nil, 0, errno.New(errno.InternalError)
	}
	return merchs, total, nil
}

func GetMerchByID(id string) (*model.Merch, error) {
	merchID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errno.New(errno.BadRequest)
	}
	merch, err := mysql.GetMerchByID(merchID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.New(errno.MerchNotFound)
		}
		return nil, errno.New(errno.InternalError)
	}
	return merch, nil
}

func CreateMerch(merch *model.Merch) error {
	if err := mysql.CreateMerch(merch); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func UpdateMerch(merch *model.Merch) error {
	if err := mysql.UpdateMerch(merch); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func DeleteMerch(id string) error {
	merchID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errno.New(errno.BadRequest)
	}
	if err := mysql.DeleteMerch(merchID); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchDeleteMerch(ids []string) error {
	var merchIDs []uint64
	for _, id := range ids {
		merchID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return errno.New(errno.BadRequest)
		}
		merchIDs = append(merchIDs, merchID)
	}
	if err := mysql.BatchDeleteMerch(merchIDs); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchUpdateMerch(reqs []struct {
	ID       uint64 `json:"id"`
	Account  string `json:"account"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}) error {
	if err := mysql.BatchUpdateMerch(reqs); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}
