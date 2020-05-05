package http

import (
	"log"
	"net/http"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/aisauth"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/labstack/echo"

	authhttp "github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth/delivery/http"
)

// const ContextUserID = "user-id"

const ContextRoleID = "role-id"

// func makeUserMiddleware(a aisauth.UseCase) func(echo.HandlerFunc) echo.HandlerFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			return next(c)
// 		}
// 	}
// }

func MakeRoleMiddleware(a aisauth.UseCase) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get(authhttp.ContextUserID).(int)

			roleID, err := a.GetUserRoleID(c.Request().Context(), &models.User{ID: userID})

			if err != nil {
				log.Printf("role of user %s not found", userID)
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			c.Set(ContextRoleID, roleID)

			return next(c)
		}
	}
}
