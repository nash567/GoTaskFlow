package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/GoTaskFlow/internal/config"
	"github.com/GoTaskFlow/pkg/db"
	"github.com/GoTaskFlow/pkg/logger"
	logModel "github.com/GoTaskFlow/pkg/logger/model"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
)

const timeout = 5 * time.Second

type Application struct {
	db         *sqlx.DB
	httpServer *http.Server
	cfg        *config.Config
	router     *mux.Router
	log        logModel.Logger
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

	router := mux.NewRouter()
	a.router = router
}

func (a *Application) Start(ctx context.Context) {
	a.router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
		AllowedHeaders:   []string{"accept", "Authorization", "content-type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler)
	a.httpServer = &http.Server{
		Addr:              ":" + fmt.Sprintf("%v", a.cfg.Server.Port),
		Handler:           a.router,
		ReadHeaderTimeout: timeout,
	}
	go func() {
		defer a.log.Warn("server stopped listening...")

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
