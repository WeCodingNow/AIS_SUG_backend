package http

import (
	"fmt"
	"net/http"

	"github.com/WeCodingNow/AIS_SUG_backend/auth"
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

func (h *Handler) SignUp(ctx echo.Context) error {
	inp := new(signInput)

	if err := ctx.Bind(inp); err != nil {
		return err
	}

	if err := h.useCase.SignUp(ctx.Request().Context(), inp.Username, inp.Password); err != nil {
		return err
	}

	return ctx.String(http.StatusOK, "Signed up")
}

type signedInJson struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func (h *Handler) SignIn(ctx echo.Context) error {
	inp := new(signInput)

	if err := ctx.Bind(inp); err != nil {
		return err
	}

	if token, err := h.useCase.SignIn(ctx.Request().Context(), inp.Username, inp.Password); err != nil {
		return err
	} else {
		return ctx.JSON(http.StatusOK, signedInJson{
			Username: inp.Username,
			Token:    token,
		})
	}
}

func (h *Handler) CheckToken(ctx echo.Context) error {
	inp := new(signedInJson)

	if err := ctx.Bind(inp); err != nil {
		return err
	}

	if user, err := h.useCase.ParseToken(ctx.Request().Context(), inp.Token); err != nil {
		return err
	} else {
		return ctx.String(http.StatusOK, fmt.Sprintf(
			"token %s gives username: %s",
			inp.Token, user.Username,
		))
	}
}
