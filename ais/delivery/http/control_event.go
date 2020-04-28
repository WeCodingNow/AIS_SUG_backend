package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

type ControlEvent struct {
	ID                int `json:"id"`
	*ControlEventType `json:"type"`
	*Discipline       `json:"discipline"`
	*Semester         `json:"semester"`
	Date              time.Time `json:"date"`
}

func toJsonControlEvent(controlEvent *models.ControlEvent) *ControlEvent {
	return &ControlEvent{
		ID:               controlEvent.ID,
		ControlEventType: toJsonControlEventType(controlEvent.ControlEventType),
		Discipline:       toJsonDiscipline(controlEvent.Discipline),
		Semester:         toJsonSemester(controlEvent.Semester),
		Date:             controlEvent.Date,
		// ID:  controlEventType.ID,
		// Def: controlEventType.Def,
	}
}
func (h *Handler) GetControlEvent(c echo.Context) error {
	controlEventIDParam := c.Param("id")
	controlEventID, err := strconv.Atoi(controlEventIDParam)

	if err != nil {
		return err
	}

	controlEventModel, err := h.useCase.GetControlEvent(c.Request().Context(), controlEventID)

	if err != nil {
		if err == ais.ErrControlEventNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, toJsonControlEvent(controlEventModel))
}

type manyControlEventsOutput struct {
	ControlEvents []*ControlEvent `json:"control_events"`
}

func (h *Handler) GetAllControlEvents(c echo.Context) error {
	controlEventModels, err := h.useCase.GetAllControlEvents(c.Request().Context())

	if err != nil {
		return err
	}

	controlEvents := make([]*ControlEvent, 0, len(controlEventModels))
	for _, controlEventModel := range controlEventModels {
		controlEvents = append(controlEvents, toJsonControlEvent(controlEventModel))
	}

	return c.JSON(http.StatusOK, manyControlEventsOutput{ControlEvents: controlEvents})
}
