package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/rinatkh/db_forum/internal/model/dto"
	"github.com/rinatkh/db_forum/internal/service"
	"github.com/sirupsen/logrus"
	"strconv"
)

type PostController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *PostController) CreatePosts(ctx echo.Context) error {
	var request []*dto.Post
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(ctx.Request().Body)
	if err != nil {
		log.Errorf("Unmarshal error: %s", err)
		return err
	}
	err = json.Unmarshal(buf.Bytes(), &request)

	slugOrID := ctx.Param("slug_or_id")
	response, err := c.registry.PostsService.CreatePosts(context.Background(), slugOrID, request)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Data)
}

func (c *PostController) GetPosts(ctx echo.Context) error {
	slugOrID := ctx.Param("slug_or_id")
	sort := ctx.QueryParam("sort")
	if sort == "" {
		sort = "flat"
	}

	sinceParam := ctx.QueryParam("since")
	if sinceParam == "" {
		sinceParam = "-1"
	}

	limitParam := ctx.QueryParam("limit")
	if sinceParam == "" {
		sinceParam = "100"
	}

	since, _ := strconv.ParseInt(sinceParam, 10, 64)
	desc, _ := strconv.ParseBool(ctx.QueryParam("desc"))
	limit, _ := strconv.ParseInt(limitParam, 10, 64)

	response, err := c.registry.PostsService.GetPosts(context.Background(), slugOrID, sort, since, desc, limit)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Data)
}

func (c *PostController) GetPostDetails(ctx echo.Context) error {
	request := new(dto.GetPostDetailsRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	request.ID = id

	response, err := c.registry.PostsService.GetPostDetails(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Data)
}

func (c *PostController) UpdatePost(ctx echo.Context) error {

	request := new(dto.EditPostRequest)
	if err := ctx.Bind(request); err != nil {
		c.log.Errorf("Bind error: %s", err)
		return err
	}
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	request.ID = id

	response, err := c.registry.PostsService.UpdatePost(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Data)
}

func NewPostController(log *logrus.Entry, registry *service.Registry) *PostController {
	return &PostController{log: log, registry: registry}
}
