package http

import (
	"fmt"
	"net/http"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/aisauth"
	auth "github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth/delivery/http"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
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

func (h *Handler) GetRole(c echo.Context) error {
	userID := c.Get(auth.ContextUserID).(int)

	role, err := h.useCase.GetUserRole(c.Request().Context(), userID)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.ToJSONRole(role, nil))
}
