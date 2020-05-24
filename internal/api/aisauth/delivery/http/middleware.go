package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/aisauth"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/labstack/echo"

	authhttp "github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth/delivery/http"
)

const ContextRoleID = "role-id"

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

func MakeRBACMiddleware(a aisauth.UseCase, rolesWithAccess []int) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get(authhttp.ContextUserID).(int)
			role, err := a.GetUserRole(c.Request().Context(), userID)

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			gotAccess := false

			for _, r := range rolesWithAccess {
				if role.ID == r {
					gotAccess = true
				}
			}

			if !gotAccess {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					fmt.Sprintf(
						"role %s doesn't have access to %s",
						role.Def, c.Request().URL.Path))
			}

			return next(c)
		}
	}
}
