package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

func (h *Handler) GetCathedra(c echo.Context) error {
	cathedraIDParam := c.Param("id")
	cathedraID, err := strconv.Atoi(cathedraIDParam)

	if err != nil {
		return err
	}

	cathedraModel, err := h.useCase.GetCathedra(c.Request().Context(), cathedraID)

	if err != nil {
		if err == ais.ErrCathedraNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, cathedraModel)
}

type manyCathedrasOutput struct {
	Cathedras []*models.Cathedra `json:"cathedras"`
}

func (h *Handler) GetAllCathedras(c echo.Context) error {
	cathedras, err := h.useCase.GetAllCathedras(c.Request().Context())

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, manyCathedrasOutput{Cathedras: cathedras})
}
