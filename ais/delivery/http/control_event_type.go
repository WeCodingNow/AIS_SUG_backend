package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

type ControlEventType struct {
	ID  int    `json:"id"`
	Def string `json:"def"`
}

func toJsonControlEventType(controlEventType *models.ControlEventType) *ControlEventType {
	return &ControlEventType{
		ID:  controlEventType.ID,
		Def: controlEventType.Def,
	}
}
func (h *Handler) GetControlEventType(c echo.Context) error {
	controlEventTypeIDParam := c.Param("id")
	controlEventTypeID, err := strconv.Atoi(controlEventTypeIDParam)

	if err != nil {
		return err
	}

	controlEventTypeModel, err := h.useCase.GetControlEventType(c.Request().Context(), controlEventTypeID)

	if err != nil {
		if err == ais.ErrControlEventTypeNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, toJsonControlEventType(controlEventTypeModel))
}

type manyControlEventTypesOutput struct {
	ControlEventTypes []*ControlEventType `json:"control_event_types"`
}

func (h *Handler) GetAllControlEventTypes(c echo.Context) error {
	controlEventTypeModels, err := h.useCase.GetAllControlEventTypes(c.Request().Context())

	if err != nil {
		return err
	}

	controlEventTypes := make([]*ControlEventType, 0, len(controlEventTypeModels))
	for _, controlEventTypeModel := range controlEventTypeModels {
		controlEventTypes = append(controlEventTypes, toJsonControlEventType(controlEventTypeModel))
	}

	return c.JSON(http.StatusOK, manyControlEventTypesOutput{ControlEventTypes: controlEventTypes})
}
