package models

import "time"

type Mark struct {
	ID    int
	date  time.Time
	value int
	*ControlEvent
	*Student
}
