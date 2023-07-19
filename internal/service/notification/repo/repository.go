package repository

import (
	"context"
	"fmt"

	"github.com/GoTaskFlow/internal/service/notification/model"
	"github.com/GoTaskFlow/internal/service/user/errors"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Add(ctx context.Context, notification []model.Notification) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("repo: newTransaction: %w", err)
	}

	query, values := buildCreateNotificationFilter(notification)
	result, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("repo: add user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repo: rowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewAddUserError()

	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) Get(ctx context.Context) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.SelectContext(ctx, &notifications, getAllQuery)
	if err != nil {
		return nil, fmt.Errorf("repo: getNotifications: %w", err)
	}
	return notifications, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (model.Notification, error) {
	var notification model.Notification
	err := r.db.GetContext(ctx, &notification, getByIdQuery, id)
	if err != nil {
		return notification, fmt.Errorf("repo: getNotifications: %w", err)
	}
	return notification, nil
}
