package types

import (
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type JSONSemester struct {
	ID        int        `json:"id"`
	Number    int        `json:"number"`
	Beginning time.Time  `json:"beginning"`
	End       *time.Time `json:"end"`

	// ControlEvents []*SemesterShortControlEvent `json:"control_events"`
	// GROUPS
}

func ToJsonSemester(s *models.Semester) *JSONSemester {
	return &JSONSemester{
		ID:        s.ID,
		Number:    s.Number,
		Beginning: s.Beginning,
		End:       s.End,
	}
}
