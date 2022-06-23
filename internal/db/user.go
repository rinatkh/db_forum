package db

import "github.com/jackc/pgx/v5/pgxpool"

type UserRepository interface {
}

type userRepositoryImpl struct {
	dbConn *pgxpool.Pool
}

func NewUserRepository(dbConn *pgxpool.Pool) *userRepositoryImpl {
	return &userRepositoryImpl{dbConn: dbConn}
}
