package mail

import (
	"context"
	"fmt"

	"github.com/GoTaskFlow/internal/notifications/mail/config"
	emailModel "github.com/GoTaskFlow/internal/notifications/mail/model"
	log "github.com/GoTaskFlow/pkg/logger/model"
	"gopkg.in/gomail.v2"
)

const (
	fromHeader      = "From"
	toHeader        = "To"
	subjectHeader   = "Subject"
	htmlContentType = "text/html"
)

type Service struct {
	dialer *gomail.Dialer
	cfg    *config.Config
}

func NewService(cfg *config.Config) *Service {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	return &Service{
		dialer: dialer,
		cfg:    cfg,
	}
}

func (s *Service) Send(_ context.Context, log log.Logger, mail emailModel.Mail) error {
	m := gomail.NewMessage()

	m.SetHeader(fromHeader, s.cfg.From)
	m.SetHeader(toHeader, mail.To()...)
	m.SetHeader(subjectHeader, mail.Subject())
	m.SetBody(htmlContentType, mail.Body())

	// Send the email
	if err := s.dialer.DialAndSend(m); err != nil {
		log.Errorf("sending email: %w", err)
		return fmt.Errorf("sending mail : %w", err)
	}
	return nil
}
