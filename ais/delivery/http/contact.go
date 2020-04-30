package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/ais/delivery/types"
	"github.com/labstack/echo"
)

// type

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

	return c.JSON(http.StatusOK, types.ToContactJsonContact(contact))
}

type manyContactsOutput struct {
	Contacts []*types.ContactJSONContact `json:"contacts"`
}

func (h *Handler) GetAllContacts(c echo.Context) error {
	contacts, err := h.useCase.GetAllContacts(c.Request().Context())

	if err != nil {
		return err
	}

	jsonContacts := make([]*types.ContactJSONContact, 0, len(contacts))
	for _, contact := range contacts {
		jsonContacts = append(jsonContacts, types.ToContactJsonContact(contact))
	}

	return c.JSON(http.StatusOK, manyContactsOutput{Contacts: jsonContacts})
}
