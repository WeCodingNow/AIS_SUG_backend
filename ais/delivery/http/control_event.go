package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/ais/delivery/types"
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

	return c.JSON(http.StatusOK, types.ToControlEventJSONControlEvent(controlEvent))
}

type manyControlEventsOutput struct {
	ControlEvents []*types.ControlEventJSONControlEvent `json:"control_events"`
}

func (h *Handler) GetAllControlEvents(c echo.Context) error {
	controlEvents, err := h.useCase.GetAllControlEvents(c.Request().Context())

	if err != nil {
		return err
	}

	jsonControlEvents := make([]*types.ControlEventJSONControlEvent, 0, len(controlEvents))
	for _, controlEvent := range controlEvents {
		jsonControlEvents = append(jsonControlEvents, types.ToControlEventJSONControlEvent(controlEvent))
	}

	return c.JSON(http.StatusOK, manyControlEventsOutput{ControlEvents: jsonControlEvents})
}
