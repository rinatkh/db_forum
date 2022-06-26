package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rinatkh/db_forum/internal/model/core"
)

type ServiceRepository interface {
	Status(ctx context.Context) (*core.ServiceInfo, error)
	Delete(ctx context.Context) error
}

type serviceRepositoryImpl struct {
	dbConn *pgxpool.Pool
}

func (repo *serviceRepositoryImpl) Delete(ctx context.Context) error {
	_, err := repo.dbConn.Exec(ctx,
		"TRUNCATE TABLE Users, Forums, Threads, Posts, ForumUsers, Votes CASCADE;")
	return err
}

func (repo *serviceRepositoryImpl) Status(ctx context.Context) (*core.ServiceInfo, error) {
	res := &core.ServiceInfo{}
	err := repo.dbConn.QueryRow(ctx,
		"SELECT (SELECT count(*) FROM Users) AS user, (SELECT count(*) FROM Forums) AS forum, (SELECT count(*) FROM Threads) AS thread, (SELECT count(*) FROM Posts) AS post;").Scan(&res.User, &res.Forum, &res.Thread, &res.Post)
	return res, err
}

func NewServiceRepository(dbConn *pgxpool.Pool) *serviceRepositoryImpl {
	return &serviceRepositoryImpl{dbConn: dbConn}
}
