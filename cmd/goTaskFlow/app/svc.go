package app

import (
	taskService "github.com/GoTaskFlow/internal/service/task"
	taskModel "github.com/GoTaskFlow/internal/service/task/model"
	taskRepo "github.com/GoTaskFlow/internal/service/task/repo"

	userService "github.com/GoTaskFlow/internal/service/user"
	userModel "github.com/GoTaskFlow/internal/service/user/model"
	userRepo "github.com/GoTaskFlow/internal/service/user/repo"

	"github.com/jmoiron/sqlx"
)

type services struct {
	userSvc userModel.Service
	taskSvc taskModel.Service
}

type repos struct {
	userRepo userModel.Repository
	taskRepo taskModel.Repository
}

func buildServices(db *sqlx.DB) *services {
	svc := &services{}
	repo := repos{}
	repo.buildRepos(db)
	svc.buildServies(repo)
	return svc

}
func (s *services) buildServies(repo repos) {
	s.userSvc = userService.NewService(repo.userRepo)
	s.taskSvc = taskService.NewService(repo.taskRepo)
}

func (r *repos) buildRepos(db *sqlx.DB) {
	r.userRepo = userRepo.NewRepository(db)
	r.taskRepo = taskRepo.NewRepository(db)
}
