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
)

type ForumService interface {
	CreateForum(ctx context.Context, request *dto.CreateForumRequest) (*dto.Response, error)
	GetForum(ctx context.Context, request *dto.GetForumRequest) (*dto.Response, error)
	GetForumThreads(ctx context.Context, request *dto.GetForumThreadsRequest) (*dto.Response, error)
	GetForumUsers(ctx context.Context, request *dto.GetForumUsersRequest) (*dto.Response, error)
}

type forumServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *forumServiceImpl) GetForum(ctx context.Context, request *dto.GetForumRequest) (*dto.Response, error) {
	forum, err := svc.db.ForumRepository.GetForum(ctx, request.Slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find forum with slug: %s", request.Slug)}, Code: http.StatusNotFound}, nil
		}
	}
	return &dto.Response{Data: forum, Code: http.StatusOK}, nil
}

func (svc *forumServiceImpl) GetForumThreads(ctx context.Context, request *dto.GetForumThreadsRequest) (*dto.Response, error) {
	if forum, err := svc.db.ForumRepository.GetForum(ctx, request.Slug); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find forum with slug: %s", request.Slug)}, Code: http.StatusNotFound}, nil
		}
	} else {
		request.Slug = forum.Slug
	}

	threads, err := svc.db.ForumRepository.GetForumThreads(ctx, request.Slug, request.Limit, request.Since, request.Desc)
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: threads, Code: http.StatusOK}, nil
}

func (svc *forumServiceImpl) CreateForum(ctx context.Context, request *dto.CreateForumRequest) (*dto.Response, error) {
	if forum, err := svc.db.ForumRepository.GetForum(ctx, request.Slug); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	} else {
		return &dto.Response{Data: forum, Code: http.StatusConflict}, nil
	}

	user, err := svc.db.UserRepository.GetUserByNickname(ctx, request.User)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.User)}, Code: http.StatusNotFound}, nil
		}
	}
	request.User = user.Nickname

	if err := svc.db.ForumRepository.CreateForum(ctx, &core.Forum{Title: request.Title, User: request.User, Slug: request.Slug}); err != nil {
		return nil, err
	}

	forum, err := svc.db.ForumRepository.GetForum(ctx, request.Slug)
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: forum, Code: http.StatusCreated}, nil
}

func (svc *forumServiceImpl) GetForumUsers(ctx context.Context, request *dto.GetForumUsersRequest) (*dto.Response, error) {
	if forum, err := svc.db.ForumRepository.GetForum(ctx, request.Slug); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find forum with slug: %s", request.Slug)}, Code: http.StatusNotFound}, nil
		}
	} else {
		request.Slug = forum.Slug
	}

	threads, err := svc.db.ForumRepository.GetForumUsers(ctx, request.Slug, request.Limit, request.Since, request.Desc)
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: threads, Code: http.StatusOK}, nil
}

func NewForumService(log *logrus.Entry, db *db.Repository) ForumService {
	return &forumServiceImpl{log: log, db: db}
}
