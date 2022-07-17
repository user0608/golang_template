package handlers

import (
	"goqrs/handlers/login"
	"goqrs/repositories"
	"goqrs/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func StartRoutes(e *echo.Echo) {
	e.GET("/health", health())
	e.GET("/healthy", health())
	loginService := services.NewLoginService(repositories.NewLoginRepository())
	e.POST("/login", login.Handler(loginService))
}
func health() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"message": "success!"})
	}
}
