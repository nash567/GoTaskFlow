package user

import (
	"context"
	"fmt"

	"github.com/GoTaskFlow/internal/service/user/model"
)

type Service struct {
	repo model.Repository
	// cfg  *config.Config
}

func NewService(repo model.Repository) model.Service {
	return &Service{repo}
}

func (s *Service) Add(ctx context.Context, user model.User) error {
	err := s.repo.Add(ctx, user)
	if err != nil {
		return fmt.Errorf("service: add user: %w ", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context) ([]model.User, error) {
	users, err := s.repo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	return users, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return user, fmt.Errorf("service: getByID: %w", err)
	}
	return user, nil
}
