package models

type Group struct {
	ID        int
	number    int
	Cathedra  *Cathedra
	Semesters []*Semester
	Students  []*Student
}
