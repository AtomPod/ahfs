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
		"enabled":            false,
		"queue_length":       512,
		"subject_prefix":     "",
		"send_as_plain_text": false,
		"from":               "",
		"host":               "",
		"username":           "",
		"password":           "",
		"skip_verify":        true,
		"use_certificate":    false,
		"cert_file":          "",
		"key_file":           "",
		"is_tls_enable":      false,
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
