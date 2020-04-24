package http

import (
	"github.com/WeCodingNow/AIS_SUG_backend/auth"
	"github.com/labstack/echo"
)

func RegisterHTTPEndpoints(e *echo.Echo, uc auth.UseCase) {
	h := NewHandler(uc)

	e.POST("/sign-up", h.SignUp)
	e.POST("/sign-in", h.SignIn)
	e.POST("/check-token", h.CheckToken, makeAuthMiddleware(uc))
}
