package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/ais/delivery/types"
	"github.com/labstack/echo"
)

func (h *Handler) GetStudent(c echo.Context) error {
	studentIDParam := c.Param("id")
	studentID, err := strconv.Atoi(studentIDParam)

	if err != nil {
		return err
	}

	student, err := h.useCase.GetStudent(c.Request().Context(), studentID)

	if err != nil {
		if err == ais.ErrStudentNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, types.ToJsonStudent(student))
}

type manyStudentsOutput struct {
	Students []*types.JSONStudent `json:"students"`
}

func (h *Handler) GetAllStudents(c echo.Context) error {
	students, err := h.useCase.GetAllStudents(c.Request().Context())

	if err != nil {
		return err
	}

	jsonStudents := make([]*types.JSONStudent, 0, len(students))
	for _, student := range students {
		jsonStudents = append(jsonStudents, types.ToJsonStudent(student))
	}

	return c.JSON(http.StatusOK, manyStudentsOutput{Students: jsonStudents})
}
