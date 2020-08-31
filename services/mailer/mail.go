package mailer

import (
	"bytes"
	"html/template"

	"github.com/czhj/ahfs/models"
	"github.com/czhj/ahfs/modules/log"
	"github.com/czhj/ahfs/modules/setting"
	"go.uber.org/zap"
)

const (
	mailAuthActivateEmail = "auth/activate_email"
	mailAuthResetPassword = "auth/reset_password"
)

var (
	bodyTemplates *template.Template
)

func InitMailTemplate(bodyTpl *template.Template) {
	bodyTemplates = bodyTpl
}

func SendUserMail(u *models.User, tpl string, code, subject string) {
	data := map[string]interface{}{
		"DisplayName":       u.Nickname,
		"ActiveCodeLives":   setting.Service.ActiveCodeLive,
		"ResetPwdCodeLives": setting.Service.ResetPasswordCodeLive,
		"Code":              code,
	}

	var content bytes.Buffer

	if err := bodyTemplates.ExecuteTemplate(&content, tpl, data); err != nil {
		log.Error("Failed to render mail template", zap.String("name", tpl), zap.Error(err))
		return
	}

	msg := NewMessage([]string{u.Email}, subject, content.String())
	SendAsync(msg)
}

func SendActiveCodeMail(email string, code string) {
	SendUserMail(&models.User{
		Email: email,
	}, mailAuthActivateEmail, code, "邮箱验证码(email verification code)")
}

func SendResetPwdCodeMail(email string, code string) {
	SendUserMail(&models.User{
		Email: email,
	}, mailAuthResetPassword, code, "修改密码验证码(reset password  verification code)")
}
