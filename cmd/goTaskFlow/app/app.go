package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/GoTaskFlow/api/http/task"
	"github.com/GoTaskFlow/api/http/user"

	"github.com/GoTaskFlow/internal/config"
	"github.com/GoTaskFlow/pkg/db"
	"github.com/GoTaskFlow/pkg/logger"
	logModel "github.com/GoTaskFlow/pkg/logger/model"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	temporal "go.temporal.io/sdk/client"
)

const timeout = 5 * time.Second

type Application struct {
	db             *sqlx.DB
	httpServer     *http.Server
	cfg            *config.Config
	router         *mux.Router
	log            logModel.Logger
	services       *services
	temporalClient temporal.Client
}

func (a *Application) Init(ctx context.Context, configFile string, migrationPath string, seedDataPath string) {
	config, err := config.Load(configFile)
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

	db, err := db.NewConnection(&config.DB)
	if err != nil {
		a.log.WithError(err).Fatal("error connecting to db")
		return
	}
	a.db = db
	a.log.WithField("host", a.cfg.DB.Host).WithField("port", a.cfg.DB.Port).Info("created database connection successfully")

	a.router = mux.NewRouter()

	temporalOptions := temporal.Options{
		HostPort: fmt.Sprintf("%s:%s", a.cfg.Temporal.Host, a.cfg.Temporal.Port),
	}
	a.temporalClient, err = temporal.Dial(temporalOptions)
	if err != nil {
		a.log.Fatalf("temporal client: %w", err)

	}
	a.services = buildServices(a.db, a.temporalClient, a.log, a.cfg)
	a.setupHandlers()
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

func (a *Application) setupHandlers() {
	user.RegisterHandlers(a.router, a.services.userSvc, a.log)
	task.RegisterHandlers(a.router, a.services.taskSvc, a.log)
}
