package notification

import (
	"context"
	"fmt"

	"github.com/GoTaskFlow/internal/service/notification/model"
)

type Service struct {
	repo model.Repository
}

func NewService(repo model.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateNotification(ctx context.Context, notification model.Notification) error {
	err := s.repo.Add(ctx,  notification)
	if err != nil {
		return fmt.Errorf("service: addNotification: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context) ([]model.Notification, error) {
	notifications, err := s.repo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: getNotifications: %w", err)
	}
	return notifications, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (model.Notification, error) {
	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return notification, fmt.Errorf("service: getNotifications: %w", err)
	}
	return notification, nil
}
