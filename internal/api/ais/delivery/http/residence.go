package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/labstack/echo"
)

type PostResidence struct {
	Address   string `json:"address"`
	City      string `json:"city"`
	Community bool   `json:"community"`
}

func (h *Handler) CreateResidence(c echo.Context) error {
	var inp PostResidence
	err := c.Bind(&inp)

	if err != nil {
		return err
	}

	newResidence, err := h.useCase.CreateResidence(c.Request().Context(), inp.Address, inp.City, inp.Community)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.ToJSONResidence(newResidence, nil))
}

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

	return c.JSON(http.StatusOK, models.ToJSONResidence(residence, nil))
}

func (h *Handler) GetAllResidences(c echo.Context) error {
	residences, err := h.useCase.GetAllResidences(c.Request().Context())

	if err != nil {
		return err
	}

	jsonResidences := make([]models.JSONMap, 0, len(residences))
	for _, residence := range residences {
		jsonResidences = append(jsonResidences, models.ToJSONResidence(residence, nil))
	}

	return c.JSON(http.StatusOK, jsonResidences)
}
