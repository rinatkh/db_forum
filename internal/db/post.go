package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
	"github.com/rinatkh/db_forum/internal/model/core"
	"github.com/rinatkh/db_forum/internal/model/dto"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type PostsRepository interface {
	CreatePosts(ctx context.Context, forum string, thread int64, posts []*dto.Post) ([]*core.Post, error)
	CheckParentPost(ctx context.Context, parent int) (int, error)
	GetPostsFlat(ctx context.Context, id int, since int64, desc bool, limit int64) ([]*core.Post, error)
	GetPostsTree(ctx context.Context, id int, since int64, desc bool, limit int64) ([]*core.Post, error)
	GetPostsParentTree(ctx context.Context, id int, since int64, desc bool, limit int64) ([]*core.Post, error)
	GetPostDetails(ctx context.Context, id int64, related string) (dto.PostDetails, error)
	GetPostByID(ctx context.Context, id int64) (*core.Post, error)
	UpdatePost(ctx context.Context, id int64, message string) (*core.Post, error)
}

type postsRepositoryImpl struct {
	dbConn *pgxpool.Pool
}

func (repo *postsRepositoryImpl) CreatePosts(ctx context.Context, forum string, thread int64, posts []*dto.Post) ([]*core.Post, error) {
	query := strings.Builder{}
	query.WriteString("INSERT INTO Posts (parent, author, message, forum, thread, created) VALUES ")

	queryArgs := make([]interface{}, 0, len(posts))
	newPosts := make([]*core.Post, 0, len(posts))
	insertTime := time.Unix(0, time.Now().UnixNano()/1e6*1e6)
	for i, post := range posts {
		p := &core.Post{Parent: post.Parent, Author: post.Author, Message: post.Message, Forum: forum, Thread: thread, Created: insertTime}
		newPosts = append(newPosts, p)
		_, err := fmt.Fprintf(&query, "($%d, $%d, $%d, $%d, $%d, $%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		if err != nil {
			return nil, err
		}
		queryArgs = append(queryArgs, post.Parent, post.Author, post.Message, forum, thread, insertTime)
	}

	qs := query.String()
	qs = qs[:len(qs)-1]
	qs += " RETURNING id;"

	rows, err := repo.dbConn.Query(ctx, qs, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		if err = rows.Scan(&newPosts[i].ID); err != nil {
			return nil, err
		}
	}

	return newPosts, nil
}

func (repo *postsRepositoryImpl) CheckParentPost(ctx context.Context, parent int) (int, error) {
	var threadID int
	err := repo.dbConn.QueryRow(ctx,
		"SELECT thread FROM Posts WHERE id = $1;", parent).Scan(&threadID)
	return threadID, err
}

