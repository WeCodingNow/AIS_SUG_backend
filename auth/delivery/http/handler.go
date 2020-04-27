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

func createUser(h *Handler, c echo.Context, role int) error {
	inp := new(signInput)

	if err := c.Bind(inp); err != nil {
		return err
	}

	return h.useCase.CreateUser(c.Request().Context(), inp.Username, inp.Password, role)
}

func (h *Handler) CreateStudent(c echo.Context) error {
	err := createUser(h, c, auth.StudentClass)

	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "Created student")
}

func (h *Handler) CreateHeadman(c echo.Context) error {
	err := createUser(h, c, auth.HeadmanClass)

	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "Created heamdan")
}

type signedInJson struct {
	Username string `json:"username"`
	Token    string `json:"token"`
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

	return c.JSON(http.StatusOK, signedInJson{
		Username: inp.Username,
		Token:    token,
	})
}

func (h *Handler) CheckToken(c echo.Context) error {
	inp := new(signedInJson)

	if err := c.Bind(inp); err != nil {
		return err
	}

	if user, err := h.useCase.ParseToken(c.Request().Context(), inp.Token); err != nil {
		return err
	} else {
		return c.String(http.StatusOK, fmt.Sprintf(
			"token %s gives username: %s",
			inp.Token, user.Username,
		))
	}
}
