package task

import (
	"context"
	"fmt"

	mailerModel "github.com/GoTaskFlow/internal/notifications/mail/model"
	notificationModel "github.com/GoTaskFlow/internal/service/notification/model"
	userModel "github.com/GoTaskFlow/internal/service/user/model"
	"github.com/google/uuid"

	"github.com/GoTaskFlow/internal/service/task/model"
	logModel "github.com/GoTaskFlow/pkg/logger/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jmoiron/sqlx"
	temporal "go.temporal.io/sdk/client"
	temporalClient "go.temporal.io/sdk/client"
)

type Service struct {
	log             logModel.Logger
	repo            model.Repository
	temporal        temporal.Client
	notificationSvc notificationModel.Service
	mailerSvc       mailerModel.Service
	userSvc         userModel.Service
}

func NewService(
	repo model.Repository,
	temporal temporal.Client,
	notification notificationModel.Service,
	mailerSvc mailerModel.Service,
	log logModel.Logger,
	userSvc userModel.Service,
) model.Service {
	return &Service{log, repo, temporal, notification, mailerSvc, userSvc}
}

func (s *Service) Add(ctx context.Context, task *model.Task) error {
	var (
		tx         *sqlx.Tx
		err        error
		workflowID *string
		runID      *string
	)
	defer func() {
		if err != nil && tx != nil {
			tx.Rollback()
			err := s.temporal.CancelWorkflow(ctx, aws.StringValue(workflowID), aws.StringValue(runID))
			if err != nil {
				s.log.Errorf("failed to cancel workflow", err)
			}
		}

	}()

	tx, err = s.repo.NewTransaction()
	if err != nil {
		return fmt.Errorf("service: failed to create transaction: %w", err)
	}
	err = s.PrepareTask(ctx, tx, task.ID)
	if err != nil {
		return fmt.Errorf("service: prepareTask: %w", err)
	}
	workflowID, runID, err = s.ExecuteTaskWorkflow(task)
	if err != nil {
		return fmt.Errorf("service: add: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("service: transactionCommit: %w", err)
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

func (s *Service) GetTaskByID(ctx context.Context, id string) (*model.Task, error) {
	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		return task, fmt.Errorf("service: getById: %w", err)
	}
	return task, nil
}

func (s *Service) UpdateTask(ctx context.Context, task *model.UpdateTask) error {
	workflowID := uuid.NewString()
	wctx := context.Background()

	// taskWorkflow
	o, err := s.temporal.ExecuteWorkflow(
		wctx,
		temporalClient.StartWorkflowOptions{
			ID:        workflowID,
			TaskQueue: "TASK_WORKER_QUEUE",
		},
		s.UpdateTaskWorkflow,
		task,
	)
	if err != nil {
		return fmt.Errorf("service: updateTask: %w", err)
	}
	if err = o.Get(wctx, nil); err != nil {
		return fmt.Errorf("service: updateTask: get %w", err)
	}

	return nil

}

func (s *Service) UpdateTaskActivity(ctx context.Context, input *model.UpdateTask) error {
	err := s.repo.UpdateTask(ctx, input)
	if err != nil {
		return fmt.Errorf("service: updateTask: %w", err)
	}
	return nil
}
