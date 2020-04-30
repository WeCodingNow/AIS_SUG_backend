package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/ais/delivery/types"
	"github.com/labstack/echo"
)

func (h *Handler) GetResidence(c echo.Context) error {
	residenceIDParam := c.Param("id")
	residenceID, err := strconv.Atoi(residenceIDParam)

	if err != nil {
		return err
	}

	residence, err := h.useCase.GetResidence(c.Request().Context(), residenceID)

	if err != nil {
		if err == ais.ErrResidenceNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, types.ToResidenceJSONResidence(residence))
}

type manyResidencesOutput struct {
	Residences []*types.ResidenceJSONResidence `json:"residences"`
}

func (h *Handler) GetAllResidences(c echo.Context) error {
	residences, err := h.useCase.GetAllResidences(c.Request().Context())

	if err != nil {
		return err
	}

	jsonResidences := make([]*types.ResidenceJSONResidence, 0, len(residences))
	for _, residence := range residences {
		jsonResidences = append(jsonResidences, types.ToResidenceJSONResidence(residence))
	}

	return c.JSON(http.StatusOK, manyResidencesOutput{Residences: jsonResidences})
}
