package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/labstack/echo"
)

func (h *Handler) GetControlEvent(c echo.Context) error {
	controlEventIDParam := c.Param("id")
	controlEventID, err := strconv.Atoi(controlEventIDParam)

	if err != nil {
		return err
	}

	controlEvent, err := h.useCase.GetControlEvent(c.Request().Context(), controlEventID)

	if err != nil {
		if err == ais.ErrControlEventNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, models.ToJSONControlEvent(controlEvent, nil))
}

func (h *Handler) GetAllControlEvents(c echo.Context) error {
	controlEvents, err := h.useCase.GetAllControlEvents(c.Request().Context())

	if err != nil {
		return err
	}

	jsonControlEvents := make([]models.JSONMap, 0, len(controlEvents))
	for _, controlEvent := range controlEvents {
		jsonControlEvents = append(jsonControlEvents, models.ToJSONControlEvent(controlEvent, nil))
	}

	return c.JSON(http.StatusOK, jsonControlEvents)
}
