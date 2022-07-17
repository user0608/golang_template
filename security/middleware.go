package security

import (
	"github.com/ksaucedo002/answer"
	"github.com/labstack/echo/v4"
)

const jwtvalueskey = "jwt-values-context-key"

func GetJWTValues(c echo.Context) (JWTValues, bool) {
	values := c.Get(jwtvalueskey)
	if values == nil {
		return JWTValues{}, false
	}
	v, ok := values.(JWTValues)
	return v, ok
}
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			token = c.QueryParam("token")
		}
		values, err := ValidateToken(token)
		if err != nil {
			return answer.ErrorResponse(c, err)
		}
		c.Set(jwtvalueskey, values)
		return next(c)
	}
}
