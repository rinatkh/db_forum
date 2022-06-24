package controllers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/rinatkh/db_forum/internal/model/dto"
	"github.com/rinatkh/db_forum/internal/service"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *UserController) CreateUser(ctx echo.Context) error {
	request := new(dto.CreateUserRequest)

	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	request.Nickname = ctx.Param("nickname")

	response, err := c.registry.UserService.CreateUser(context.Background(), request)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Code, response.Data)
}

func (c *UserController) GetUserProfile(ctx echo.Context) error {
	request := new(dto.GetUserProfileRequest)

	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	request.Nickname = ctx.Param("nickname")

	response, err := c.registry.UserService.GetUserProfile(context.Background(), request)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Code, response.Data)
}

func (c *UserController) EditUserProfile(ctx echo.Context) error {
	request := new(dto.EditUserProfileRequest)

	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	request.Nickname = ctx.Param("nickname")

	response, err := c.registry.UserService.EditUserProfile(context.Background(), request)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Code, response.Data)
}

func NewUserController(log *logrus.Entry, registry *service.Registry) *UserController {
	return &UserController{log: log, registry: registry}
}
