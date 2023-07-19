package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	taskService "github.com/GoTaskFlow/internal/service/task"
	taskModel "github.com/GoTaskFlow/internal/service/task/model"
	taskRepo "github.com/GoTaskFlow/internal/service/task/repo"
	"github.com/GoTaskFlow/pkg/db"
	"github.com/GoTaskFlow/pkg/logger"
	logModel "github.com/GoTaskFlow/pkg/logger/model"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	mailerService "github.com/GoTaskFlow/internal/notifications/mail"
	notificationService "github.com/GoTaskFlow/internal/service/notification"
	userService "github.com/GoTaskFlow/internal/service/user"

	notificationModel "github.com/GoTaskFlow/internal/service/notification/model"
	notificationRepo "github.com/GoTaskFlow/internal/service/notification/repo"
	userModel "github.com/GoTaskFlow/internal/service/user/model"
	userRepo "github.com/GoTaskFlow/internal/service/user/repo"
	temporal "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const timeout = 5 * time.Second

type Application struct {
	httpServer      *http.Server
	router          *mux.Router
	db              *sqlx.DB
	cfg             *Config
	log             logModel.Logger
	temporalClient  temporal.Client
	taskService     taskModel.Service
	notificationSvc notificationModel.Service
	mailerSvc       *mailerService.Service
	userSvc         userModel.Service
}

func (a *Application) Init(ctx context.Context, configFile string) {
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
	a.log = a.log.WithFields(logModel.Fields{
		"appName": a.cfg.AppName,
		"env":     a.cfg.Env,
	})
	a.router = mux.NewRouter()
	a.setupRoutes()
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
}

func (a *Application) Start(ctx context.Context) {
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Accept", "Authorization", "Content-Type"}),
	)
	a.router.Use(corsHandler)

	a.httpServer = &http.Server{
		Addr:              ":" + fmt.Sprintf("%v", a.cfg.Server.Port),
		Handler:           a.router,
		ReadHeaderTimeout: timeout,
	}
	go func() {
		defer a.log.Error("server stopped listening...")
		if err := a.httpServer.ListenAndServe(); err != nil {
			a.log.WithError(err).Fatal("failed to listen and serve")
			return
		}
	}()
	a.log.Infof("http server started on port %d..", a.cfg.Server.Port)
}
func (a *Application) Stop(ctx context.Context) {
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
	a.log.Warn("shutting down....")
}

func (a *Application) setupRoutes() {
	a.router.HandleFunc("taskWorker/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func (a *Application) RegisterWorkflow() worker.Worker {

	// Todo find from where to get workflowid and other params

	// _, err := a.temporalClient.ResetWorkflowExecution(context.TODO(), &workflowservice.ResetWorkflowExecutionRequest{
	// 	Namespace: "default",
	// })
	// if err != nil {
	// 	a.log.Errorf("temporal resetWorkflow: %v", err)
	// }

	w := worker.New(a.temporalClient, a.cfg.Temporal.TaskWorkerQueue, worker.Options{
		MaxConcurrentActivityExecutionSize: 300,
		EnableSessionWorker:                true,
	})
	notificationRepo := notificationRepo.NewRepository(a.db)
	a.notificationSvc = notificationService.NewService(notificationRepo)
	a.mailerSvc = mailerService.NewService(&a.cfg.Mailer)
	a.userSvc = userService.NewService(userRepo.NewRepository(a.db))
	a.taskService = taskService.NewService(taskRepo.NewRepository(a.db), a.temporalClient, a.notificationSvc, a.mailerSvc, a.log, a.userSvc)
	w.RegisterWorkflowWithOptions(a.taskService.UpdateTaskWorkflow, workflow.RegisterOptions{
		Name: "updateTaskWorkflow",
	})
	w.RegisterWorkflowWithOptions(a.taskService.TaskWorkflow, workflow.RegisterOptions{
		Name: "taskWorkflow",
	})
	w.RegisterActivity(a.taskService.CreateTask)
	w.RegisterActivity(a.taskService.UpdateTaskStep)
	w.RegisterActivity(a.notificationSvc.CreateNotification)
	w.RegisterActivity(a.taskService.SendMail)
	w.RegisterActivity(a.userSvc.GetUserByID)
	w.RegisterActivity(a.taskService.UpdateTaskActivity)
	w.RegisterActivity(a.taskService.GetTaskByID)
	w.RegisterActivity(a.userSvc.GetUsersByID)

	return w
}

func (a *Application) RunWorker(ctx context.Context, w worker.Worker) {
	defer func() {
		a.Stop(ctx)
		a.temporalClient.Close()
	}()
	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("failed to start temporal worker: %v", err)
	}
}
