package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

type Cathedra struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

func toJsonCathedra(cathedra *models.Cathedra) *Cathedra {
	return &Cathedra{
		cathedra.ID,
		cathedra.Name,
		cathedra.ShortName,
	}
}

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

	return c.JSON(http.StatusOK, toJsonCathedra(cathedraModel))
}

type manyCathedrasOutput struct {
	Cathedras []*Cathedra `json:"cathedras"`
}

func (h *Handler) GetAllCathedras(c echo.Context) error {
	cathedraModels, err := h.useCase.GetAllCathedras(c.Request().Context())

	cathedras := make([]*Cathedra, 0, len(cathedraModels))
	for _, cathedraModel := range cathedraModels {
		cathedras = append(cathedras, toJsonCathedra(cathedraModel))
	}

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, manyCathedrasOutput{Cathedras: cathedras})
}
