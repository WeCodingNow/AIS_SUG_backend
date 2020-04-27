package models

// binding between user model and student
type UserStudent struct {
	ID int
	User
	Student
}
