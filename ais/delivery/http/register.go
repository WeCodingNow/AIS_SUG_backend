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
}
