package logic

import (
	"centraliz-backend/dao/mysql"
	"centraliz-backend/pkg/errno"
)

func GetTotalUserCount() (int64, error) {
	count, err := mysql.GetTotalUserCount()
	if err != nil {
		return 0, errno.New(errno.InternalError)
	}
	return count, nil
}