func (repo *postsRepositoryImpl) GetPostsFlat(ctx context.Context, id int, since int64, desc bool, limit int64) ([]*core.Post, error) {
	query := "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 "

	if since != -1 {
		if desc {
			query += "AND id < $2 "
		} else {
			query += "AND id > $2 "
		}
	}

	if desc {
		query += "ORDER BY created DESC, id DESC "
	} else {
		query += "ORDER BY created ASC, id ASC "
	}

	query += fmt.Sprintf("LIMIT %d ", limit)

	var rows pgx.Rows
	var err error
	if since == -1 {
		rows, err = repo.dbConn.Query(ctx, query, id)
	} else {
		rows, err = repo.dbConn.Query(ctx, query, id, since)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*core.Post, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		post := &core.Post{}
		if err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, err
}

func (repo *postsRepositoryImpl) GetPostsTree(ctx context.Context, id int, since int64, desc bool, limit int64) ([]*core.Post, error) {
	query := "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 "

	if since != -1 {
		if desc {
			query += "and path < "
		} else {
			query += "and path > "
		}
		query += fmt.Sprintf("(SELECT path FROM posts WHERE id = %d) ", since)
	}

	if desc {
		query += "ORDER BY path desc "
	} else {
		query += "ORDER BY path asc, id "
	}

	query += fmt.Sprintf("LIMIT NULLIF(%d, 0) ", limit)

	rows, err := repo.dbConn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*core.Post, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		post := &core.Post{}
		if err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (repo *postsRepositoryImpl) GetPostsParentTree(ctx context.Context, id int, since int64, desc bool, limit int64) ([]*core.Post, error) {
	var rows pgx.Rows
	var err error
	if since == -1 {
		if desc {
			rows, err = repo.dbConn.Query(ctx,
				`SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts
					WHERE path[1] IN (SELECT id FROM Posts WHERE thread = $1 AND parent = 0 ORDER BY id DESC LIMIT $2)
					ORDER BY path[1] DESC, path ASC, id ASC;`,
				id, limit)
		} else {
			rows, err = repo.dbConn.Query(ctx,
				`SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts
					WHERE path[1] IN (SELECT id FROM Posts WHERE thread = $1 AND parent = 0 ORDER BY id ASC LIMIT $2)
					ORDER BY path ASC, id ASC;`,
				id, limit)
		}
	} else {
		if desc {
			rows, err = repo.dbConn.Query(ctx,
				`SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts
					WHERE path[1] IN (SELECT id FROM Posts WHERE thread = $1 AND parent = 0 AND path[1] < (SELECT path[1] FROM posts WHERE id = $2)
					ORDER BY id DESC LIMIT $3) ORDER BY path[1] DESC, path ASC, id ASC;`,
				id, since, limit)
		} else {
			rows, err = repo.dbConn.Query(ctx,
				`SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts
					WHERE path[1] IN (SELECT id FROM Posts WHERE thread = $1 AND parent = 0 AND path[1] >
					(SELECT path[1] FROM Posts WHERE id = $2) ORDER BY id ASC LIMIT $3) 
					ORDER BY path ASC, id ASC;`,
				id, since, limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*core.Post, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		post := &core.Post{}
		if err := rows.Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (repo *postsRepositoryImpl) GetPostDetails(ctx context.Context, id int64, related string) (dto.PostDetails, error) {
	postDetails := dto.PostDetails{}
	for _, arg := range strings.Split(related, ",") {
		switch arg {
		case "user":
			author := &core.User{}
			err := repo.dbConn.QueryRow(ctx,
				"SELECT a.nickname, a.fullname, a.about, a.email FROM Posts JOIN Users a ON a.nickname = Posts.author WHERE posts.id = $1;", id).Scan(&author.Nickname, &author.Fullname, &author.About, &author.Email)
			if err != nil {
				return dto.PostDetails{}, err
			}
			postDetails.Author = author

		case "thread":
			thread := &core.Thread{}
			err := repo.dbConn.QueryRow(ctx,
				"SELECT th.id, th.title, th.author, th.forum, th.message, th.votes, th.slug, th.created FROM Posts JOIN Threads th ON th.id = Posts.thread WHERE posts.id = $1;",
				id).
				Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)

			if err != nil {
				return dto.PostDetails{}, err
			}
			postDetails.Thread = thread

		case "forum":
			forum := &core.Forum{}
			err := repo.dbConn.QueryRow(ctx,
				"SELECT f.title, f.user, f.slug, f.posts, f.threads FROM Posts JOIN Forums f ON f.slug = Posts.forum WHERE Posts.id = $1;", id).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
			if err != nil {
				return dto.PostDetails{}, err
			}
			postDetails.Forum = forum
		}
	}
	log.Infof("%v", postDetails)
	return postDetails, nil
}

func (repo *postsRepositoryImpl) GetPostByID(ctx context.Context, id int64) (*core.Post, error) {
	post := &core.Post{}
	err := repo.dbConn.QueryRow(ctx,
		"SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE id = $1;",
		id).
		Scan(&post.ID, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	return post, err
}

func (repo *postsRepositoryImpl) UpdatePost(ctx context.Context, id int64, message string) (*core.Post, error) {
	post := &core.Post{}
	err := repo.dbConn.QueryRow(ctx,
		"UPDATE Posts SET message = $2, isEdited = true WHERE id = $1 RETURNING id, parent, author, message, isEdited, forum, thread, created;",
		id, message).
		Scan(&post.ID, &post.Parent, &post.Author, &post.Message,
			&post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func NewPostsRepository(dbConn *pgxpool.Pool) *postsRepositoryImpl {
	return &postsRepositoryImpl{dbConn: dbConn}
}
