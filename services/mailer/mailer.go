package mailer

import (
	"context"
	"crypto/tls"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/czhj/ahfs/modules/log"
	"github.com/czhj/ahfs/modules/queue"
	"github.com/czhj/ahfs/modules/setting"
	"github.com/jaytaylor/html2text"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type Message struct {
	Subject     string
	FromAddress string
	FromName    string
	To          []string
	Date        time.Time
	Body        string
	Headers     map[string][]string
}

func (m *Message) ToMessage() *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", m.FromAddress, m.FromName)
	msg.SetHeader("To", m.To...)

	for header := range m.Headers {
		msg.SetHeader(header, m.Headers[header]...)
	}

	if len(setting.MailerService.SubjectPrefix) > 0 {
		msg.SetHeader("Subject", setting.MailerService.SubjectPrefix+" "+m.Subject)
	} else {
		msg.SetHeader("Subject", m.Subject)
	}
	msg.SetDateHeader("Date", m.Date)

	msg.SetHeader("X-Auto-Response-Suppress", "All")

	plainBody, err := html2text.FromString(m.Body)
	if err != nil || setting.MailerService.SendAsPlainText {
		if strings.Contains(m.Body, "<html>") {
			log.Warn("Mail contains HTML but configurad to send as a plain text.")
		}
		msg.SetBody("text/plain", plainBody)
	} else {
		msg.SetBody("text/plain", plainBody)
		msg.AddAlternative("text/html", m.Body)
	}
	return msg
}

func (m *Message) SetHeader(field string, value ...string) {
	m.Headers[field] = value
}

func NewMessageFrom(to []string, fromName, fromAddress, subject, body string) *Message {
	return &Message{
		Subject:     subject,
		Body:        body,
		To:          to,
		FromAddress: fromAddress,
		FromName:    fromName,
		Date:        time.Now(),
		Headers:     map[string][]string{},
	}
}

func NewMessage(to []string, subject, body string) *Message {
	return NewMessageFrom(to, setting.MailerService.FromName, setting.MailerService.FromEmail, subject, body)
}

var (
	mailerQueue queue.Queue
	Sender      *gomail.Dialer
)

func NewContext() {
	if setting.MailerService == nil || mailerQueue != nil {
		return
	}

	Sender = smtpSender()
	if Sender == nil {
		return
	}

	mailerQueue = queue.CreateQueue("mail", func(data ...queue.Data) {
		for _, dat := range data {
			msg := dat.(*Message)
			message := msg.ToMessage()

			if err := Sender.DialAndSend(message); err != nil {
				log.Error("Failed to send emails", zap.Strings("To", message.GetHeader("To")), zap.Error(err))
			} else {
				log.Debug("E-mails sent", zap.Strings("To", message.GetHeader("To")))
			}
		}
	}, &Message{})

	mailerQueue.Run(func(c context.Context, f func()) {
		f()
	}, func(c context.Context, f func()) {
		f()
	})
	log.Debug("Mailer service is running")
}
func smtpSender() *gomail.Dialer {
	cfg := setting.MailerService

	host, port, err := net.SplitHostPort(cfg.Host)
	if err != nil {
		log.Fatal("Failed to split smtp host", zap.String("host", cfg.Host), zap.Error(err))
	}

	iport, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal("Cannot parse port to int", zap.String("port", port), zap.Error(err))
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.SkipVerify,
		ServerName:         host,
	}

	if cfg.UseCertificate {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			log.Fatal("Failed to load x509 key pair for mailer", zap.Error(err))
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	dialer := &gomail.Dialer{
		Host:     host,
		Port:     iport,
		Username: cfg.Username,
		Password: cfg.Password,
	}

	if cfg.IsTLSEnable || (strings.HasSuffix(port, "465")) {
		dialer.SSL = true
		dialer.TLSConfig = tlsConfig
	}

	return dialer
}

func SendAsync(msg *Message) {
	go func() {
		_ = mailerQueue.Push(msg)
	}()
}
