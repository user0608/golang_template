package utils

import (
	"github.com/ksaucedo002/answer/errores"
	"github.com/labstack/echo/v4"
)

func JSON(c echo.Context, payload any) error {
	if err := (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
		return errores.NewBadRequestf(err, "json document invalido")
	}
	return nil
}
