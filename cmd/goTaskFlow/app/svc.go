package app

import (
	"github.com/GoTaskFlow/internal/config"
	taskService "github.com/GoTaskFlow/internal/service/task"
	taskModel "github.com/GoTaskFlow/internal/service/task/model"
	taskRepo "github.com/GoTaskFlow/internal/service/task/repo"
	logModel "github.com/GoTaskFlow/pkg/logger/model"

	mailerService "github.com/GoTaskFlow/internal/notifications/mail"
	mailerModel "github.com/GoTaskFlow/internal/notifications/mail/model"
	notificationService "github.com/GoTaskFlow/internal/service/notification"
	notificationModel "github.com/GoTaskFlow/internal/service/notification/model"
	notificationRepo "github.com/GoTaskFlow/internal/service/notification/repo"
	userService "github.com/GoTaskFlow/internal/service/user"
	userModel "github.com/GoTaskFlow/internal/service/user/model"
	userRepo "github.com/GoTaskFlow/internal/service/user/repo"
	temporal "go.temporal.io/sdk/client"

	"github.com/jmoiron/sqlx"
)

type services struct {
	cfg             *config.Config
	userSvc         userModel.Service
	taskSvc         taskModel.Service
	notificationSvc notificationModel.Service
	mailerSvc       mailerModel.Service
}

type repos struct {
	userRepo         userModel.Repository
	taskRepo         taskModel.Repository
	notificationRepo notificationModel.Repository
}

func buildServices(db *sqlx.DB, temporal temporal.Client, log logModel.Logger, cfg *config.Config) *services {
	svc := &services{
		cfg: cfg,
	}
	repo := repos{}
	repo.buildRepos(db)
	svc.buildServies(repo, temporal, log)
	return svc

}
func (s *services) buildServies(repo repos, temporal temporal.Client, log logModel.Logger) {
	s.userSvc = userService.NewService(repo.userRepo)
	s.notificationSvc = notificationService.NewService(repo.notificationRepo)
	s.mailerSvc = mailerService.NewService(&s.cfg.Mailer)
	s.taskSvc = taskService.NewService(repo.taskRepo, temporal, s.notificationSvc, s.mailerSvc, log, s.userSvc)
}

func (r *repos) buildRepos(db *sqlx.DB) {
	r.userRepo = userRepo.NewRepository(db)
	r.taskRepo = taskRepo.NewRepository(db)
	r.notificationRepo = notificationRepo.NewRepository(db)
}
