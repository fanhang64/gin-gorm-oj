package helper

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"time"
)

// Send
// 发送邮件
func Send(from, to, subject, content string) error {
	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(content)
	// 返回EOF错误时候，关闭SSL并重试
	err := e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "fanhangzhou@163.com", "EUFUZDBKBCDPGQIJ", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"},
	)
	return err
}

// SendVerifyCode
// 发送验证码
func SendVerifyCode(to string, verifyCode string) error {
	now := time.Now()
	rand.Seed(now.UnixNano())
	content := fmt.Sprintf(`您的验证码是：<b>%v</b>，有效期5分钟。`, verifyCode)

	return Send("fanhangzhou@163.com", to, "验证码测试", content)
}
