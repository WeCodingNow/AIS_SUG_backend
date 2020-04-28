package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

type Discipline struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Hours int    `json:"hours"`
}

func toJsonDiscipline(discipline *models.Discipline) *Discipline {
	return &Discipline{
		discipline.ID,
		discipline.Name,
		discipline.Hours,
	}
}

func (h *Handler) GetDiscipline(c echo.Context) error {
	disciplineIDParam := c.Param("id")
	disciplineID, err := strconv.Atoi(disciplineIDParam)

	if err != nil {
		return err
	}

	disciplineModel, err := h.useCase.GetDiscipline(c.Request().Context(), disciplineID)

	if err != nil {
		if err == ais.ErrDisciplineNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, toJsonDiscipline(disciplineModel))
}

type manyDisciplinesOutput struct {
	Disciplines []*Discipline `json:"disciplines"`
}

func (h *Handler) GetAllDisciplines(c echo.Context) error {
	disciplineModels, err := h.useCase.GetAllDisciplines(c.Request().Context())

	if err != nil {
		return err
	}

	disciplines := make([]*Discipline, 0, len(disciplineModels))
	for _, disciplineModel := range disciplineModels {
		disciplines = append(disciplines, toJsonDiscipline(disciplineModel))
	}

	return c.JSON(http.StatusOK, manyDisciplinesOutput{Disciplines: disciplines})
}
