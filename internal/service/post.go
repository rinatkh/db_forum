package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/rinatkh/db_forum/internal/db"
	"github.com/rinatkh/db_forum/internal/model/core"
	"github.com/rinatkh/db_forum/internal/model/dto"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type PostsService interface {
	CreatePosts(ctx context.Context, soi string, posts []*dto.Post) (*dto.Response, error)
	GetPosts(ctx context.Context, soi string, sort string, since int64, desc bool, limit int64) (*dto.Response, error)
	GetPostDetails(ctx context.Context, request *dto.GetPostDetailsRequest) (*dto.Response, error)
	EditPost(ctx context.Context, request *dto.EditPostRequest) (*dto.Response, error)
}

type postsServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *postsServiceImpl) EditPost(ctx context.Context, request *dto.EditPostRequest) (*dto.Response, error) {
	post, err := svc.db.PostsRepository.GetPostByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find post by id: %d", request.ID)}, Code: http.StatusNotFound}, nil
		}
		return nil, err
	}

	if len(request.Message) == 0 || request.Message == post.Message {
		return &dto.Response{Data: post, Code: http.StatusOK}, nil
	}

	updatedPost, err := svc.db.PostsRepository.EditPost(ctx, request.ID, request.Message)
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: updatedPost, Code: http.StatusOK}, nil
}

func (svc *postsServiceImpl) CreatePosts(ctx context.Context, soi string, posts []*dto.Post) (*dto.Response, error) {
	var id int
	var err error
	id, err = strconv.Atoi(soi)

	var thread *core.Thread
	if err != nil {
		if thread, err = svc.db.ThreadRepository.GetThreadBySlug(ctx, soi); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", soi)}, Code: http.StatusNotFound}, nil
			}
		} else {
			id = int(thread.ID)
		}
	} else {
		if thread, err = svc.db.ThreadRepository.GetThreadByID(ctx, int64(id)); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by id: %d", id)}, Code: http.StatusNotFound}, nil
			}
		}
	}

	if len(posts) == 0 {
		return &dto.Response{Data: []struct{}{}, Code: http.StatusCreated}, nil
	}

	if posts[0].Parent != 0 {
		parentThreadID, err := svc.db.PostsRepository.CheckParentPost(ctx, int(posts[0].Parent))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &dto.Response{Data: dto.ErrorResponse{Message: "Parent post was created in another thread"}, Code: http.StatusConflict}, nil
			}
		}

		if parentThreadID != id {
			return &dto.Response{Data: dto.ErrorResponse{Message: "Parent post was created in another thread"}, Code: http.StatusConflict}, nil
		}
	}

	if _, err := svc.db.UserRepository.GetUserByNickname(ctx, posts[0].Author); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", posts[0].Author)}, Code: http.StatusNotFound}, nil
		}
	}

	insertedPosts, err := svc.db.PostsRepository.CreatePosts(ctx, thread.Forum, int64(id), posts)
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: insertedPosts, Code: http.StatusCreated}, nil
}

func (svc *postsServiceImpl) GetPosts(ctx context.Context, soi string, sort string, since int64, desc bool, limit int64) (*dto.Response, error) {
	id, err := strconv.Atoi(soi)
	if err != nil {
		if thread, err := svc.db.ThreadRepository.GetThreadBySlug(ctx, soi); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", soi)}, Code: http.StatusNotFound}, nil
			}
		} else {
			id = int(thread.ID)
		}
	}

	if _, err := svc.db.ThreadRepository.GetThreadByID(ctx, int64(id)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by id: %d", id)}, Code: http.StatusNotFound}, nil
		}
	}

	var posts []*core.Post
	switch sort {
	case "flat":
		posts, err = svc.db.PostsRepository.GetPostsFlat(ctx, id, since, desc, limit)
	case "tree":
		posts, err = svc.db.PostsRepository.GetPostsTree(ctx, id, since, desc, limit)
	case "parent_tree":
		posts, err = svc.db.PostsRepository.GetPostsParentTree(ctx, id, since, desc, limit)
	default:
		posts, err = svc.db.PostsRepository.GetPostsFlat(ctx, id, since, desc, limit)
	}
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: posts, Code: http.StatusOK}, nil
}

func (svc *postsServiceImpl) GetPostDetails(ctx context.Context, request *dto.GetPostDetailsRequest) (*dto.Response, error) {
	post, err := svc.db.PostsRepository.GetPostByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find post by id: %d", request.ID)}, Code: http.StatusNotFound}, nil
		}
		return nil, err
	}

	postDetails, err := svc.db.PostsRepository.GetPostDetails(ctx, request.ID, request.Related)
	if err != nil {
		return nil, err
	}
	postDetails.Post = post

	return &dto.Response{Data: postDetails, Code: http.StatusOK}, nil
}

func NewPostsService(log *logrus.Entry, db *db.Repository) PostsService {
	return &postsServiceImpl{log: log, db: db}
}
