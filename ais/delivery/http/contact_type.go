package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

type ContactType struct {
	ID  int    `json:"id"`
	Def string `json:"def"`
}

func toJsonContactType(contactType *models.ContactType) *ContactType {
	return &ContactType{
		ID:  contactType.ID,
		Def: contactType.Def,
	}
}

func (h *Handler) GetContactType(c echo.Context) error {
	contactTypeIDParam := c.Param("id")
	contactTypeID, err := strconv.Atoi(contactTypeIDParam)

	if err != nil {
		return err
	}

	contactTypeModel, err := h.useCase.GetContactType(c.Request().Context(), contactTypeID)

	if err != nil {
		if err == ais.ErrContactTypeNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, contactTypeModel)
}

type manyContactTypesOutput struct {
	ContactTypes []*ContactType `json:"contact_types"`
}

func (h *Handler) GetAllContactTypes(c echo.Context) error {
	contactTypeModels, err := h.useCase.GetAllContactTypes(c.Request().Context())

	if err != nil {
		return err
	}

	contactTypes := make([]*ContactType, 0, len(contactTypeModels))
	for _, contactTypeModel := range contactTypeModels {
		contactTypes = append(contactTypes, toJsonContactType(contactTypeModel))
	}

	return c.JSON(http.StatusOK, manyContactTypesOutput{ContactTypes: contactTypes})
}
