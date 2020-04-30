package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/ais/delivery/types"

	"github.com/labstack/echo"
)

func (h *Handler) GetContactType(c echo.Context) error {
	contactTypeIDParam := c.Param("id")
	contactTypeID, err := strconv.Atoi(contactTypeIDParam)

	if err != nil {
		return err
	}

	contactType, err := h.useCase.GetContactType(c.Request().Context(), contactTypeID)

	if err != nil {
		if err == ais.ErrContactTypeNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, types.ToJsonContactType(contactType))
}

type manyContactTypesOutput struct {
	ContactTypes []*types.JSONContactType `json:"contact_types"`
}

func (h *Handler) GetAllContactTypes(c echo.Context) error {
	contactTypes, err := h.useCase.GetAllContactTypes(c.Request().Context())

	if err != nil {
		return err
	}

	jsonContactTypes := make([]*types.JSONContactType, 0, len(contactTypes))
	for _, contactType := range contactTypes {
		jsonContactTypes = append(jsonContactTypes, types.ToJsonContactType(contactType))
	}

	return c.JSON(http.StatusOK, manyContactTypesOutput{ContactTypes: jsonContactTypes})
}
