package api

import (
	"bytes"
	"fmt"
	"github.com/labstack/gommon/log"

	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type validatorImpl struct {
	validator *validator.Validate
}

func (v *validatorImpl) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func NewValidator() echo.Validator {
	return &validatorImpl{validator: validator.New()}
}

type binderImpl struct{}

func (b *binderImpl) Bind(i interface{}, ctx echo.Context) error {
	db := new(echo.DefaultBinder)

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(ctx.Request().Body)
	if err != nil {
		log.Errorf("Unmarshal error: %s", err)
		return err
	}
	err = sonic.Unmarshal(buf.Bytes(), i)
	if err != nil {
		log.Errorf("Unmarshal error: %s", err)
	}
	if err := db.BindQueryParams(ctx, i); err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := db.BindPathParams(ctx, i); err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := db.BindHeaders(ctx, i); err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := ctx.Validate(i); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func NewBinder() echo.Binder {
	return &binderImpl{}
}
