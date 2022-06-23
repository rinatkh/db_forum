package api

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	controllers "github.com/rinatkh/db_forum/internal/api/contollers"
	"github.com/rinatkh/db_forum/internal/db"
	"github.com/rinatkh/db_forum/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type APIService struct {
	log    *logrus.Entry
	router *echo.Echo
	debug  bool
}

func (svc *APIService) Serve() {
	svc.log.Info("Starting HTTP server")
	listenAddr := viper.GetString("service.bind.address") + ":" + viper.GetString("service.bind.port")
	svc.log.Fatal(svc.router.Start(listenAddr))
}

func (svc *APIService) Shutdown(ctx context.Context) error {
	if err := svc.router.Shutdown(ctx); err != nil {
		svc.log.Fatal(err)
	}
	return nil
}

func NewAPIService(log *logrus.Entry, dbConn *pgxpool.Pool, debug bool) (*APIService, error) {
	svc := &APIService{
		log:    log,
		router: echo.New(),
		debug:  debug,
	}

	svc.router.Validator = NewValidator()
	svc.router.Binder = NewBinder()

	repository, err := db.NewRepository(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	registry := service.NewRegistry(log, repository)
	userCtrl := controllers.NewUserController(log, registry)

	svc.router.HTTPErrorHandler = svc.httpErrorHandler
	svc.router.Use(svc.XRequestIDMiddleware(), svc.LoggingMiddleware())

	api := svc.router.Group("/api")

	authAPI := api.Group("/auth")

	authAPI.GET("/get", userCtrl.GetUserData)

	return svc, nil
}
