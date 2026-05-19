package mail

import (
	"fmt"
	"net/smtp"
)

// MailSender 邮件发送器
type MailSender struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	Auth     smtp.Auth
}

// Send 发送邮件
func (m *MailSender) Send(to []string, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)
	err := smtp.SendMail(addr, m.Auth, m.From, to, []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)))
	if err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}
	return nil
}
