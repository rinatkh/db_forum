package controllers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/rinatkh/db_forum/internal/model/dto"
	"github.com/rinatkh/db_forum/internal/service"
	"github.com/sirupsen/logrus"
)

type ThreadController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *ThreadController) CreateThread(ctx echo.Context) error {
	request := new(dto.CreateThreadRequest)

	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	request.Forum = ctx.Param("slug")

	response, err := c.registry.ThreadService.CreateThread(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Data)
}

func (c *ThreadController) CountVote(ctx echo.Context) error {
	request := &dto.EditVoteRequest{}
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	soi := ctx.Param("slug_or_id")
	response, err := c.registry.ThreadService.CountVote(context.Background(), soi, request)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Data)
}

func (c *ThreadController) GetThread(ctx echo.Context) error {
	soi := ctx.Param("slug_or_id")

	response, err := c.registry.ThreadService.GetThread(context.Background(), soi)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Data)
}

func (c *ThreadController) EditThread(ctx echo.Context) error {
	request := &dto.EditThreadRequest{}
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	soi := ctx.Param("slug_or_id")

	response, err := c.registry.ThreadService.EditThread(context.Background(), soi, request)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Data)
}

func NewThreadController(log *logrus.Entry, registry *service.Registry) *ThreadController {
	return &ThreadController{log: log, registry: registry}
}
