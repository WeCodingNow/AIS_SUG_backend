package models

type Group struct {
	ID        int
	Cathedra  *Cathedra
	number    int
	Semesters []*Semester
	Students  []*Student
}
