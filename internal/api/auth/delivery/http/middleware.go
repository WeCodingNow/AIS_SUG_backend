package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth"
	"github.com/labstack/echo"
)

const ContextUserToken = "user-token"
const ContextUserID = "user-id"

// const ContextUserID = "user-id"

func MakeAuthMiddleware(a auth.UseCase) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)

			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 {
				log.Print("no 2 parts in auth")
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			if headerParts[0] != "Bearer" {
				log.Print("No bearer in header")
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			token := headerParts[1]
			user, err := a.ParseToken(c.Request().Context(), token)

			c.Set(ContextUserToken, token)
			c.Set(ContextUserID, user.ID)

			if err != nil {
				log.Print("another error: ", err.Error())
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}
