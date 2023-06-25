package repository

import (
	"context"
	"fmt"

	internalErrors "github.com/GoTaskFlow/internal/errors"
	taskError "github.com/GoTaskFlow/internal/service/task/errors"
	"github.com/GoTaskFlow/internal/service/task/model"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) model.Repository {
	return &Repository{db}
}

func (r *Repository) Get(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.SelectContext(ctx, &tasks, getQuery)
	if err != nil {
		return nil, fmt.Errorf("repo: get: %w", err)
	}
	return tasks, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*model.Task, error) {
	var task *model.Task
	err := r.db.GetContext(ctx, task, getByIdQuery)
	if err != nil {
		fmt.Println("asdada", task)
		if task == nil {
			return nil, internalErrors.NewInvalidIDError(id)
		}
		return nil, fmt.Errorf("repo: getById: %w", err)
	}
	return task, nil

}
func (r *Repository) Add(ctx context.Context, task model.Task) error {
	reult, err := r.db.ExecContext(
		ctx,
		addQuery,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.AssignedBy,
		task.AssignedTo,
		pq.Array(task.Comments),
		task.DueDate,
	)
	if err != nil {
		return fmt.Errorf("repo: add: %w", err)
	}

	rowsAffected, err := reult.RowsAffected()
	if err != nil {
		return fmt.Errorf("repo: rowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return taskError.NewAddTaskError()

	}
	return nil

}
