package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rinatkh/db_forum/internal/model/core"
)

type ThreadRepository interface {
	CreateThread(ctx context.Context, thread *core.Thread) (*core.Thread, error)
	GetThreadByID(ctx context.Context, id int64) (*core.Thread, error)
	GetThreadBySlug(ctx context.Context, slug string) (*core.Thread, error)
	UpdateThreadByID(ctx context.Context, id int64, title string, message string) (*core.Thread, error)
}

type threadRepositoryImpl struct {
	dbConn *pgxpool.Pool
}

func (repo *threadRepositoryImpl) CreateThread(ctx context.Context, thread *core.Thread) (*core.Thread, error) {
	t := &core.Thread{}
	err := repo.dbConn.QueryRow(ctx, "INSERT INTO Threads (title, author, forum, message, slug, created) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, title, author, forum, message, votes, slug, created;", thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created).
		Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
	return t, err
}

func (repo *threadRepositoryImpl) GetThreadByID(ctx context.Context, id int64) (*core.Thread, error) {
	t := &core.Thread{}
	err := repo.dbConn.QueryRow(ctx, "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE id = $1;", id).Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
	return t, err
}

func (repo *threadRepositoryImpl) GetThreadBySlug(ctx context.Context, slug string) (*core.Thread, error) {
	t := &core.Thread{}
	err := repo.dbConn.QueryRow(ctx, "SELECT id, title, author, forum, message, votes, slug, created FROM Threads WHERE slug = $1;", slug).Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
	return t, err
}

func (repo *threadRepositoryImpl) UpdateThreadByID(ctx context.Context, id int64, title string, message string) (*core.Thread, error) {
	t := &core.Thread{}
	err := repo.dbConn.QueryRow(ctx, "UPDATE Threads SET title = $2, message = $3 WHERE id = $1 RETURNING id, title, author, forum, message, votes, slug, created;", id, title, message).Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
	return t, err
}

func NewThreadRepository(dbConn *pgxpool.Pool) *threadRepositoryImpl {
	return &threadRepositoryImpl{dbConn: dbConn}
}
