package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

// MailSender 邮件发送器
type MailSender struct {
	Host       string
	Port       int
	User       string
	Password   string
	From       string
	VerifyType string // SSL or TLS
}

// Send 发送邮件
func (m *MailSender) Send(to []string, subject, body string) error {
	// 构建邮件内容，使用"BSLD"作为发件人显示名称
	from := m.From
	if from == "" {
		from = m.User
	}
	
	// 设置发件人显示名称为"BSLD"
	fromWithDisplayName := fmt.Sprintf("BSLD <%s>", from)
	
	headers := make(map[string]string)
	headers["From"] = fromWithDisplayName
	headers["To"] = to[0]
	headers["Subject"] = "验证码" // 修改邮件标题
	headers["Content-Type"] = "text/html; charset=UTF-8" // 使用HTML格式

	var msg string
	for k, v := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n"

	// 修改邮件正文内容
	body = `
    <html>
    <body>
        <p>验证码</p>
        <p>您的验证码是：</p>
        <div style="font-size: 24px; font-weight: bold; color: #1E90FF; text-align: center; padding: 20px; background-color: #F5F5F5;">
            ` + body + `
        </div>
        <p>验证码5分钟内有效，请勿泄露给他人。</p>
        <p>如果这不是您本人的操作，请忽略此邮件。</p>
    </body>
    </html>
    `

	msg += body

	var conn net.Conn
	var err error

	// 根据验证类型选择连接方式
	if m.VerifyType == "SSL" {
		// SSL连接（用于465端口）
		conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", m.Host, m.Port), &tls.Config{
			ServerName: m.Host,
		})
		if err != nil {
			return fmt.Errorf("SSL连接失败: %v", err)
		}
	} else {
		// 普通TCP连接（用于25或587端口，然后升级到TLS）
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", m.Host, m.Port))
		if err != nil {
			return fmt.Errorf("TCP连接失败: %v", err)
		}
		
		// 如果是TLS，尝试升级连接
		if m.VerifyType == "TLS" {
			host, _, _ := net.SplitHostPort(conn.LocalAddr().String())
			tlsConn := tls.Client(conn, &tls.Config{ServerName: host})
			if err := tlsConn.Handshake(); err != nil {
				conn.Close()
				return fmt.Errorf("TLS握手失败: %v", err)
			}
			conn = tlsConn
		}
	}

	defer conn.Close()

	client, err := smtp.NewClient(conn, m.Host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}

	// 如果不是SSL连接，尝试STARTTLS
	if m.VerifyType != "SSL" {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err = client.StartTLS(&tls.Config{ServerName: m.Host}); err != nil {
				client.Close()
				return fmt.Errorf("STARTTLS失败: %v", err)
			}
		}
	}

	// 认证
	if m.User != "" && m.Password != "" {
		auth := smtp.PlainAuth("", m.User, m.Password, m.Host)
		if err = client.Auth(auth); err != nil {
			client.Close()
			return fmt.Errorf("SMTP认证失败: %v", err)
		}
	}

	// 发送邮件
	if err = client.Mail(m.User); err != nil {
		client.Close()
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			client.Close()
			return fmt.Errorf("设置收件人失败: %v", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		client.Close()
		return fmt.Errorf("获取数据写入器失败: %v", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		w.Close()
		client.Close()
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = w.Close()
	if err != nil {
		client.Close()
		return fmt.Errorf("关闭数据写入器失败: %v", err)
	}

	err = client.Quit()
	if err != nil {
		return fmt.Errorf("关闭SMTP客户端失败: %v", err)
	}

	return nil
}