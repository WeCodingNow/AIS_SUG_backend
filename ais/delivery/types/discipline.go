package types

import (
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type JSONDiscipline struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Hours int    `json:"hours"`
	// ControlEvents []*DisciplineShortControlEvent `json:"control_events"`
}

func ToJsonDiscipline(discipline *models.Discipline) *JSONDiscipline {
	// jsonControlEvents := make([]*DisciplineShortControlEvent, 0)

	// for _, controlEvent := range discipline.ControlEvents {
	// 	jsonControlEvents = append(jsonControlEvents, toJsonDisciplineShortControlEvent(controlEvent))
	// }

	return &JSONDiscipline{
		ID:    discipline.Hours,
		Name:  discipline.Name,
		Hours: discipline.Hours,
		// ControlEvents: jsonControlEvents,
	}
}
