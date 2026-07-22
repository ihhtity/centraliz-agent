package logic

import (
	"centraliz-backend/dao/mysql"
	"centraliz-backend/model"
	"centraliz-backend/pkg/errno"
	"strconv"

	"gorm.io/gorm"
)

func GetDeviceCount() (int64, error) {
	count, err := mysql.GetDeviceCount()
	if err != nil {
		return 0, errno.New(errno.InternalError)
	}
	return count, nil
}

func GetDeviceListFiltered(merchsID int32, groupsID int32, name string, status string, deviceType string, page, pageSize int) ([]model.Device, int64, error) {
	devices, total, err := mysql.GetDeviceListFiltered(merchsID, groupsID, name, status, deviceType, page, pageSize)
	if err != nil {
		return nil, 0, errno.New(errno.InternalError)
	}
	return devices, total, nil
}

func GetDeviceByID(id string) (*model.Device, error) {
	deviceID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errno.New(errno.BadRequest)
	}
	device, err := mysql.GetDeviceByID(deviceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.New(errno.DeviceNotFound)
		}
		return nil, errno.New(errno.InternalError)
	}
	return device, nil
}

func CreateDevice(device *model.Device) error {
	if err := mysql.CreateDevice(device); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func UpdateDevice(device *model.Device) error {
	if err := mysql.UpdateDevice(device); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func DeleteDevice(id string) error {
	deviceID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errno.New(errno.BadRequest)
	}
	if err := mysql.DeleteDevice(deviceID); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchDeleteDevice(ids []string) error {
	var deviceIDs []uint64
	for _, id := range ids {
		deviceID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return errno.New(errno.BadRequest)
		}
		deviceIDs = append(deviceIDs, deviceID)
	}
	if err := mysql.BatchDeleteDevice(deviceIDs); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchUpdateDevice(reqs []struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	LockCount int32  `json:"lock_count"`
}) error {
	if err := mysql.BatchUpdateDevice(reqs); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}