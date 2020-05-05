package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
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

	return c.JSON(http.StatusOK, models.ToJSONControlEventType(controlEventType, nil))
}

func (h *Handler) GetAllControlEventTypes(c echo.Context) error {
	controlEventTypes, err := h.useCase.GetAllControlEventTypes(c.Request().Context())

	if err != nil {
		return err
	}

	jsonControlEventTypes := make([]models.JSONMap, 0, len(controlEventTypes))
	for _, controlEventType := range controlEventTypes {
		jsonControlEventTypes = append(jsonControlEventTypes, models.ToJSONControlEventType(controlEventType, nil))
	}

	return c.JSON(http.StatusOK, jsonControlEventTypes)
}
