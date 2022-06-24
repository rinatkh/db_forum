package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/rinatkh/db_forum/internal/db"
	"github.com/rinatkh/db_forum/internal/model/core"
	"github.com/rinatkh/db_forum/internal/model/dto"
	"net/http"

	"github.com/sirupsen/logrus"
)

type UserService interface {
	CreateUser(ctx context.Context, request *dto.CreateUserRequest) (*dto.Response, error)
	GetUserProfile(ctx context.Context, request *dto.GetUserProfileRequest) (*dto.Response, error)
	EditUserProfile(ctx context.Context, request *dto.EditUserProfileRequest) (*dto.Response, error)
}

type userServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *userServiceImpl) EditUserProfile(ctx context.Context, request *dto.EditUserProfileRequest) (*dto.Response, error) {
	if len(request.Email) > 0 {
		if user, err := svc.db.UserRepository.GetUserByEmail(ctx, request.Email); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}
		} else if user.Nickname != request.Nickname {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("This email is already registered by user: %s", user.Nickname)}, Code: http.StatusConflict}, nil
		}
	}

	user := &core.User{Nickname: request.Nickname, Fullname: request.Fullname, About: request.About, Email: request.Email}
	updatedUser, err := svc.db.UserRepository.EditUser(ctx, user)
	if err != nil {
		return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.Nickname)}, Code: http.StatusNotFound}, nil
	}
	return &dto.Response{Data: updatedUser, Code: http.StatusOK}, nil
}

func (svc *userServiceImpl) CreateUser(ctx context.Context, request *dto.CreateUserRequest) (*dto.Response, error) {
	if users, err := svc.db.UserRepository.GetUsersByEmailOrNickname(ctx, request.Email, request.Nickname); err != nil {
		return nil, err
	} else if len(users) > 0 {
		return &dto.Response{Data: users, Code: http.StatusConflict}, nil
	}

	user := &core.User{Nickname: request.Nickname, Fullname: request.Fullname, About: request.About, Email: request.Email}
	if err := svc.db.UserRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &dto.Response{Data: user, Code: http.StatusCreated}, nil
}

func (svc *userServiceImpl) GetUserProfile(ctx context.Context, request *dto.GetUserProfileRequest) (*dto.Response, error) {
	user, err := svc.db.UserRepository.GetUserByNickname(ctx, request.Nickname)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.Response{Data: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.Nickname)}, Code: http.StatusNotFound}, nil
		}
		return nil, err
	}
	return &dto.Response{Data: user, Code: http.StatusOK}, nil
}
func NewUserService(log *logrus.Entry, db *db.Repository) UserService {
	return &userServiceImpl{log: log, db: db}
}
