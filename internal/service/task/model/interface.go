package model

import (
	"context"

	userModel "github.com/GoTaskFlow/internal/service/user/model"
	"github.com/jmoiron/sqlx"
	"go.temporal.io/sdk/workflow"
)

type Service interface {
	Add(context.Context, *Task) error
	Get(context.Context) ([]Task, error)
	GetByID(context.Context, string) (*Task, error)
	TaskWorkflow(ctx workflow.Context, input *Task) error
	CreateTask(ctx context.Context, task *Task) (*string, error)
	UpdateTaskStep(ctx context.Context, taskStep *TaskStep) error
	SendMail(ctx context.Context, to []string, subject, body string) error
	NotifyDueDateWorkflow(ctx workflow.Context) error
	GetTasksWithDueDate(ctx context.Context) (map[string]userModel.User, error)
}
type Repository interface {
	Add(context.Context, *Task) (*string, error)
	Get(context.Context) ([]Task, error)
	GetByID(context.Context, string) (*Task, error)
	AddTaskStep(ctx context.Context, taskStep *TaskStep) error
	UpdateTaskStep(ctx context.Context, taskStep *TaskStep) error
	NewTransaction() (*sqlx.Tx, error)
	GetTasksWithDueDate(ctx context.Context) ([]Task, error)
}
