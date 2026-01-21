package schema

import "strings"

type Captcha struct {
	CaptchaID string
}

type LoginForm struct {
	Username    string
	Password    string
	CaptchaID   string
	CaptchaCode string
}

func (a *LoginForm) Trim() *LoginForm {
	a.Username = strings.TrimSpace(a.Username)
	a.CaptchaCode = strings.TrimSpace(a.CaptchaCode)
	return a
}

type UpdateLoginPassword struct {
	OldPassword string
	NewPassword string
}

type LoginToken struct {
	AccessToken string
	TokenType   string
	ExpiresAt   int64
}

type UpdateCurrentUser struct {
	Name   string
	Phone  string
	Email  string
	Remark string
}
