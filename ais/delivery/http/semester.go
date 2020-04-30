package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/ais/delivery/types"
	"github.com/labstack/echo"
)

func (h *Handler) GetSemester(c echo.Context) error {
	semesterIDParam := c.Param("id")
	semesterID, err := strconv.Atoi(semesterIDParam)

	if err != nil {
		return err
	}

	semester, err := h.useCase.GetSemester(c.Request().Context(), semesterID)

	if err != nil {
		if err == ais.ErrSemesterNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, types.ToJsonSemester(semester))
}

type manySemestersOutput struct {
	Semesters []*types.JSONSemester `json:"semesters"`
}

func (h *Handler) GetAllSemesters(c echo.Context) error {
	semesters, err := h.useCase.GetAllSemesters(c.Request().Context())

	if err != nil {
		return err
	}

	semesterJSONs := make([]*types.JSONSemester, 0, len(semesters))
	for _, semester := range semesters {
		semesterJSONs = append(semesterJSONs, types.ToJsonSemester(semester))
	}

	return c.JSON(http.StatusOK, manySemestersOutput{Semesters: semesterJSONs})
}
