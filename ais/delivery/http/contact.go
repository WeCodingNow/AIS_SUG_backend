package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

// type

type Contact struct {
	ID          int    `json:"id"`
	ContactType string `json:"type"`
	Def         string `json:"def"`
}

func toJsonContact(contact *models.Contact) *Contact {
	return &Contact{
		ID:          contact.ID,
		Def:         contact.Def,
		ContactType: contact.ContactType.Def,
	}
}

func (h *Handler) GetContact(c echo.Context) error {
	contactIDParam := c.Param("id")
	contactID, err := strconv.Atoi(contactIDParam)

	if err != nil {
		return err
	}

	contactModel, err := h.useCase.GetContact(c.Request().Context(), contactID)

	if err != nil {
		if err == ais.ErrContactNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, toJsonContact(contactModel))
}

type manyContactsOutput struct {
	Contacts []*Contact `json:"contacts"`
}

func (h *Handler) GetAllContacts(c echo.Context) error {
	contactModels, err := h.useCase.GetAllContacts(c.Request().Context())

	if err != nil {
		return err
	}

	contacts := make([]*Contact, 0, len(contactModels))
	for _, contactModel := range contactModels {
		contacts = append(contacts, toJsonContact(contactModel))
	}

	return c.JSON(http.StatusOK, manyContactsOutput{Contacts: contacts})
}
