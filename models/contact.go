package models

type Contact struct {
	ID int
	*ContactType
	*Student
	Def string
}
