package models

import "time"

type Mark struct {
	ID int
	*ControlEvent
	*Student
	date  time.Time
	value int
}
