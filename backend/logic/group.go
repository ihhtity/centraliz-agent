package logic

import (
	"centraliz-backend/dao/mysql"
	"centraliz-backend/model"
	"centraliz-backend/pkg/errno"
	"strconv"

	"gorm.io/gorm"
)

func GetGroupCount() (int64, error) {
	count, err := mysql.GetGroupCount()
	if err != nil {
		return 0, errno.New(errno.InternalError)
	}
	return count, nil
}

func GetGroupListFiltered(merchsID int32, name string, groupType string, page, pageSize int) ([]model.Group, int64, error) {
	groups, total, err := mysql.GetGroupListFiltered(merchsID, name, groupType, page, pageSize)
	if err != nil {
		return nil, 0, errno.New(errno.InternalError)
	}
	return groups, total, nil
}

func GetGroupByID(id string) (*model.Group, error) {
	groupID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errno.New(errno.BadRequest)
	}
	group, err := mysql.GetGroupByID(groupID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.New(errno.GroupNotFound)
		}
		return nil, errno.New(errno.InternalError)
	}
	return group, nil
}

func CreateGroup(group *model.Group) error {
	if err := mysql.CreateGroup(group); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func UpdateGroup(group *model.Group) error {
	if err := mysql.UpdateGroup(group); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func DeleteGroup(id string) error {
	groupID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errno.New(errno.BadRequest)
	}
	if err := mysql.DeleteGroup(groupID); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchDeleteGroup(ids []string) error {
	var groupIDs []uint64
	for _, id := range ids {
		groupID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return errno.New(errno.BadRequest)
		}
		groupIDs = append(groupIDs, groupID)
	}
	if err := mysql.BatchDeleteGroup(groupIDs); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchUpdateGroup(reqs []struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Location    string `json:"location"`
	Phone       string `json:"phone"`
	BindNumber  string `json:"bind_number"`
	ConsumePush string `json:"consume_push"`
}) error {
	if err := mysql.BatchUpdateGroup(reqs); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchUpdateGroupByIDs(ids []string, data map[string]interface{}) error {
	var groupIDs []uint64
	for _, id := range ids {
		groupID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return errno.New(errno.BadRequest)
		}
		groupIDs = append(groupIDs, groupID)
	}
	if err := mysql.BatchUpdateGroupByIDs(groupIDs, data); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}
