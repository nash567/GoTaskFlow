package task

import (
	"context"
	"fmt"

	"github.com/GoTaskFlow/internal/service/task/model"
)

type Service struct {
	repo model.Repository
}

func NewService(repo model.Repository) model.Service {
	return &Service{repo}
}

func (s *Service) Add(ctx context.Context, task model.Task) error {
	err := s.repo.Add(ctx, task)
	if err != nil {
		return fmt.Errorf("service: add: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context) ([]model.Task, error) {
	tasks, err := s.repo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: get: %w", err)
	}

	return tasks, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*model.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return task, fmt.Errorf("service: getById: %w", err)
	}
	return task, nil
}
