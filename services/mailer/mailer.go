package mailer

import (
	"crypto/tls"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/czhj/ahfs/modules/log"
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
	mailQueue chan *Message
)

func NewContext() {
	if setting.MailerService == nil || mailQueue != nil {
		return
	}

	mailQueue = make(chan *Message, setting.MailerService.QueueLength)
	go smtpSender()

	log.Debug("Mailer service is running")
}

func smtpSender() {
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

	var sender gomail.SendCloser
	opened := false

	for {
		select {
		case m, ok := <-mailQueue:
			if !ok {
				return
			}
			if !opened {
				if sender, err = dialer.Dial(); err != nil {
					log.Error("Failed to dial stmp server", zap.Error(err))
				} else {
					opened = true
				}
			}
			message := m.ToMessage()
			log.Debug("New e-mail sending request", zap.Strings("To", message.GetHeader("To")))
			if err := gomail.Send(sender, message); err != nil {
				log.Error("Failed to send emails", zap.Strings("To", message.GetHeader("To")), zap.Error(err))
			} else {
				log.Debug("E-mails sent", zap.Strings("To", message.GetHeader("To")))
			}
		case <-time.After(30 * time.Second):
			if opened {
				if err := sender.Close(); err != nil {
					log.Error("Failed to close stmp connection", zap.Error(err))
				}
				opened = false
			}
		}
	}
}

func SendAsync(msg *Message) {
	go func() {
		mailQueue <- msg
	}()
}
