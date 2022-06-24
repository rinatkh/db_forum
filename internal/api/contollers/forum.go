package controllers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/rinatkh/db_forum/internal/model/dto"
	"github.com/rinatkh/db_forum/internal/service"
	"github.com/sirupsen/logrus"
)

type ForumController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *ForumController) CreateForum(ctx echo.Context) error {
	request := new(dto.CreateForumRequest)

	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}

	response, err := c.registry.ForumService.CreateForum(context.Background(), request)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Code, response.Data)
}

func (c *ForumController) GetForum(ctx echo.Context) error {
	request := new(dto.GetForumBySlugRequest)

	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	request.Slug = ctx.Param("slug")

	response, err := c.registry.ForumService.GetForum(context.Background(), request)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Code, response.Data)
}

func (c *ForumController) GetForumThreads(ctx echo.Context) error {
	request := new(dto.GetForumThreadsRequest)

	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	request.Slug = ctx.Param("slug")

	response, err := c.registry.ForumService.GetForumThreads(context.Background(), request)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Code, response.Data)
}

func (c *ForumController) GetForumUsers(ctx echo.Context) error {
	request := new(dto.GetForumUsersRequest)

	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	request.Slug = ctx.Param("slug")

	response, err := c.registry.ForumService.GetForumUsers(context.Background(), request)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Code, response.Data)
}

func NewForumController(log *logrus.Entry, registry *service.Registry) *ForumController {
	return &ForumController{log: log, registry: registry}
}
