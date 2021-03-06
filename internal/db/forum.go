package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rinatkh/db_forum/internal/model/core"

	"github.com/jackc/pgx/v5"
)

type ForumRepository interface {
	CreateForum(ctx context.Context, forum *core.Forum) error
	GetForum(ctx context.Context, slug string) (*core.Forum, error)
	GetForumUsers(ctx context.Context, slug string, limit int64, since string, desc bool) ([]*core.User, error)
	GetForumThreads(ctx context.Context, slug string, limit int64, since string, desc bool) ([]*core.Thread, error)
}

type forumRepositoryImpl struct {
	dbConn *pgxpool.Pool
}

func (repo *forumRepositoryImpl) GetForum(ctx context.Context, slug string) (*core.Forum, error) {
	forum := &core.Forum{}
	err := repo.dbConn.QueryRow(ctx,
		`SELECT title, "user", slug, posts, threads FROM Forums WHERE slug = $1;`, slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	return forum, err
}

func (repo *forumRepositoryImpl) GetForumUsers(ctx context.Context, slug string, limit int64, since string, desc bool) ([]*core.User, error) {
	query := constructGetForumUsersQuery(limit, since, desc)
	rows, err := repo.dbConn.Query(ctx, query, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*core.User, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		u := &core.User{}
		if err := rows.Scan(&u.Nickname, &u.Fullname, &u.About, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (repo *forumRepositoryImpl) CreateForum(ctx context.Context, forum *core.Forum) error {
	_, err := repo.dbConn.Exec(ctx,
		`INSERT INTO Forums (title, "user", slug) VALUES ($1, $2, $3);`,
		&forum.Title, &forum.User, &forum.Slug)
	return err
}

func NewForumRepository(dbConn *pgxpool.Pool) *forumRepositoryImpl {
	return &forumRepositoryImpl{dbConn: dbConn}
}

func constructGetForumUsersQuery(limit int64, since string, desc bool) string {
	query := "SELECT u.nickname, u.fullname, u.about, u.email from ForumUsers u where u.forum = $1 "

	if len(since) > 0 {
		if desc {
			query += fmt.Sprintf("and u.nickname < '%s' ", since)
		} else {
			query += fmt.Sprintf("and u.nickname > '%s' ", since)
		}
	}

	query += "ORDER BY u.nickname "
	if desc {
		query += "DESC "
	}
	query += fmt.Sprintf("LIMIT %d ", limit)

	return query
}

func (repo *forumRepositoryImpl) GetForumThreads(ctx context.Context, slug string, limit int64, since string, desc bool) ([]*core.Thread, error) {
	var rows pgx.Rows
	var err error

	query := "SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created FROM Threads as t LEFT JOIN Forums f ON t.forum = f.slug WHERE f.slug = $1 "

	queryOrderBy := "ORDER BY t.created "
	if desc {
		queryOrderBy += "DESC "
	}
	if limit > 0 {
		queryOrderBy += fmt.Sprintf("LIMIT %d ", limit)
	}

	if since != "" {
		querySince := "AND t.created >= $2 "
		if since != "" && desc {
			querySince = "AND t.created <= $2 "
		} else if since != "" && !desc {
			querySince = "AND t.created >= $2 "
		}

		query += querySince + queryOrderBy
		rows, err = repo.dbConn.Query(ctx, query, slug, since)
	} else {
		query += queryOrderBy
		rows, err = repo.dbConn.Query(ctx, query, slug)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := make([]*core.Thread, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		t := &core.Thread{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created); err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}

	return threads, nil
}
