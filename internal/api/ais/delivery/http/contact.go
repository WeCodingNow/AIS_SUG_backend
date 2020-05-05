package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/labstack/echo"
)

func (h *Handler) GetContact(c echo.Context) error {
	contactIDParam := c.Param("id")
	contactID, err := strconv.Atoi(contactIDParam)

	if err != nil {
		return err
	}

	contact, err := h.useCase.GetContact(c.Request().Context(), contactID)

	if err != nil {
		if err == ais.ErrContactNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, models.ToJSONContact(contact, nil))
}

func (h *Handler) GetAllContacts(c echo.Context) error {
	contacts, err := h.useCase.GetAllContacts(c.Request().Context())

	if err != nil {
		return err
	}

	jsonContacts := make([]models.JSONMap, 0, len(contacts))
	for _, contact := range contacts {
		jsonContacts = append(jsonContacts, models.ToJSONContact(contact, nil))
	}

	return c.JSON(http.StatusOK, jsonContacts)
}
