package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	UserRepository  UserRepository
	ForumRepository ForumRepository
}

func NewRepository(dbConn *pgxpool.Pool) (*Repository, error) {
	repository := &Repository{}

	repository.UserRepository = NewUserRepository(dbConn)
	repository.ForumRepository = NewForumRepository(dbConn)
	return repository, nil
}
