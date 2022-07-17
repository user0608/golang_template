package login

import (
	"goqrs/security"
	"goqrs/services"
	"goqrs/utils"

	"github.com/ksaucedo002/answer"
	"github.com/labstack/echo/v4"
)

func Handler(s services.LoginService) echo.HandlerFunc {
	type jsonLogin struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(c echo.Context) error {
		var jsonlogin jsonLogin
		if err := utils.JSON(c, &jsonlogin); err != nil {
			return nil
		}
		account, err := s.Login(c.Request().Context(), jsonlogin.Username, jsonlogin.Password)
		if err != nil {
			return answer.ErrorResponse(c, err)
		}
		tokenString, err := security.GenToken(security.JWTValues{Username: account.Username})
		if err != nil {
			return answer.ErrorResponse(c, err)
		}
		return answer.OK(c, echo.Map{"token": tokenString, "account": account})
	}
}
