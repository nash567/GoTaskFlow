package repository

import (
	"context"
	"fmt"
	"reflect"

	internalErrors "github.com/GoTaskFlow/internal/errors"
	"github.com/GoTaskFlow/internal/service/user/errors"

	"github.com/GoTaskFlow/internal/service/user/model"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) model.Repository {
	return &Repository{db}
}

func (r *Repository) Get(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.db.SelectContext(ctx, &users, getAllQuery)
	if err != nil {
		return nil, fmt.Errorf("repo: get: %w", err)
	}
	return users, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, getByIdQuery, id)
	if err != nil {
		emptyUser := model.User{}
		if reflect.DeepEqual(user, emptyUser) {
			return user, internalErrors.NewInvalidIDError(id)
		}
		return user, fmt.Errorf("repo: getByID: %w", err)
	}
	return user, nil
}
func (r *Repository) Add(ctx context.Context, user model.User) error {
	result, err := r.db.ExecContext(ctx, addQuery, user.ID, user.Name, user.Email, user.Password, user.Active)
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
	return nil

}
