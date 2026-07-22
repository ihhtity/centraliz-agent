package logic

import (
	"centraliz-backend/dao/mysql"
	"centraliz-backend/model"
	"centraliz-backend/pkg/errno"
	"strconv"

	"gorm.io/gorm"
)

func GetDeviceLogListFiltered(merchsID, devicesID, roomID int32, code, deviceType, control, status string, page, pageSize int) ([]model.Devicelog, int64, error) {
	logs, total, err := mysql.GetDeviceLogListFiltered(merchsID, devicesID, roomID, code, deviceType, control, status, page, pageSize)
	if err != nil {
		return nil, 0, errno.New(errno.InternalError)
	}
	return logs, total, nil
}

func GetDeviceLogByID(id string) (*model.Devicelog, error) {
	logID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errno.New(errno.BadRequest)
	}
	log, err := mysql.GetDeviceLogByID(logID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.New(errno.DeviceLogNotFound)
		}
		return nil, errno.New(errno.InternalError)
	}
	return log, nil
}

func DeleteDeviceLog(id string) error {
	logID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errno.New(errno.BadRequest)
	}
	if err := mysql.DeleteDeviceLog(logID); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchDeleteDeviceLog(ids []string) error {
	var logIDs []uint64
	for _, id := range ids {
		logID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return errno.New(errno.BadRequest)
		}
		logIDs = append(logIDs, logID)
	}
	if err := mysql.BatchDeleteDeviceLog(logIDs); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}
