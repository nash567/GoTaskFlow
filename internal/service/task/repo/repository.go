package repository

import (
	"context"
	"fmt"

	taskError "github.com/GoTaskFlow/internal/service/task/errors"
	"github.com/GoTaskFlow/internal/service/task/model"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) model.Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.SelectContext(ctx, &tasks, getQuery)
	if err != nil {
		return nil, fmt.Errorf("repo: get: %w", err)
	}
	return tasks, nil
}

func (r *Repository) GetTaskByID(ctx context.Context, id string) (*model.Task, error) {
	var task model.Task
	err := r.db.GetContext(ctx, &task, getByIdQuery, id)
	if err != nil {
		return nil, fmt.Errorf("repo: getById: %w", err)
	}
	return &task, nil

}
func (r *Repository) Add(ctx context.Context, task *model.Task) (*string, error) {
	tx, err := r.NewTransaction()
	if err != nil {
		return nil, fmt.Errorf("repo: newTransaction: %w", err)
	}

	defer func() {
		if err != nil && tx != nil {
			if err = tx.Rollback(); err != nil {
				fmt.Println("i m rolling back transaction", err)
			}
		}
	}()
	result, err := tx.ExecContext(
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
		return nil, fmt.Errorf("repo: add: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("repo: rowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, taskError.NewAddTaskError()

	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &task.ID, nil

}
func (r *Repository) AddTaskStep(ctx context.Context, taskStep *model.TaskStep) error {
	result, err := r.db.ExecContext(ctx,
		createTaskStepQuery,
		taskStep.Name,
		taskStep.Status,
		taskStep.StepOrder,
		taskStep.TaskID,
	)
	if err != nil {
		return fmt.Errorf("repo: createTaskStep %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repo: rowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return taskError.NewAddTaskStepError()

	}
	return nil
}

func (r *Repository) UpdateTask(ctx context.Context, filter *model.UpdateTask) error {
	result, err := r.db.ExecContext(ctx, buildUpdateTaskFilter(filter))
	if err != nil {
		return fmt.Errorf("repo: updateTask %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repo: rowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return taskError.NewUpdateTaskError()

	}
	return nil
}

func (r *Repository) UpdateTaskStep(ctx context.Context, taskStep *model.TaskStep) error {
	tx, err := r.NewTransaction()
	if err != nil {
		return fmt.Errorf("repo: newTransaction: %w", err)
	}
	defer func() {
		if err != nil && tx != nil {
			if err = tx.Rollback(); err != nil {
				fmt.Println("i m rolling back transaction", err)
			}
		}
	}()

	result, err := tx.ExecContext(ctx, updateTaskStepQuery, taskStep.Status, taskStep.TaskID, taskStep.Name)
	if err != nil {
		return fmt.Errorf("repo: createTaskStep %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repo: rowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return taskError.NewUpdateTaskStepError()

	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) NewTransaction() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *Repository) GetTasksWithDueDate(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.SelectContext(ctx, &tasks, getTasksWithDueDate)
	if err != nil {
		return nil, fmt.Errorf("repo: getUserIDsWithDueDate: %w", err)
	}

	return tasks, nil

}
