package logic

import (
	"centraliz-backend/dao/mysql"
	"centraliz-backend/model"
	"centraliz-backend/pkg/errno"
	"strconv"

	"gorm.io/gorm"
)

func GetRuleListFiltered(merchsID int32, name string, ruleType string, page, pageSize int) ([]model.Rule, int64, error) {
	rules, total, err := mysql.GetRuleListFiltered(merchsID, name, ruleType, page, pageSize)
	if err != nil {
		return nil, 0, errno.New(errno.InternalError)
	}
	return rules, total, nil
}

func GetRuleByID(id string) (*model.Rule, error) {
	ruleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, errno.New(errno.BadRequest)
	}
	rule, err := mysql.GetRuleByID(uint32(ruleID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.New(errno.RuleNotFound)
		}
		return nil, errno.New(errno.InternalError)
	}
	return rule, nil
}

func CreateRule(rule *model.Rule) error {
	if err := mysql.CreateRule(rule); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func UpdateRule(rule *model.Rule) error {
	if err := mysql.UpdateRule(rule); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func DeleteRule(id string) error {
	ruleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return errno.New(errno.BadRequest)
	}
	if err := mysql.DeleteRule(uint32(ruleID)); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchDeleteRule(ids []string) error {
	var ruleIDs []uint32
	for _, id := range ids {
		ruleID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return errno.New(errno.BadRequest)
		}
		ruleIDs = append(ruleIDs, uint32(ruleID))
	}
	if err := mysql.BatchDeleteRule(ruleIDs); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchUpdateRule(reqs []struct {
	ID          uint32  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Mode        string  `json:"mode"`
	Price       float32 `json:"price"`
	Deposit     float32 `json:"deposit"`
	Description string  `json:"description"`
}) error {
	if err := mysql.BatchUpdateRule(reqs); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}

func BatchUpdateRuleByIDs(ids []string, data map[string]interface{}) error {
	var ruleIDs []uint32
	for _, id := range ids {
		ruleID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return errno.New(errno.BadRequest)
		}
		ruleIDs = append(ruleIDs, uint32(ruleID))
	}
	if err := mysql.BatchUpdateRuleByIDs(ruleIDs, data); err != nil {
		return errno.New(errno.InternalError)
	}
	return nil
}