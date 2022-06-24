package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	UserRepository   UserRepository
	ForumRepository  ForumRepository
	ThreadRepository ThreadRepository
	VotesRepository  VotesRepository
}

func NewRepository(dbConn *pgxpool.Pool) (*Repository, error) {
	repository := &Repository{}

	repository.UserRepository = NewUserRepository(dbConn)
	repository.ForumRepository = NewForumRepository(dbConn)
	repository.ThreadRepository = NewThreadRepository(dbConn)
	repository.VotesRepository = NewVotesRepository(dbConn)
	return repository, nil
}
