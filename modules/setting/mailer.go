package setting

import (
	"net/mail"

	"github.com/czhj/ahfs/modules/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Mailer struct {
	QueueLength     int
	SubjectPrefix   string
	SendAsPlainText bool
	From            string
	FromName        string
	FromEmail       string

	//SMTP
	Host           string
	Username       string
	Password       string
	SkipVerify     bool
	UseCertificate bool
	CertFile       string
	KeyFile        string
	IsTLSEnable    bool
}

var (
	MailerService *Mailer
)

func newMailService() {
	viper.SetDefault("mailer", map[string]interface{}{
		"enabled":         false,
		"queueLength":     512,
		"subjectPrefix":   "",
		"sendAsPlainText": false,
		"from":            "",
		"host":            "",
		"username":        "",
		"password":        "",
		"skipVerify":      true,
		"useCertificate":  false,
		"certFile":        "",
		"keyFile":         "",
		"isTLSEnable":     false,
	})

	mailerCfg := viper.Sub("mailer")

	if !mailerCfg.GetBool("enabled") {
		return
	}

	MailerService = &Mailer{}
	if err := mailerCfg.Unmarshal(MailerService); err != nil {
		log.Fatal("Failed to unmarshal mailer config", zap.Error(err))
	}

	if len(MailerService.From) == 0 {
		MailerService.From = MailerService.Username
	}

	parsed, err := mail.ParseAddress(MailerService.From)
	if err != nil {
		log.Fatal("Invalid mailer.FROM", zap.String("from", MailerService.From), zap.Error(err))
	}

	MailerService.FromName = parsed.Name
	MailerService.FromEmail = parsed.Address

	log.Info("Mail Service Enabled")
}
