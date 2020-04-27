package models

import "time"

type ControlEvent struct {
	ID   int
	date time.Time
	*Discipline
	*ControlEventType
	*Semester
}
