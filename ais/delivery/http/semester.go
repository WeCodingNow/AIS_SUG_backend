package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

type Semester struct {
	ID        int        `json:"id"`
	Number    int        `json:"number"`
	Beginning time.Time  `json:"beginning"`
	End       *time.Time `json:"end"`
}

func toJsonSemester(s *models.Semester) *Semester {
	return &Semester{
		s.ID,
		s.Number,
		s.Beginning,
		s.End,
	}
}

func (h *Handler) GetSemester(c echo.Context) error {
	semesterIDParam := c.Param("id")
	semesterID, err := strconv.Atoi(semesterIDParam)

	if err != nil {
		return err
	}

	semesterModel, err := h.useCase.GetSemester(c.Request().Context(), semesterID)

	if err != nil {
		if err == ais.ErrSemesterNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, semesterModel)
}

type manySemestersOutput struct {
	Semesters []*models.Semester `json:"semesters"`
}

func (h *Handler) GetAllSemesters(c echo.Context) error {
	semesters, err := h.useCase.GetAllSemesters(c.Request().Context())

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, manySemestersOutput{Semesters: semesters})
}
