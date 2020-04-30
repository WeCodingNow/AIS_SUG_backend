package models

type Contact struct {
	ID  int
	Def string

	*ContactType
	*Student
}
