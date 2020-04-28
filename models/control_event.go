package models

import "time"

type ControlEvent struct {
	ID int
	*ControlEventType
	*Discipline
	*Semester
	Date time.Time
}
