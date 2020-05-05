package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
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

	return c.JSON(http.StatusOK, models.ToJSONStudent(student, nil))
}

func (h *Handler) GetAllStudents(c echo.Context) error {
	students, err := h.useCase.GetAllStudents(c.Request().Context())

	if err != nil {
		return err
	}

	jsonStudents := make([]models.JSONMap, 0, len(students))
	for _, student := range students {
		jsonStudents = append(jsonStudents, models.ToJSONStudent(student, nil))
	}

	return c.JSON(http.StatusOK, jsonStudents)
}
