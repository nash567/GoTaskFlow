package app

import (
	"context"
	"fmt"
	"log"

	"github.com/GoTaskFlow/pkg/db"
	"github.com/GoTaskFlow/pkg/logger"
	"github.com/jmoiron/sqlx"

	mailerService "github.com/GoTaskFlow/internal/notifications/mail"
	notificationService "github.com/GoTaskFlow/internal/service/notification"
	taskService "github.com/GoTaskFlow/internal/service/task"
	taskModel "github.com/GoTaskFlow/internal/service/task/model"
	taskRepo "github.com/GoTaskFlow/internal/service/task/repo"
	userService "github.com/GoTaskFlow/internal/service/user"
	logModel "github.com/GoTaskFlow/pkg/logger/model"
	temporal "go.temporal.io/sdk/client"

	cron "github.com/GoTaskFlow/internal/service/cron"
	notificationModel "github.com/GoTaskFlow/internal/service/notification/model"
	notificationRepo "github.com/GoTaskFlow/internal/service/notification/repo"
	userModel "github.com/GoTaskFlow/internal/service/user/model"
	userRepo "github.com/GoTaskFlow/internal/service/user/repo"
)

type Application struct {
	cfg             *Config
	log             logModel.Logger
	taskSvc         taskModel.Service
	db              *sqlx.DB
	temporalClient  temporal.Client
	notificationSvc notificationModel.Service
	mailerSvc       *mailerService.Service
	userSvc         userModel.Service
	cronSvc         cron.Service
}

func (a *Application) Init(ctx context.Context, configFile string, migrationPath string, seedDataPath string) {
	config, err := Load(configFile)
	if err != nil {
		log.Fatal("failed to read config")
		return
	}
	a.cfg = config
	a.log, err = logger.NewZapLogger(&a.cfg.Log)
	if err != nil {
		panic(err)
	}
	db, err := db.NewConnection(&config.DB)
	if err != nil {
		a.log.WithError(err).Fatal("error connecting to db")
		return
	}
	a.db = db
	temporalOptions := temporal.Options{
		HostPort: fmt.Sprintf("%s:%s", a.cfg.Temporal.Host, a.cfg.Temporal.Port),
	}
	a.temporalClient, err = temporal.Dial(temporalOptions)
	if err != nil {
		a.log.Fatalf("temporal client: %w", err)

	}
	notificationRepo := notificationRepo.NewRepository(a.db)
	a.notificationSvc = notificationService.NewService(notificationRepo)
	a.mailerSvc = mailerService.NewService(&a.cfg.Mailer)
	a.userSvc = userService.NewService(userRepo.NewRepository(a.db))
	a.taskSvc = taskService.NewService(taskRepo.NewRepository(a.db), a.temporalClient, a.notificationSvc, a.mailerSvc, a.log, a.userSvc)
	a.cronSvc = cron.NewService(ctx, a.taskSvc)
}

func (a *Application) Start(ctx context.Context) {
	fmt.Println("scheduling workflow................................")
	a.cronSvc.ScheduleWorkflow(ctx)

}
