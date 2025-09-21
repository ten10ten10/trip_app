package email

import (
	"context"
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Sender interface {
	SendVerificationEmail(ctx context.Context, recipientEmail, rawToken, rawPassword string) error
}

type emailSender struct {
	smtpHost     string
	smtpPort     int
	smtpUser     string
	smtpPassword string
	fromEmail    string
}

func NewSender() Sender {
	return &emailSender{}
}

func NewEmailSender(host string, portStr, user, appPassword, from string) (Sender, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid smtp port: %w", err)
	}
	return &emailSender{
		smtpHost:     host,
		smtpPort:     port,
		smtpUser:     user,
		smtpPassword: appPassword,
		fromEmail:    from,
	}, nil
}

func (e *emailSender) SendVerificationEmail(ctx context.Context, recipientEmail, rawToken, rawPassword string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.fromEmail)
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", "【Trip App】アカウント有効化のご案内")

	body := fmt.Sprintf(`
	<p>Trip Appへのご登録ありがとうございます。</p>
	<p>ご登録のメールアドレスをご確認いただき、お間違いなければ、下のリンクをクリックしてメールアドレスの認証を完了してください。</p>
	<hr>
	<p><b>認証トークン:</b> %s</p>
	<hr>
	<p>メールアドレスの認証完了後、以下の初回パスワードを使用してログインしてください。</p>
	<hr>
	<p><b>初回パスワード:</b> %s</p>
	<hr>
	<p>※このパスワードは初回ログイン後に変更してください。</p>
	<p>※認証トークンの有効期限は30分です。有効期限を過ぎた場合は、再度サインアップをお願いいたします。</p>
	<p>このメールにお心当たりがない場合は、お手数ですが本メールを破棄してください。</p>
	`, rawToken, rawPassword)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(e.smtpHost, e.smtpPort, e.smtpUser, e.smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("✅ Verification email sent to %s\n", recipientEmail)
	return nil
}
