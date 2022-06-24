package service

import (
	"github.com/rinatkh/db_forum/internal/db"
	"github.com/sirupsen/logrus"
)

type Registry struct {
	UserService  UserService
	ForumService ForumService
}

func NewRegistry(log *logrus.Entry, repository *db.Repository) *Registry {
	registry := new(Registry)

	registry.UserService = NewUserService(log, repository)
	registry.ForumService = NewForumService(log, repository)
	return registry
}
