package types

import (
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type JSONMark struct {
	ID    int       `json:"id"`
	Date  time.Time `json:"date"`
	Value int       `json:"value"`
	// *JSONControlEvent `json:"control_event"`
	// *ShortStudent     `json:"student"`
}

func ToJsonMark(mark *models.Mark) *JSONMark {
	return &JSONMark{
		ID:    mark.ID,
		Date:  mark.Date,
		Value: mark.Value,
		// ControlEvent: toJsonControlEvent(mark.ControlEvent),
		// ShortStudent: toJsonShortStudent(mark.Student),
	}
}
