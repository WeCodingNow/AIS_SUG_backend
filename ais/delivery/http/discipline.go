package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

func (h *Handler) GetDiscipline(c echo.Context) error {
	disciplineIDParam := c.Param("id")
	disciplineID, err := strconv.Atoi(disciplineIDParam)

	if err != nil {
		return err
	}

	discipline, err := h.useCase.GetDiscipline(c.Request().Context(), disciplineID)

	if err != nil {
		if err == ais.ErrDisciplineNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, models.ToJSONDiscipline(discipline, nil))
}

func (h *Handler) GetAllDisciplines(c echo.Context) error {
	disciplines, err := h.useCase.GetAllDisciplines(c.Request().Context())

	if err != nil {
		return err
	}

	jsonDisciplines := make([]models.JSONMap, 0, len(disciplines))
	for _, discipline := range disciplines {
		jsonDisciplines = append(jsonDisciplines, models.ToJSONDiscipline(discipline, nil))
	}

	return c.JSON(http.StatusOK, jsonDisciplines)
}
