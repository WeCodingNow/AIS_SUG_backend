package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

func (h *Handler) GetGroup(c echo.Context) error {
	groupIDParam := c.Param("id")
	groupID, err := strconv.Atoi(groupIDParam)

	if err != nil {
		return err
	}

	group, err := h.useCase.GetGroup(c.Request().Context(), groupID)

	if err != nil {
		if err == ais.ErrGroupNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, models.ToJSONGroup(group, nil))
}

func (h *Handler) GetAllGroups(c echo.Context) error {
	groups, err := h.useCase.GetAllGroups(c.Request().Context())

	if err != nil {
		return err
	}

	jsonGroups := make([]models.JSONMap, 0, len(groups))
	for _, group := range groups {
		jsonGroups = append(jsonGroups, models.ToJSONGroup(group, nil))
	}

	return c.JSON(http.StatusOK, jsonGroups)
}
