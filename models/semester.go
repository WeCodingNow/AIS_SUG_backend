package models

import "time"

type Semester struct {
	ID        int        `json:"id"`
	Number    int        `json:"number"`
	Beginning time.Time  `json:"beginning"`
	End       *time.Time `json:"end"`

	Groups []*Group
	// ControlEvents []*ControlEvent
}
