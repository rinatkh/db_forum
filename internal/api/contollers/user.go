package controllers

import (
	"github.com/rinatkh/db_forum/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *UserController) GetUserData(ctx echo.Context) error {
	//request := new(dto.GetUserRequest)
	//if err := ctx.Bind(request); err != nil {
	//	c.log.Errorf("Bind error: %s", err)
	//	return err
	//}

	//if len(request.UserID) == 0 {
	//	request.UserID = ctx.Request().Header.Get(constants.HeaderKeyUserID)
	//}
	//
	//response, err := c.registry.UserService.GetUserData(context.Background(), request.UserID)
	//if err != nil {
	//	return err
	//}
	response := 1
	return ctx.JSON(http.StatusOK, response)
}

func NewUserController(log *logrus.Entry, registry *service.Registry) *UserController {
	return &UserController{log: log, registry: registry}
}
