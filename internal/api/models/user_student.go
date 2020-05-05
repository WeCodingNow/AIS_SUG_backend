package models

// binding between user model and student
type UserStudentBinding struct {
	ID      int
	User    *User
	Student *Student
}
