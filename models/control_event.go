package models

import "time"

type ControlEvent struct {
	ID   int
	Date time.Time

	ControlEventType *ControlEventType
	Discipline       *Discipline
	// Semester         *Semester
	Marks []*Mark
}
