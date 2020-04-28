package models

type Group struct {
	ID        int
	Cathedra  *Cathedra
	Number    int
	Semesters []*Semester
	Students  []*Student
}
