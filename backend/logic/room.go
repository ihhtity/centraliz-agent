package logic

import (
	"centraliz-backend/dao/mysql"
	"centraliz-backend/model"
	"centraliz-backend/pkg/errno"
	"strconv"

	"gorm.io/gorm"
)

func GetRoomCount() (int64, error) {
	count, err := mysql.GetRoomCount()
	if err != nil {
		return 0, errno.New(errno.InternalError)
	}
	return count, nil
}

func GetRoomListFiltered(merchsID int32, groupsID int32, name string, status string, boardNo string, lockNo string, tag string, page, pageSize int) ([]model.Room, int64, error) {
	rooms, total, err := mysql.GetRoomListFiltered(merchsID, groupsID, name, status, boardNo, lockNo, tag, page, pageSize)
	if err != nil {
		return nil, 0, errno.New(errno.InternalError)
	}
	return rooms, total, nil
}

func GetRoomByID(id string) (*model.Room, error) {
	roomID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errno.New(errno.BadRequest)
	}
	room, err := mysql.GetRoomByID(roomID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.New(errno.RoomNotFound)
		}
		return nil, errno.New(errno.InternalError)
	}
	return room, nil
}

func CreateRoom(room *model.Room) error {
	if err := mysql.CreateRoom(room); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func UpdateRoom(room *model.Room) error {
	if err := mysql.UpdateRoom(room); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func DeleteRoom(id string) error {
	roomID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errno.New(errno.BadRequest)
	}
	if err := mysql.DeleteRoom(roomID); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchDeleteRoom(ids []string) error {
	var roomIDs []uint64
	for _, id := range ids {
		roomID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return errno.New(errno.BadRequest)
		}
		roomIDs = append(roomIDs, roomID)
	}
	if err := mysql.BatchDeleteRoom(roomIDs); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchUpdateRoom(reqs []struct {
	ID     uint64  `json:"id"`
	Name   string  `json:"name"`
	Tag    string  `json:"tag"`
	Status string  `json:"status"`
	LockNo string  `json:"lock_no"`
	Price  float32 `json:"price"`
	Image  string  `json:"image"`
}) error {
	if err := mysql.BatchUpdateRoom(reqs); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchUpdateRoomByIDs(ids []string, data map[string]interface{}) error {
	var roomIDs []uint64
	for _, id := range ids {
		roomID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return errno.New(errno.BadRequest)
		}
		roomIDs = append(roomIDs, roomID)
	}
	if err := mysql.BatchUpdateRoomByIDs(roomIDs, data); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}
