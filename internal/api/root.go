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

	repository, err := db.NewRepository(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	svc.router.Validator = NewValidator()
	svc.router.Binder = NewBinder()

	registry := service.NewRegistry(log, repository)
	userCtrl := controllers.NewUserController(log, registry)
	forumCtrl := controllers.NewForumController(log, registry)
	//threadCtrl := controllers.NewThreadController(log, registry)
	//postCtrl := controllers.NewPostController(log, registry)
	//serviceCtrl := controllers.NewServiceController(log, registry)

	api := svc.router.Group("/api")

	api.POST("/forum/create", forumCtrl.CreateForum)
	api.GET("/forum/:slug/details", forumCtrl.GetForum)
	//api.POST("/forum/:slug/create", threadCtrl.CreateForumThread)
	api.GET("/forum/:slug/users", forumCtrl.GetForumUsers)
	api.GET("/forum/:slug/threads", forumCtrl.GetForumThreads)

	//api.GET("/post/:id/details", postCtrl.GetPostDetails)
	//api.POST("/post/:id/details", postCtrl.UpdatePost)
	//
	//api.POST("/service/clear", serviceCtrl.Delete)
	//api.GET("/service/status", serviceCtrl.Status)
	//
	//api.POST("/thread/:slug_or_id/create", threadCtrl.CreatePosts)
	//api.GET("/thread/:slug_or_id/details", threadCtrl.GetForumThreadDetails)
	//api.POST("/thread/:slug_or_id/details", threadCtrl.EditForumThread)
	//api.GET("/thread/:slug_or_id/posts", threadCtrl.GetPosts)
	//api.POST("/thread/:slug_or_id/vote", threadCtrl.CountVote)

	api.POST("/user/:nickname/create", userCtrl.CreateUser)
	api.GET("/user/:nickname/profile", userCtrl.GetUserProfile)
	api.POST("/user/:nickname/profile", userCtrl.EditUserProfile)

	return svc, nil
}
