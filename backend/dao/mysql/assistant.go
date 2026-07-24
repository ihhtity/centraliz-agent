package mysql

import (
	"centraliz-backend/model"

	"gorm.io/gorm"
)

func CreateAssistantSession(session *model.AssistantSession) error {
	return GetDB().Create(session).Error
}

func GetAssistantSessionsByUserID(userID uint32) ([]model.AssistantSession, error) {
	var sessions []model.AssistantSession
	err := GetDB().Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).Where("user_id = ?", userID).Order("updated_at DESC").Find(&sessions).Error
	return sessions, err
}

func GetAssistantSessionDetail(sessionID uint32) (*model.AssistantSession, error) {
	var session model.AssistantSession
	err := GetDB().Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(&session, sessionID).Error
	return &session, err
}

func UpdateAssistantSession(session *model.AssistantSession) error {
	return GetDB().Save(session).Error
}

func DeleteAssistantSession(sessionID uint32) error {
	tx := GetDB().Begin()

	var session model.AssistantSession
	if err := tx.First(&session, sessionID).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Model(&session).Association("Messages").Clear()

	if err := tx.Delete(&session).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func CreateAssistantMessage(message *model.AssistantMessage) error {
	return GetDB().Create(message).Error
}
