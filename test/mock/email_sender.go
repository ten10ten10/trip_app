package mock

import (
	"context"
	"trip_app/internal/infrastructure/email"
)

// MockEmailSender はテスト用のメール送信モック
// 実際にはメールを送信せず、トークンとパスワードを保存するだけ
type MockEmailSender struct {
	lastToken    string
	lastPassword string
}

// NewMockEmailSender はMockEmailSenderの新しいインスタンスを作成
func NewMockEmailSender() *MockEmailSender {
	return &MockEmailSender{}
}

// SendVerificationEmail はメール送信をシミュレートし、トークンとパスワードを保存
func (m *MockEmailSender) SendVerificationEmail(ctx context.Context, recipientEmail, rawToken, rawPassword string) error {
	m.lastToken = rawToken
	m.lastPassword = rawPassword
	return nil
}

// GetLastToken は最後に送信されたトークンを返す
func (m *MockEmailSender) GetLastToken() string {
	return m.lastToken
}

// GetLastPassword は最後に送信されたパスワードを返す
func (m *MockEmailSender) GetLastPassword() string {
	return m.lastPassword
}

// コンパイル時にinterfaceを実装していることを確認
var _ email.Sender = (*MockEmailSender)(nil)
