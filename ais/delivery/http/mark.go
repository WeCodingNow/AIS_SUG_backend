package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

// type

type Mark struct {
	ID            int `json:"id"`
	*ControlEvent `json:"control_event"`
	// *ShortStudent `json:"student"`
	Date  time.Time `json:"date"`
	Value int       `json:"value"`
}

// type ShortMark struct {
// 	ID            int `json:"id"`
// 	*ControlEvent `json:"control_event"`
// 	Date          time.Time `json:"date"`
// 	Value         int       `json:"value"`
// }

func toJsonMark(mark *models.Mark) *Mark {
	return &Mark{
		ID:           mark.ID,
		ControlEvent: toJsonControlEvent(mark.ControlEvent),
		// ShortStudent: toJsonShortStudent(mark.Student),
		Date:  mark.Date,
		Value: mark.Value,
	}
}

// func toJsonShortMark(mark *models.Mark) *ShortMark {
// 	return &ShortMark{
// 		ID:           mark.ID,
// 		ControlEvent: toJsonControlEvent(mark.ControlEvent),
// 		Date:         mark.Date,
// 		Value:        mark.Value,
// 	}
// }

func (h *Handler) GetMark(c echo.Context) error {
	markIDParam := c.Param("id")
	markID, err := strconv.Atoi(markIDParam)

	if err != nil {
		return err
	}

	markModel, err := h.useCase.GetMark(c.Request().Context(), markID)

	if err != nil {
		if err == ais.ErrMarkNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, toJsonMark(markModel))
}

type manyMarksOutput struct {
	Marks []*Mark `json:"marks"`
}

func (h *Handler) GetAllMarks(c echo.Context) error {
	markModels, err := h.useCase.GetAllMarks(c.Request().Context())

	if err != nil {
		return err
	}

	marks := make([]*Mark, 0, len(markModels))
	for _, contactModel := range markModels {
		marks = append(marks, toJsonMark(contactModel))
	}

	return c.JSON(http.StatusOK, manyMarksOutput{Marks: marks})
}
