package models

// binding between user model and role
type UserRole struct {
	ID        int
	UserID    int
	RoleID    int
	Confirmed bool
}
