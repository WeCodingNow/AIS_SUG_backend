package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/ais/delivery/types"
	"github.com/labstack/echo"
)

func (h *Handler) GetControlEventType(c echo.Context) error {
	controlEventTypeIDParam := c.Param("id")
	controlEventTypeID, err := strconv.Atoi(controlEventTypeIDParam)

	if err != nil {
		return err
	}

	controlEventType, err := h.useCase.GetControlEventType(c.Request().Context(), controlEventTypeID)

	if err != nil {
		if err == ais.ErrControlEventTypeNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, types.ToJsonControlEventType(controlEventType))
}

type manyControlEventTypesOutput struct {
	ControlEventTypes []*types.JSONControlEventType `json:"control_event_types"`
}

func (h *Handler) GetAllControlEventTypes(c echo.Context) error {
	jsonControlEventTypes, err := h.useCase.GetAllControlEventTypes(c.Request().Context())

	if err != nil {
		return err
	}

	controlEventTypes := make([]*types.JSONControlEventType, 0, len(jsonControlEventTypes))
	for _, controlEventType := range jsonControlEventTypes {
		controlEventTypes = append(controlEventTypes, types.ToJsonControlEventType(controlEventType))
	}

	return c.JSON(http.StatusOK, manyControlEventTypesOutput{ControlEventTypes: controlEventTypes})
}
