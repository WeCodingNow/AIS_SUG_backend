package models

type Contact struct {
	ID int
	*ContactType
	*Student
	def string
}
