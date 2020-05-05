package http

import (
	"fmt"
	"net/http"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/aisauth"
	"github.com/labstack/echo"
)

type Handler struct {
	useCase aisauth.UseCase
}

func NewHandler(useCase aisauth.UseCase) *Handler {
	return &Handler{
		useCase,
	}
}

func (h *Handler) CheckRole(c echo.Context) error {
	roleID := c.Get(ContextRoleID)

	if roleID == nil {
		return fmt.Errorf("no role")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"role_id": roleID,
	})
}
