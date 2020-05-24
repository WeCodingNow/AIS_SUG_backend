package http

import (
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/aisauth"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth"
	authhttp "github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth/delivery/http"
	"github.com/labstack/echo"
)

func RegisterHTTPEndpoints(e *echo.Echo, uc auth.UseCase, auc aisauth.UseCase) {
	h := NewHandler(auc)

	e.GET("/role", h.GetRole, authhttp.MakeAuthMiddleware(uc))
	e.GET("/check-role", h.CheckRole, authhttp.MakeAuthMiddleware(uc), MakeRoleMiddleware(auc))

	e.GET("/roles", h.GetRoles, authhttp.MakeAuthMiddleware(uc), MakeRBACMiddleware(auc, []int{auth.AdminClass}))
	e.POST("/roles", h.PromoteUser, authhttp.MakeAuthMiddleware(uc), MakeRBACMiddleware(auc, []int{auth.AdminClass}))

	e.GET("/info", h.GetInfo, authhttp.MakeAuthMiddleware(uc))

	e.GET("/student_users", h.GetStudentsWithUsers, authhttp.MakeAuthMiddleware(uc), MakeRBACMiddleware(auc, []int{auth.AdminClass}))

	e.POST("/register", h.AssignOrCreateStudentWithCreds)

}
