package models

import "time"

type Mark struct {
	ID int
	*ControlEvent
	*Student
	Date  time.Time
	Value int
}
