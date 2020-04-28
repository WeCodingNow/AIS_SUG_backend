package http

import (
	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/labstack/echo"
)

func RegisterHTTPEndpoints(e *echo.Echo, ais ais.UseCase) {
	h := NewHandler(ais)

	e.GET("/cathedra", h.GetAllCathedras)
	e.GET("/cathedra/:id", h.GetCathedra)

	e.GET("/semester", h.GetAllSemesters)
	e.GET("/semester/:id", h.GetSemester)

	e.GET("/discipline", h.GetAllDisciplines)
	e.GET("/discipline/:id", h.GetDiscipline)

	e.GET("/contact_type", h.GetAllContactTypes)
	e.GET("/contact_type/:id", h.GetContactType)

	e.GET("/contact", h.GetAllContacts)
	e.GET("/contact/:id", h.GetContact)

	e.GET("/residence", h.GetAllResidences)
	e.GET("/residence/:id", h.GetResidence)

	e.GET("/control_event_type", h.GetAllControlEventTypes)
	e.GET("/control_event_type/:id", h.GetControlEventType)
}
