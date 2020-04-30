package types

import (
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type JSONControlEvent struct {
	ID   int       `json:"id"`
	Date time.Time `json:"date"`
	// *JSONControlEventType `json:"type"`
	// *JSONDiscipline       `json:"discipline"`
	// *Semester         `json:"semester"`
}

func ToJsonControlEvent(controlEvent *models.ControlEvent) *JSONControlEvent {
	return &JSONControlEvent{
		ID:   controlEvent.ID,
		Date: controlEvent.Date,
		// ControlEventType: toJsonControlEventType(controlEvent.ControlEventType),
		// Discipline:       toJsonDiscipline(controlEvent.Discipline),
		// Semester:         toJsonSemester(controlEvent.Semester),
	}
}
