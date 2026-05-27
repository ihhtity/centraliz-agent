package mail

import (
	"fmt"

	"centraliz-backend/pkg/config"
)

// 全局邮件发送器实例
var sender *MailSender

// InitMailService 初始化邮件服务
func InitMailService() error {
	mailConfig := config.AppConfig.Mail
	if mailConfig == nil {
		return fmt.Errorf("邮件配置缺失")
	}

	// 创建邮件发送器
	sender = &MailSender{
		Host:       mailConfig.Host,
		Port:       mailConfig.Port,
		User:       mailConfig.User,
		Password:   mailConfig.Password,
		From:       mailConfig.From,
		VerifyType: mailConfig.VerifyType,
	}

	// log.Logger.Info("邮件服务已初始化", zap.String("host", mailConfig.Host), zap.Int("port", mailConfig.Port), zap.String("from", mailConfig.From))
	return nil
}

// GetMailSender 获取邮件发送器实例
func GetMailSender() *MailSender {
	return sender
}
