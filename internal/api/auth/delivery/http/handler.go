package http

import (
	"fmt"
	"net/http"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth"
	"github.com/labstack/echo"
)

type Handler struct {
	useCase auth.UseCase
}

func NewHandler(useCase auth.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type signInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signedInJSON struct {
	Token string `json:"token"`
}

func (h *Handler) SignIn(c echo.Context) error {
	inp := new(signInput)

	if err := c.Bind(inp); err != nil {
		return err
	}

	token, err := h.useCase.SignIn(c.Request().Context(), inp.Username, inp.Password)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, signedInJSON{
		Token: token,
	})
}

func (h *Handler) CheckToken(c echo.Context) error {
	token := c.Get(ContextUserToken).(string)

	if user, err := h.useCase.ParseToken(c.Request().Context(), token); err != nil {
		return err
	} else {
		return c.String(http.StatusOK, fmt.Sprintf(
			"token %s gives user id: %d",
			token, user.ID,
		))
	}
}
