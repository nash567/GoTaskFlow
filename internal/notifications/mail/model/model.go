package model

import (
	"context"

	log "github.com/GoTaskFlow/pkg/logger/model"
)

type Mailer interface {
	To() []string
	Subject() string
	Body() string
}

type Service interface {
	Send(_ context.Context, log log.Logger, mail Mail) error
}

type Mail struct {
	to      []string
	subject string
	body    string
}

func NewMail(to []string, subject string, body string) Mail {
	return Mail{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (m *Mail) To() []string {
	return m.to
}

func (m *Mail) Subject() string {
	return m.subject
}

func (m *Mail) Body() string {
	return m.body
}
