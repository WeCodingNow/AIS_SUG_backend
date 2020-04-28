package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

type Group struct {
	ID       int        `json:"id"`
	Cathedra *Cathedra  `json:"cathedra"`
	Number   int        `json:"number"`
	Students []*Student `json:"students"`
	// Semesters []*Semester
}

func toJsonGroup(group *models.Group) *Group {
	studentJSONs := make([]*Student, 0, len(group.Students))

	for _, studentModel := range group.Students {
		log.Print(*studentModel)
		studentJSONs = append(studentJSONs, toJsonStudent(studentModel))
	}

	return &Group{
		ID:       group.ID,
		Cathedra: toJsonCathedra(group.Cathedra),
		Students: studentJSONs,
		Number:   group.Number,
	}
}

func (h *Handler) GetGroup(c echo.Context) error {
	groupIDParam := c.Param("id")
	groupID, err := strconv.Atoi(groupIDParam)

	if err != nil {
		return err
	}

	groupModel, err := h.useCase.GetGroup(c.Request().Context(), groupID)

	if err != nil {
		if err == ais.ErrGroupNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, toJsonGroup(groupModel))
}

type manyGroupsOutput struct {
	Groups []*Group `json:"groups"`
}

func (h *Handler) GetAllGroups(c echo.Context) error {
	groupModels, err := h.useCase.GetAllGroups(c.Request().Context())

	if err != nil {
		return err
	}

	groups := make([]*Group, 0, len(groupModels))
	for _, groupModel := range groupModels {
		groups = append(groups, toJsonGroup(groupModel))
	}

	return c.JSON(http.StatusOK, manyGroupsOutput{Groups: groups})
}
