package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

// type

// TODO: добавить сюда студентов
type Residence struct {
	ID        int      `json:"id"`
	Address   string   `json:"address"`
	City      string   `json:"city"`
	Community bool     `json:"community"`
	Students  []string `json:"students"`
	// Students  []*Student
}

type ShortResidence struct {
	ID        int    `json:"id"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Community bool   `json:"community"`
}

func toJsonResidence(residence *models.Residence) *Residence {
	return &Residence{
		ID:        residence.ID,
		Address:   residence.Address,
		City:      residence.City,
		Community: residence.Community,
		Students:  []string{"student1", "student2"},
	}
}

func toJsonShortResidence(residence *models.Residence) *ShortResidence {
	return &ShortResidence{
		ID:        residence.ID,
		Address:   residence.Address,
		City:      residence.City,
		Community: residence.Community,
	}
}

func (h *Handler) GetResidence(c echo.Context) error {
	residenceIDParam := c.Param("id")
	residenceID, err := strconv.Atoi(residenceIDParam)

	if err != nil {
		return err
	}

	residenceModel, err := h.useCase.GetResidence(c.Request().Context(), residenceID)

	if err != nil {
		if err == ais.ErrResidenceNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, toJsonResidence(residenceModel))
}

type manyResidencesOutput struct {
	Residences []*Residence `json:"residences"`
}

func (h *Handler) GetAllResidences(c echo.Context) error {
	residenceModels, err := h.useCase.GetAllResidences(c.Request().Context())

	if err != nil {
		return err
	}

	residences := make([]*Residence, 0, len(residenceModels))
	for _, residenceModel := range residenceModels {
		residences = append(residences, toJsonResidence(residenceModel))
	}

	return c.JSON(http.StatusOK, manyResidencesOutput{Residences: residences})
}
