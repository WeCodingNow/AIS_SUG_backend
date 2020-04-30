package models

type Group struct {
	ID     int
	Number int

	Cathedra  *Cathedra
	Students  []*Student
	Semesters []*Semester
}
