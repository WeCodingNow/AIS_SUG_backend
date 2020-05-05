package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

func (h *Handler) GetContactType(c echo.Context) error {
	contactTypeIDParam := c.Param("id")
	contactTypeID, err := strconv.Atoi(contactTypeIDParam)

	if err != nil {
		return err
	}

	contactType, err := h.useCase.GetContactType(c.Request().Context(), contactTypeID)

	if err != nil {
		if err == ais.ErrContactTypeNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, models.ToJSONContactType(contactType, nil))
}

func (h *Handler) GetAllContactTypes(c echo.Context) error {
	contactTypes, err := h.useCase.GetAllContactTypes(c.Request().Context())

	if err != nil {
		return err
	}

	jsonContactTypes := make([]models.JSONMap, 0, len(contactTypes))
	for _, controlEventType := range contactTypes {
		jsonContactTypes = append(jsonContactTypes, models.ToJSONContactType(controlEventType, nil))
	}

	return c.JSON(http.StatusOK, jsonContactTypes)
}
