package models

// binding between user model and role
type UserRoleBinding struct {
	ID   int
	User *User
	Role *Role
}
