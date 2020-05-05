package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

func (h *Handler) GetMark(c echo.Context) error {
	markIDParam := c.Param("id")
	markID, err := strconv.Atoi(markIDParam)

	if err != nil {
		return err
	}

	mark, err := h.useCase.GetMark(c.Request().Context(), markID)

	if err != nil {
		if err == ais.ErrMarkNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, models.ToJSONMark(mark, nil))
}

func (h *Handler) GetAllMarks(c echo.Context) error {
	marks, err := h.useCase.GetAllMarks(c.Request().Context())

	if err != nil {
		return err
	}

	jsonMarks := make([]models.JSONMap, 0, len(marks))
	for _, mark := range marks {
		jsonMarks = append(jsonMarks, models.ToJSONMark(mark, nil))
	}

	return c.JSON(http.StatusOK, jsonMarks)
}
