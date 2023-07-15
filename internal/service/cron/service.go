package cron

import (
	"context"
	"log"
	"time"

	taskModel "github.com/GoTaskFlow/internal/service/task/model"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
)

type Service interface {
	ScheduleWorkflow(ctx context.Context)
}
type Cron struct {
	taskSvc taskModel.Service
}

func NewService(ctx context.Context, taskSvc taskModel.Service) Service {
	return &Cron{
		taskSvc: taskSvc,
	}
}

func (c *Cron) ScheduleWorkflow(ctx context.Context) {
	temporalClient, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal Client", err)
	}
	defer temporalClient.Close()

	// Create Schedule and Workflow IDs
	scheduleID := "schedule_" + uuid.New()
	workflowID := "schedule_workflow_" + uuid.New()
	// Create the schedule.
	scheduleHandle, err := temporalClient.ScheduleClient().Create(ctx, client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			Calendars: []client.ScheduleCalendarSpec{
				{
					Hour: []client.ScheduleRange{
						{
							Start: 0,
							End:   0,
						},
					},
					Minute: []client.ScheduleRange{
						{
							Start: 0,
							End:   0,
						},
					},
				},
			},
			Intervals: []client.ScheduleIntervalSpec{
				{
					Every: 24 * time.Hour,
				},
			},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  c.taskSvc.NotifyDueDateWorkflow,
			TaskQueue: "TASK_SCHEDULE_QUEUE",
		},
	})
	if err != nil {
		log.Fatalln("Unable to create schedule", err)
	}
	log.Println("Schedule created", "ScheduleID", scheduleID)
	_, _ = scheduleHandle.Describe(ctx)
}
