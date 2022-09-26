package test

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "fanhangzhou <fanhangzhou@163.com>"
	e.To = []string{"fanhangzhou@163.com"}
	e.Subject = "验证码测试"
	e.HTML = []byte("你的验证码为：<b>12356</b>")
	//err := e.Send("smtp.163.com:465",
	//	smtp.PlainAuth("fan6512519", "test@gmail.com", "", "smtp.163.com"))
	// 返回EOF错误时候，关闭SSL并重试
	err := e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "fanhangzhou@163.com", "EUFUZDBKBCDPGQIJ", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"},
	)
	if err != nil {
		t.Fatal(err)
	}
}
