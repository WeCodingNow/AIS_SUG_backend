package models

type Discipline struct {
	ID    int
	Name  string
	Hours int

	ControlEvents []*ControlEvent
}
