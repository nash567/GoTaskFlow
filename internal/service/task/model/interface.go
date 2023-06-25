package model

import "context"

type Service interface {
	Add(context.Context, Task) error
	Get(context.Context) ([]Task, error)
	GetByID(context.Context, string) (*Task, error)
}
type Repository interface {
	Add(context.Context, Task) error
	Get(context.Context) ([]Task, error)
	GetByID(context.Context, string) (*Task, error)
}
