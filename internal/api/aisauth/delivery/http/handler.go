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

func (h *Handler) GetInfo(c echo.Context) error {
	userID := c.Get(auth.ContextUserID).(int)

	studentID, err := h.useCase.GetUserStudentID(c.Request().Context(), userID)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":    userID,
		"student_id": studentID,
	})
}

type studentUserInput struct {
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	FirstName   string  `json:"first_name"`
	SecondName  string  `json:"second_name"`
	ThirdName   *string `json:"third_name"`
	GroupID     int     `json:"group_id"`
	ResidenceID int     `json:"residence_id"`
}

func (h *Handler) AssignOrCreateStudentWithCreds(c echo.Context) error {
	inp := new(studentUserInput)

	if err := c.Bind(inp); err != nil {
		return err
	}

	err := h.useCase.AssignOrCreateStudentWithCreds(
		c.Request().Context(),
		&models.User{
			Username: inp.Username,
			Password: inp.Password},
		&models.Student{
			Name:       inp.FirstName,
			SecondName: inp.SecondName,
			ThirdName:  inp.ThirdName,
			Residence:  &models.Residence{ID: inp.ResidenceID},
			Group:      &models.Group{ID: inp.GroupID}},
	)

	if err != nil {
		if err == aisauth.ErrDuplicateUser {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return err
	}

	return c.String(http.StatusOK, "succesfully registered student")
}

func (h *Handler) GetStudentsWithUsers(c echo.Context) error {
	studentsWithUsers, err := h.useCase.GetStudentsWithUsers(c.Request().Context())

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, studentsWithUsers)
}

func (h *Handler) GetRoles(c echo.Context) error {
	roles, err := h.useCase.GetRoles(c.Request().Context())

	if err != nil {
		return err
	}

	jsonRoles := make([]map[string]interface{}, len(roles))

	for i, role := range roles {
		jsonRoles[i] = map[string]interface{}{
			"def": role.Def,
			"id":  role.ID,
		}
	}

	return c.JSON(http.StatusOK, jsonRoles)
}

type PromoteInput struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
}

func (h *Handler) PromoteUser(c echo.Context) error {
	var inp PromoteInput
	err := c.Bind(&inp)

	if err != nil {
		return err
	}

	err = h.useCase.PromoteUser(c.Request().Context(), inp.UserID, inp.RoleID)

	if err != nil {
		return err
	}

	return c.String(http.StatusOK, fmt.Sprintf("succesfully promoted user with %d to role id %d", inp.UserID, inp.RoleID))
}
