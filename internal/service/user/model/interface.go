package model

import (
	"context"
)

type Service interface {
	Add(context.Context, User) error
	Get(context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	GetUsersByID(ctx context.Context, ids []string) ([]User, error)
}
type Repository interface {
	Add(context.Context, User) error
	Get(context.Context) ([]User, error)
	GetUserByID(ctx context.Context, filter *Filter) (User, error)
	GetUsersByID(ctx context.Context, filter *Filter) ([]User, error)
}
