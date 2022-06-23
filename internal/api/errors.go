package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/rinatkh/db_forum/internal/constants"
	"github.com/rinatkh/db_forum/internal/model/dto"
)

func (svc *APIService) httpErrorHandler(err error, c echo.Context) {
	e := err
	msg := err.Error()
	for e != nil {
		if ce, ok := e.(*constants.CodedError); ok {
			code := ce.Code()
			if !svc.debug {
				if code == http.StatusInternalServerError {
					msg = "internal server error"
				} else {
					msg = e.Error()
				}
			}

			_ = c.JSON(code, dto.ErrorResponse{
				Message: msg,
				Code:    code,
			})

			return
		} else {
			e = errors.Unwrap(e)
		}
	}

	if !svc.debug {
		msg = "internal server error"
	}

	_ = c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
		Message: msg,
		Code:    http.StatusInternalServerError,
	})
}
