package postgres

import "time"

type ControlEvent struct {
	ID int
	*ControlEventType
	*Discipline
	*Semester
	date time.Time
}
