package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
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

	return c.JSON(http.StatusOK, models.ToJSONCathedra(cathedra, nil))
}

func (h *Handler) GetAllCathedras(c echo.Context) error {
	cathedras, err := h.useCase.GetAllCathedras(c.Request().Context())

	if err != nil {
		return err
	}

	jsonCathedras := make([]models.JSONMap, 0, len(cathedras))
	for _, cathedra := range cathedras {
		jsonCathedras = append(jsonCathedras, models.ToJSONCathedra(cathedra, nil))
	}

	return c.JSON(http.StatusOK, jsonCathedras)
}
