package model

import "context"

type Service interface {
	Add(context.Context, User) error
	Get(context.Context) ([]User, error)
	GetByID(context.Context, string) (User, error)
}
type Repository interface {
	Add(context.Context, User) error
	Get(context.Context) ([]User, error)
	GetByID(context.Context, string) (User, error)
}
