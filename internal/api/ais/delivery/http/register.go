package http

import (
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/aisauth"
	aisauthhttp "github.com/WeCodingNow/AIS_SUG_backend/internal/api/aisauth/delivery/http"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth"
	authhttp "github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth/delivery/http"
	"github.com/labstack/echo"
)

func RegisterHTTPEndpoints(e *echo.Echo, ais ais.UseCase, au auth.UseCase, auc aisauth.UseCase) {
	h := NewHandler(ais)

	e.GET("/contact_type", h.GetAllContactTypes)
	e.GET("/contact_type/:id", h.GetContactType)

	e.GET("/control_event_type", h.GetAllControlEventTypes)
	e.GET("/control_event_type/:id", h.GetControlEventType)

	e.GET("/contact", h.GetAllContacts)
	e.GET("/contact/:id", h.GetContact)

	e.GET("/control_event", h.GetAllControlEvents)
	e.GET("/control_event/:id", h.GetControlEvent)

	e.GET("/discipline", h.GetAllDisciplines)
	e.GET("/discipline/:id", h.GetDiscipline)

	e.POST("/mark", h.CreateMark, authhttp.MakeAuthMiddleware(au), aisauthhttp.MakeRBACMiddleware(auc, []int{auth.AdminClass}))
	e.GET("/mark", h.GetAllMarks)
	e.GET("/mark/:id", h.GetMark)

	e.POST("/residence", h.CreateResidence)
	e.GET("/residence", h.GetAllResidences)
	e.GET("/residence/:id", h.GetResidence)

	e.POST("/backlog", h.CreateBacklog)
	e.GET("/backlog", h.GetAllBacklogs)
	e.GET("/backlog/:id", h.GetBacklog)

	e.POST("/student", h.CreateStudent)
	e.GET("/student", h.GetAllStudents)
	e.GET("/student/:id", h.GetStudent)

	e.GET("/group", h.GetAllGroups)
	e.GET("/group/:id", h.GetGroup)

	e.GET("/cathedra", h.GetAllCathedras)
	e.GET("/cathedra/:id", h.GetCathedra)

	e.GET("/semester", h.GetAllSemesters)
	e.GET("/semester/:id", h.GetSemester)
}
