package models

import "time"

type Mark struct {
	ID    int
	Date  time.Time
	Value int

	*ControlEvent
	*Student
}
