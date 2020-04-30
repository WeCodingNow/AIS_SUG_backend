package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/ais/delivery/types"
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

	return c.JSON(http.StatusOK, types.ToMarkJsonMark(mark))
}

type manyMarksOutput struct {
	Marks []*types.MarkJSONMark `json:"marks"`
}

func (h *Handler) GetAllMarks(c echo.Context) error {
	markModels, err := h.useCase.GetAllMarks(c.Request().Context())

	if err != nil {
		return err
	}

	marks := make([]*types.MarkJSONMark, 0, len(markModels))
	for _, mark := range markModels {
		marks = append(marks, types.ToMarkJsonMark(mark))
	}

	return c.JSON(http.StatusOK, manyMarksOutput{Marks: marks})
}
