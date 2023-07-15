package model

import (
	"context"
)

type Service interface {
	CreateNotification(context.Context, Notification) error
	Get(context.Context) ([]Notification, error)
	GetByID(context.Context, string) (Notification, error)
}

type Repository interface {
	Add(context.Context,  Notification) error
	Get(context.Context) ([]Notification, error)
	GetByID(context.Context, string) (Notification, error)
}
