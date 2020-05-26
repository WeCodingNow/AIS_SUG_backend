package aisauth

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
)

type StudentWithUserAndRole struct {
	StudentID int  `json:"student_id"`
	UserID    *int `json:"user_id"`
	RoleID    *int `json:"role_id"`
}

type StudentInfo struct {
	UserID    int  `json:"user_id"`
	StudentID *int `json:"student_id"`
	GroupID   *int `json:"group_id"`
}

type UseCase interface {
	CreateStudentWithCreds(ctx context.Context, user *models.User, role *models.Role, student *models.Student) error
	AssignOrCreateStudentWithCreds(ctx context.Context, user *models.User, targetStudent *models.Student) error

	GetUserRoleID(ctx context.Context, user *models.User) (int, error)
	GetUserRole(ctx context.Context, userID int) (*models.Role, error)
	GetRoles(ctx context.Context) ([]*models.Role, error)

	GetInfo(ctx context.Context, userID int) (*StudentInfo, error)

	GetUserStudentID(ctx context.Context, userID int) (*int, error)
	GetStudentsWithUsers(ctx context.Context) ([]StudentWithUserAndRole, error)

	PromoteUser(ctx context.Context, userID int, roleID int) error
}
