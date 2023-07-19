package task

import (
	"context"
	"fmt"
	"reflect"
	"time"

	mailerModel "github.com/GoTaskFlow/internal/notifications/mail/model"
	notificationModel "github.com/GoTaskFlow/internal/service/notification/model"
	"github.com/GoTaskFlow/internal/service/task/model"
	userModel "github.com/GoTaskFlow/internal/service/user/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	temporalClient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	createTask          = "CREATE TASK"
	createNotification  = "CREATE NOTIFICATION"
	getUserByID         = "GET USER BY ID"
	sendMail            = "SEND EMAIL"
	inQueue             = "IN_QUEUE"
	failed              = "FAILED"
	inProgress          = "IN_PROGRESS"
	completed           = "COMPLETED"
	NotifyDueDate       = "NOTIFY_DUE_DATE"
	GetUsersWithDueDate = "GET_USERS_WITH_DUE_DATE"
	unread              = "UNREAD"
)

func (s *Service) ExecuteTaskWorkflow(task *model.Task) (*string, *string, error) {
	workflowID := uuid.NewString()
	wctx := context.Background()
	// taskWorkflow
	o, err := s.temporal.ExecuteWorkflow(
		wctx,
		temporalClient.StartWorkflowOptions{
			ID:        workflowID,
			TaskQueue: "TASK_WORKER_QUEUE",
		},
		s.TaskWorkflow,
		task,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("service: executeTaskWorkflow: %w", err)
	}
	if err = o.Get(wctx, nil); err != nil {
		return nil, nil, fmt.Errorf("service: executeTaskWorkflow get: %w", err)
	}

	return &workflowID, aws.String(o.GetRunID()), nil

}
func (s *Service) TaskWorkflow(ctx workflow.Context, input *model.Task) error {
	var (
		retryPolicy = temporal.RetryPolicy{
			InitialInterval:    time.Minute,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    1,
		}

		opts = workflow.ActivityOptions{
			StartToCloseTimeout: time.Hour,
			RetryPolicy:         &retryPolicy,
		}

		stepName    *string
		err         error
		mailSubject = "Task Assignment"
	)
	ctx = workflow.WithActivityOptions(ctx, opts)

	sessionOption := &workflow.SessionOptions{
		CreationTimeout:  time.Minute,
		ExecutionTimeout: time.Hour * 3,
	}
	ctx, err = workflow.CreateSession(ctx, sessionOption)
	if err != nil {
		return fmt.Errorf("taskWorkflow: createSession: %w", err)
	}

	defer func(ctx workflow.Context) {
		defer workflow.CompleteSession(ctx)
		if err != nil {
			if stepName != nil {
				dctx, _ := workflow.NewDisconnectedContext(ctx)
				f := workflow.ExecuteActivity(
					dctx, s.UpdateTaskStep,
					&model.TaskStep{
						TaskID: &input.ID,
						Name:   stepName,
						Status: aws.String(failed),
					},
				)
				updateErr := f.Get(dctx, nil)
				if updateErr != nil {
					err = updateErr

				}
			}
		}
	}(ctx)

	// create task
	stepName = aws.String(createTask)
	createTaskFuture := workflow.ExecuteActivity(
		ctx, s.CreateTask, input,
	)
	err = s.updateStepStatus(ctx, &input.ID, stepName, aws.String(inProgress))
	if err != nil {
		return fmt.Errorf("taskWorkflow: updateStepStatus stepName:%s: %w", aws.StringValue(stepName), err)

	}
	var taskID *string
	if err = createTaskFuture.Get(ctx, &taskID); err != nil {
		return fmt.Errorf("taskWorkflow: createTaskFuture get: %w", err)
	}
	err = s.updateStepStatus(ctx, &input.ID, stepName, aws.String(completed))
	if err != nil {
		return fmt.Errorf("taskWorkflow: updateStepStatus stepName:%s: %w", aws.StringValue(stepName), err)
	}

	// create notification  for user
	stepName = aws.String(createNotification)
	createNotificationFuture := workflow.ExecuteActivity(ctx, "CreateNotification", []notificationModel.Notification{
		{
			Message: "task created successfully",
			UserID:  input.AssignedTo,
			TaskID:  aws.StringValue(taskID),
			Status:  notificationModel.StatusUnread,
		},
	})
	err = s.updateStepStatus(ctx, &input.ID, stepName, aws.String(inProgress))
	if err != nil {
		return fmt.Errorf("taskWorkflow: updateStepStatus stepName:%s: %w", aws.StringValue(stepName), err)

	}
	if err = createNotificationFuture.Get(ctx, nil); err != nil {
		// TODO : delete task created if error occured here using compensating activities
		return fmt.Errorf("taskWorkflow: createNotificationFuture get: %w", err)
	}
	err = s.updateStepStatus(ctx, &input.ID, stepName, aws.String(completed))
	if err != nil {
		return fmt.Errorf("taskWorkflow: updateStepStatus stepName:%s: %w", aws.StringValue(stepName), err)

	}

	// get userby id
	stepName = aws.String(getUserByID)
	var user userModel.User
	getUserByIDFuture := workflow.ExecuteActivity(ctx, "GetUserByID", input.AssignedTo)
	err = s.updateStepStatus(ctx, &input.ID, stepName, aws.String(inProgress))
	if err != nil {
		return fmt.Errorf("taskWorkflow: updateStepStatus stepName:%s: %w", aws.StringValue(stepName), err)
	}
	if err = getUserByIDFuture.Get(ctx, &user); err != nil {
		return fmt.Errorf("taskWorkflow: getUserByIDFuture get: %w", err)
	}
	err = s.updateStepStatus(ctx, &input.ID, stepName, aws.String(completed))
	if err != nil {
		return fmt.Errorf("taskWorkflow: updateStepStatus stepName:%s: %w", aws.StringValue(stepName), err)
	}

	//  send invitation mail to user
	stepName = aws.String(sendMail)
	sendMailFuture := workflow.ExecuteActivity(ctx, s.SendMail, []string{user.Email}, mailSubject, s.getBody(createTaskMail))
	err = s.updateStepStatus(ctx, &input.ID, stepName, aws.String(inProgress))
	if err != nil {
		return fmt.Errorf("taskWorkflow: updateStepStatus stepName:%s: %w", aws.StringValue(stepName), err)
	}
	if err = sendMailFuture.Get(ctx, nil); err != nil {
		return fmt.Errorf("taskWorkflow: sendMailFuture get: %w", err)
	}
	err = s.updateStepStatus(ctx, &input.ID, stepName, aws.String(completed))
	if err != nil {
		return fmt.Errorf("taskWorkflow: updateStepStatus stepName:%s: %w", aws.StringValue(stepName), err)
	}

	return err

}

