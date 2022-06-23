package service

import (
	"github.com/rinatkh/db_forum/internal/db"

	"github.com/sirupsen/logrus"
)

type UserService interface {
}

type userServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func NewUserService(log *logrus.Entry, db *db.Repository) UserService {
	return &userServiceImpl{log: log, db: db}
}
