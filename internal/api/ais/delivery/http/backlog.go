package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/labstack/echo"
)

type PostBacklog struct {
	Description  string `json:"desc"`
	StudentID    int    `json:"student_id"`
	DisciplineID int    `json:"discipline_id"`
}

func (h *Handler) CreateBacklog(c echo.Context) error {
	var inp PostBacklog
	err := c.Bind(&inp)

	if err != nil {
		return err
	}

	newBacklog, err := h.useCase.CreateBacklog(c.Request().Context(), inp.Description, inp.DisciplineID, inp.StudentID)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.ToJsonBacklog(newBacklog, nil))
}

func (h *Handler) GetBacklog(c echo.Context) error {
	backlogIDParam := c.Param("id")
	backlogID, err := strconv.Atoi(backlogIDParam)

	if err != nil {
		return err
	}

	backlog, err := h.useCase.GetBacklog(c.Request().Context(), backlogID)

	if err != nil {
		if err == ais.ErrBacklogNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, models.ToJsonBacklog(backlog, nil))
}

func (h *Handler) GetAllBacklogs(c echo.Context) error {
	backlogs, err := h.useCase.GetAllBacklogs(c.Request().Context())

	if err != nil {
		return err
	}

	jsonBacklogs := make([]models.JSONMap, 0, len(backlogs))
	for _, backlog := range backlogs {
		jsonBacklogs = append(jsonBacklogs, models.ToJsonBacklog(backlog, nil))
	}

	return c.JSON(http.StatusOK, jsonBacklogs)
}