func (s *Service) CreateTask(ctx context.Context, task *model.Task) (*string, error) {
	// task.ID = uuid.NewString()
	id, err := s.repo.Add(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("service: createTask: %v", err)
	}
	return id, nil
}
func (s *Service) UpdateTaskStep(ctx context.Context, taskStep *model.TaskStep) error {
	err := s.repo.UpdateTaskStep(ctx, taskStep)
	if err != nil {
		return fmt.Errorf("service: createTask: %v", err)
	}
	return nil
}

func (s *Service) PrepareTask(ctx context.Context, tx *sqlx.Tx, taskID string) error {
	err := s.PrepareTaskSteps(ctx, tx, taskID)
	if err != nil {
		return fmt.Errorf("workFlow: prepareTask: %w", err)
	}
	return err
}

func (s *Service) updateStepStatus(ctx workflow.Context, id, stepName, status *string) error {
	f := workflow.ExecuteActivity(ctx,
		s.UpdateTaskStep,
		&model.TaskStep{
			TaskID: id,
			Name:   stepName,
			Status: status,
		})

	err := f.Get(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) SendMail(ctx context.Context, to []string, subject, body string) error {
	mail := mailerModel.NewMail(to, subject, body)
	err := s.mailerSvc.Send(ctx, s.log, mail)
	if err != nil {
		return fmt.Errorf("workflow: SendMail: %w", err)
	}
	s.log.Info("Mail sent successfully")
	return nil
}

func (s *Service) NotifyDueDateWorkflow(ctx workflow.Context) error {
	var (
		retryPolicy = temporal.RetryPolicy{
			InitialInterval:    time.Minute,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		}

		opts = workflow.ActivityOptions{
			StartToCloseTimeout: time.Hour,
			RetryPolicy:         &retryPolicy,
		}

		stepName    *string
		err         error
		mailSubject = "Task Deadline"
	)
	ctx = workflow.WithActivityOptions(ctx, opts)

	sessionOption := &workflow.SessionOptions{
		CreationTimeout:  time.Minute,
		ExecutionTimeout: time.Hour * 3,
	}
	defer func(ctx workflow.Context) {
		defer workflow.CompleteSession(ctx)
	}(ctx)
	ctx, err = workflow.CreateSession(ctx, sessionOption)
	if err != nil {
		return err
	}
	stepName = aws.String(GetUsersWithDueDate)
	userMap := make(map[string]userModel.User)
	future := workflow.ExecuteActivity(ctx, s.GetTasksWithDueDate)
	// TODO: check how to update the step status
	if err = future.Get(ctx, &userMap); err != nil {
		s.log.Errorf("notifyUserWorkflow: notifyDueDateFuture get: %w", stepName, err)
		return err
	}

	stepName = aws.String(NotifyDueDate)
	userMails := make([]string, 0, len(userMap))
	for _, value := range userMap {
		userMails = append(userMails, value.Email)
	}

	future = workflow.ExecuteActivity(ctx, s.SendMail, userMails, mailSubject, s.getBody(notifyDueDate))
	if err = future.Get(ctx, nil); err != nil {
		s.log.Errorf("notifyUserWorkflow: notifyDueDateFuture get: %w", aws.StringValue(stepName), err)
		return err
	}

	return nil
}

func (s *Service) PrepareTaskSteps(ctx context.Context, tx *sqlx.Tx, taskID string) error {
	steps := []*model.TaskStep{
		{
			Name:      aws.String(createTask),
			Status:    aws.String(inQueue),
			StepOrder: aws.Int(1),
			TaskID:    aws.String(taskID),
		},
		{
			Name:      aws.String(createNotification),
			Status:    aws.String(inQueue),
			StepOrder: aws.Int(2),
			TaskID:    aws.String(taskID),
		},
		{
			Name:      aws.String(getUserByID),
			Status:    aws.String(inQueue),
			StepOrder: aws.Int(3),
			TaskID:    aws.String(taskID),
		},
		{
			Name:      aws.String(sendMail),
			Status:    aws.String(inQueue),
			StepOrder: aws.Int(4),
			TaskID:    aws.String(taskID),
		},
	}

	for _, step := range steps {
		err := s.repo.AddTaskStep(ctx, step)
		if err != nil {
			return fmt.Errorf("taskSerice: prepareTaskSteps: addTaskStep: %w", err)
		}
	}
	return nil

}

func (s *Service) GetTasksWithDueDate(ctx context.Context) (map[string]userModel.User, error) {
	tasks, err := s.repo.GetTasksWithDueDate(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: getTasksWithDueDate: %w", err)
	}
	userIDs := make([]string, 0, len(tasks))
	for _, task := range tasks {
		userIDs = append(userIDs, task.AssignedTo)
	}

	users, err := s.userSvc.GetUsersByID(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("service: getUsersByID: %w", err)
	}

	userMap := make(map[string]userModel.User)

	for _, task := range tasks {
		for _, user := range users {
			if user.ID == task.AssignedTo {
				userMap[task.ID] = user
			}
		}
	}

	return userMap, nil

}

func (s *Service) UpdateTaskWorkflow(ctx workflow.Context, input *model.UpdateTask) error {
	var (
		retryPolicy = temporal.RetryPolicy{
			InitialInterval:    time.Minute,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    1,
		}

		opts = workflow.ActivityOptions{
			StartToCloseTimeout: time.Hour,
			RetryPolicy:         &retryPolicy,
		}

		err         error
		mailSubject = "Task Update"
	)
	ctx = workflow.WithActivityOptions(ctx, opts)

	sessionOption := &workflow.SessionOptions{
		CreationTimeout:  time.Minute,
		ExecutionTimeout: time.Hour * 3,
	}
	ctx, err = workflow.CreateSession(ctx, sessionOption)
	if err != nil {
		return fmt.Errorf("UpdateTaskWorkflow: createSession: %w", err)
	}

	defer func(ctx workflow.Context) {
		defer workflow.CompleteSession(ctx)
	}(ctx)

	// get task by id
	var originalTask *model.Task
	getTaskFuture := workflow.ExecuteActivity(ctx, s.GetTaskByID, input.ID)
	if err = getTaskFuture.Get(ctx, &originalTask); err != nil {
		return fmt.Errorf("UpdateTaskWorkflow: getTaskFuture get: %w", err)
	}

	// prepare task for update
	updateTaskInput := input.PrepareUpateTask(originalTask)

	// update task
	updateTaskFuture := workflow.ExecuteActivity(ctx, s.UpdateTaskActivity, updateTaskInput)
	if err = updateTaskFuture.Get(ctx, nil); err != nil {
		return fmt.Errorf("UpdateTaskWorkflow: updateTaskFuture get: %w", err)
	}

	// create notification for each field update
	assignedTo := originalTask.AssignedTo
	if updateTaskInput.AssignedTo != "" {
		assignedTo = updateTaskInput.AssignedTo
	}

	notifications := s.prepareNotifications(updateTaskInput, assignedTo)
	createNotificationFuture := workflow.ExecuteActivity(ctx, "CreateNotification", notifications)
	if err = createNotificationFuture.Get(ctx, nil); err != nil {
		// TODO : delete task created if error occured here using compensating activities
		return fmt.Errorf("taskWorkflow: createNotificationFuture get: %w", err)
	}
	// get user by id
	var users []userModel.User
	getUserByIDFuture := workflow.ExecuteActivity(ctx, "GetUsersByID", []string{assignedTo, originalTask.AssignedTo})
	if err = getUserByIDFuture.Get(ctx, &users); err != nil {
		return fmt.Errorf("taskWorkflow: getUserByIDFuture get: %w", err)
	}
	var userEmails []string
	for _, user := range users {
		userEmails = append(userEmails, user.Email)
	}

	if updateTaskInput.AssignedTo != "" && updateTaskInput.AssignedTo != originalTask.AssignedTo {
		if len(userEmails) <= 1 {
			// send mail to user that task has been updated
			sendMailFuture := workflow.ExecuteActivity(ctx, s.SendMail, []string{userEmails[0], userEmails[1]}, mailSubject, s.getBody(createTaskMail))
			if err = sendMailFuture.Get(ctx, nil); err != nil {
				return fmt.Errorf("taskWorkflow: sendMailFuture get: %w", err)
			}
		}
	}

	sendMailFuture := workflow.ExecuteActivity(ctx, s.SendMail, []string{userEmails[0]}, mailSubject, s.getBody(updateTaskMail))
	if err = sendMailFuture.Get(ctx, nil); err != nil {
		return fmt.Errorf("taskWorkflow: sendMailFuture get: %w", err)
	}

	return err
}

func (s *Service) prepareNotifications(task *model.UpdateTask, userID string) []notificationModel.Notification {
	v := reflect.ValueOf(*task)
	t := reflect.TypeOf(*task)

	var notifications []notificationModel.Notification

	for i := 0; i < v.NumField(); i++ {
		// Get the value of the current field.
		fieldName := t.Field(i).Name
		fieldValue := v.Field(i).Interface()
		if !reflect.DeepEqual(fieldValue, reflect.Zero(t.Field(i).Type).Interface()) && fieldName != "ID" {
			notifications = append(notifications, notificationModel.Notification{
				Message: fmt.Sprintf(`Task Field, %s has been updated to %v`, fieldName, fieldValue),
				UserID:  userID,
				TaskID:  task.ID,
				Status:  notificationModel.StatusUnread,
			})

		}
	}

	return notifications

}
