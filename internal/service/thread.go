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

type ThreadService interface {
	CreateThread(ctx context.Context, request *dto.CreateThreadRequest) (*dto.Response, error)
	CountVote(ctx context.Context, slugOrID string, request *dto.UpdateVoteRequest) (*dto.Response, error)
	GetThread(ctx context.Context, slugOrID string) (*dto.Response, error)
	EditThread(ctx context.Context, slugOrID string, request *dto.UpdateThreadRequest) (*dto.Response, error)
}

type threadServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *threadServiceImpl) CreateThread(ctx context.Context, request *dto.CreateThreadRequest) (*dto.Response, error) {
	user, err := svc.db.UserRepository.GetUserByNickname(ctx, request.Author)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.Author)}, Code: http.StatusNotFound}, nil
		}
	}
	request.Author = user.Nickname

	if forum, err := svc.db.ForumRepository.GetForum(ctx, request.Forum); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", request.Forum)}, Code: http.StatusNotFound}, nil
		}
	} else {
		request.Forum = forum.Slug
	}

	if request.Slug != "" {
		if thread, err := svc.db.ThreadRepository.GetThreadBySlug(ctx, request.Slug); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}
		} else {
			return &dto.Response{Data: thread, Code: http.StatusConflict}, nil
		}
	}

	reqThread := &core.Thread{Forum: request.Forum, Title: request.Title, Author: request.Author, Message: request.Message, Slug: request.Slug, Created: request.Created}
	thread, err := svc.db.ThreadRepository.CreateThread(ctx, reqThread)
	if err != nil {
		return nil, err
	}

	return &dto.Response{Data: thread, Code: http.StatusCreated}, nil
}

func (svc *threadServiceImpl) CountVote(ctx context.Context, slugOrID string, request *dto.UpdateVoteRequest) (*dto.Response, error) {
	var id int
	var err error
	id, err = strconv.Atoi(slugOrID)

	var thread *core.Thread
	if err != nil {
		if thread, err = svc.db.ThreadRepository.GetThreadBySlug(ctx, slugOrID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", slugOrID)}, Code: http.StatusNotFound}, nil
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

	user, err := svc.db.UserRepository.GetUserByNickname(ctx, request.Nickname)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.Nickname)}, Code: http.StatusNotFound}, nil
		}
	}
	request.Nickname = user.Nickname

	exists, err := svc.db.VotesRepository.VoteExists(ctx, request.Nickname, thread.ID)
	if err != nil {
		return nil, err
	}

	if exists {
		if ok, err := svc.db.VotesRepository.UpdateVote(ctx, thread.ID, request.Nickname, request.Voice); err != nil {
			return nil, err
		} else if ok {
			thread.Votes += request.Voice * 2
		}
	} else {
		newVote := &core.Vote{
			Nickname: request.Nickname,
			ThreadID: thread.ID,
			Voice:    request.Voice,
		}

		if err := svc.db.VotesRepository.CreateVote(ctx, newVote); err != nil {
			return nil, err
		}

		thread.Votes += request.Voice
	}

	return &dto.Response{Data: thread, Code: http.StatusOK}, nil
}

func (svc *threadServiceImpl) GetThread(ctx context.Context, slugOrID string) (*dto.Response, error) {
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		if thread, err := svc.db.ThreadRepository.GetThreadBySlug(ctx, slugOrID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", slugOrID)}, Code: http.StatusNotFound}, nil
			}
			return nil, err
		} else {
			return &dto.Response{Data: thread, Code: http.StatusOK}, nil
		}
	}

	thread, err := svc.db.ThreadRepository.GetThreadByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by id: %d", id)}, Code: http.StatusNotFound}, nil
		}
	}

	return &dto.Response{Data: thread, Code: http.StatusOK}, nil
}

func (svc *threadServiceImpl) EditThread(ctx context.Context, slugOrID string, request *dto.UpdateThreadRequest) (*dto.Response, error) {
	var thread *core.Thread
	var err error

	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		if thread, err = svc.db.ThreadRepository.GetThreadBySlug(ctx, slugOrID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", slugOrID)}, Code: http.StatusNotFound}, nil
			}
			return nil, err
		} else {
			id = int(thread.ID)
		}
	}

	if thread, err = svc.db.ThreadRepository.GetThreadByID(ctx, int64(id)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by id: %d", id)}, Code: http.StatusNotFound}, nil
		}
	}

	if len(request.Title) == 0 {
		request.Title = thread.Title
	}

	if len(request.Message) == 0 {
		request.Message = thread.Message
	}

	thread, err = svc.db.ThreadRepository.UpdateThreadByID(ctx, int64(id), request.Title, request.Message)
	return &dto.Response{Data: thread, Code: http.StatusOK}, err
}

func NewThreadService(log *logrus.Entry, db *db.Repository) ThreadService {
	return &threadServiceImpl{log: log, db: db}
}
