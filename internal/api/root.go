package api

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	controllers "github.com/rinatkh/db_forum/internal/api/contollers"
	"github.com/rinatkh/db_forum/internal/db"
	"github.com/rinatkh/db_forum/internal/service"
	"github.com/sirupsen/logrus"
)

type APIService struct {
	log    *logrus.Entry
	router *echo.Echo
}

func (svc *APIService) Serve() {
	svc.log.Info("Starting HTTP server")
	listenAddr := "0.0.0.0:5000"
	svc.log.Fatal(svc.router.Start(listenAddr))
}

func (svc *APIService) Shutdown(ctx context.Context) error {
	if err := svc.router.Shutdown(ctx); err != nil {
		svc.log.Fatal(err)
	}
	return nil
}

func NewAPIService(log *logrus.Entry, dbConn *pgxpool.Pool) (*APIService, error) {
	svc := &APIService{
		log:    log,
		router: echo.New(),
	}

	repository, err := db.NewRepository(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	svc.router.Validator = NewValidator()
	svc.router.Binder = NewBinder()
	svc.router.Use(svc.LoggingMiddleware())

	registry := service.NewRegistry(log, repository)
	userCtrl := controllers.NewUserController(log, registry)
	forumCtrl := controllers.NewForumController(log, registry)
	threadCtrl := controllers.NewThreadController(log, registry)
	postCtrl := controllers.NewPostController(log, registry)
	serviceCtrl := controllers.NewServiceController(log, repository)

	api := svc.router.Group("/api")

	api.POST("/forum/create", forumCtrl.CreateForum)
	api.GET("/forum/:slug/details", forumCtrl.GetForum)
	api.POST("/forum/:slug/create", threadCtrl.CreateThread)
	api.GET("/forum/:slug/users", forumCtrl.GetForumUsers)
	api.GET("/forum/:slug/threads", forumCtrl.GetForumThreads)

	api.GET("/post/:id/details", postCtrl.GetPostDetails)
	api.POST("/post/:id/details", postCtrl.UpdatePost)

	api.POST("/service/clear", serviceCtrl.Clear)
	api.GET("/service/status", serviceCtrl.Status)

	api.POST("/thread/:slug_or_id/create", postCtrl.CreatePosts)
	api.GET("/thread/:slug_or_id/details", threadCtrl.GetThread)
	api.POST("/thread/:slug_or_id/details", threadCtrl.EditThread)
	api.GET("/thread/:slug_or_id/posts", postCtrl.GetPosts)
	api.POST("/thread/:slug_or_id/vote", threadCtrl.CountVote)

	api.POST("/user/:nickname/create", userCtrl.CreateUser)
	api.GET("/user/:nickname/profile", userCtrl.GetUserProfile)
	api.POST("/user/:nickname/profile", userCtrl.EditUserProfile)

	return svc, nil
}
