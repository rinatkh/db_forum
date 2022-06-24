package controllers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/rinatkh/db_forum/internal/db"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ServiceController struct {
	log *logrus.Entry
	db  *db.Repository
}

func (c *ServiceController) Status(ctx echo.Context) error {
	response, err := c.db.ServiceRepository.Status(context.Background())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, response)
}

func (c *ServiceController) Clear(ctx echo.Context) error {
	err := c.db.ServiceRepository.Delete(context.Background())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, nil)
}

func NewServiceController(log *logrus.Entry, db *db.Repository) *ServiceController {
	return &ServiceController{log: log, db: db}
}
