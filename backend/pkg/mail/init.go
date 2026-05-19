package mail

import (
	"fmt"
	"net/smtp"

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

	// 根据SMTP验证方式选择TLS或SSL
	var auth smtp.Auth
	if mailConfig.VerifyType == "SSL" {
		auth = smtp.PlainAuth("", mailConfig.User, mailConfig.Password, mailConfig.Host)
	} else if mailConfig.VerifyType == "TLS" {
		auth = smtp.PlainAuth("", mailConfig.User, mailConfig.Password, mailConfig.Host)
	} else {
		return fmt.Errorf("不支持的SMTP验证类型: %s", mailConfig.VerifyType)
	}

	// 创建邮件发送器
	sender = &MailSender{
		Host:     mailConfig.Host,
		Port:     mailConfig.Port,
		User:     mailConfig.User,
		Password: mailConfig.Password,
		From:     mailConfig.From,
		Auth:     auth,
	}

	// log.Logger.Info("邮件服务已初始化", zap.String("host", mailConfig.Host), zap.Int("port", mailConfig.Port), zap.String("from", mailConfig.From))
	return nil
}

// GetMailSender 获取邮件发送器实例
func GetMailSender() *MailSender {
	return sender
}
