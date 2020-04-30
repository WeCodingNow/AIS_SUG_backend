package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/ais/delivery/types"
	"github.com/labstack/echo"
)

func (h *Handler) GetCathedra(c echo.Context) error {
	cathedraIDParam := c.Param("id")
	cathedraID, err := strconv.Atoi(cathedraIDParam)

	if err != nil {
		return err
	}

	cathedra, err := h.useCase.GetCathedra(c.Request().Context(), cathedraID)

	if err != nil {
		if err == ais.ErrCathedraNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, types.ToJsonCathedra(cathedra))
}

type manyCathedrasOutput struct {
	Cathedras []*types.JSONCathedra `json:"cathedras"`
}

func (h *Handler) GetAllCathedras(c echo.Context) error {
	cathedras, err := h.useCase.GetAllCathedras(c.Request().Context())

	jsonCathedras := make([]*types.JSONCathedra, 0, len(cathedras))
	for _, cathedra := range cathedras {
		jsonCathedras = append(jsonCathedras, types.ToJsonCathedra(cathedra))
	}

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, manyCathedrasOutput{Cathedras: jsonCathedras})
}
