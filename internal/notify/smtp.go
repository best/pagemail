package notify

import (
	"crypto/tls"
	"fmt"
	"io"

	"gopkg.in/gomail.v2"
)

type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromName  string
	FromEmail string
	UseTLS    bool
}

type SMTPSender struct {
	config *SMTPConfig
}

func NewSMTPSender(cfg *SMTPConfig) *SMTPSender {
	return &SMTPSender{config: cfg}
}

type EmailMessage struct {
	To          []string
	Subject     string
	Body        string
	HTMLBody    string
	Attachments []Attachment
}

type Attachment struct {
	Filename    string
	ContentType string
	Reader      io.Reader
}

func (s *SMTPSender) Send(msg *EmailMessage) error {
	m := gomail.NewMessage()

	from := s.config.FromEmail
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	}
	m.SetHeader("From", from)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)

	if msg.HTMLBody != "" {
		m.SetBody("text/html", msg.HTMLBody)
		if msg.Body != "" {
			m.AddAlternative("text/plain", msg.Body)
		}
	} else {
		m.SetBody("text/plain", msg.Body)
	}

	for _, att := range msg.Attachments {
		m.Attach(att.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, att.Reader)
			return err
		}))
	}

	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)
	if s.config.UseTLS {
		d.TLSConfig = &tls.Config{ServerName: s.config.Host}
	}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *SMTPSender) TestConnection() error {
	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)
	if s.config.UseTLS {
		d.TLSConfig = &tls.Config{ServerName: s.config.Host}
	}

	closer, err := d.Dial()
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer closer.Close()

	return nil
}
