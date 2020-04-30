package types

import (
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type JSONControlEvent struct {
	ID               int                   `json:"id"`
	Date             time.Time             `json:"date"`
	ControlEventType *JSONControlEventType `json:"type"`
}

func toJsonControlEvent(controlEvent *models.ControlEvent) *JSONControlEvent {
	return &JSONControlEvent{
		ID:               controlEvent.ID,
		Date:             controlEvent.Date,
		ControlEventType: ToJsonControlEventType(controlEvent.ControlEventType),
	}
}

type ControlEventJSONMark struct {
	*JSONMark
	Student *MarkJSONStudent `json:"student"`
}

func toControlEventJsonMark(mark *models.Mark) *ControlEventJSONMark {
	return &ControlEventJSONMark{
		JSONMark: toJsonMark(mark),
		Student:  toMarkJsonStudent(mark.Student),
	}
}

type ControlEventJSONControlEvent struct {
	*JSONControlEvent
	Discipline *JSONDiscipline
	Marks      []*ControlEventJSONMark
}

func ToControlEventJSONControlEvent(controlEvent *models.ControlEvent) *ControlEventJSONControlEvent {
	markJSONs := make([]*ControlEventJSONMark, 0, len(controlEvent.Marks))

	for _, mark := range controlEvent.Marks {
		markJSONs = append(markJSONs, toControlEventJsonMark(mark))
	}

	return &ControlEventJSONControlEvent{
		JSONControlEvent: toJsonControlEvent(controlEvent),
		Discipline:       ToJsonDiscipline(controlEvent.Discipline),
		Marks:            markJSONs,
	}

}
