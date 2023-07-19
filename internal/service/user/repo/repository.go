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
	err := r.db.SelectContext(ctx, &users, getUsers)
	if err != nil {
		return nil, fmt.Errorf("repo: get: %w", err)
	}
	return users, nil
}

func (r *Repository) GetUserByID(ctx context.Context, filter *model.Filter) (model.User, error) {
	var user model.User
	filterQueries, values := buildFilter(filter)
	err := r.db.GetContext(ctx, &user, getUsers+filterQueries, values...)
	if err != nil {
		// TODO: check if this check is needed or not
		emptyUser := model.User{}
		if reflect.DeepEqual(user, emptyUser) {
			return user, internalErrors.NewInvalidIDError(filter.ID[0])
		}
		return user, fmt.Errorf("repo: getUserByID: %w", err)
	}
	return user, nil
}
func (r *Repository) Add(ctx context.Context, user *model.User) error {
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
func (r *Repository) GetUsersByID(ctx context.Context, filter *model.Filter) ([]model.User, error) {
	var users []model.User
	filterQueries, values := buildFilter(filter)
	err := r.db.SelectContext(ctx, &users, getUsers+filterQueries, values...)
	if err != nil {
		return users, fmt.Errorf("repo: getUsersByID: %w", err)
	}
	return users, nil
}
