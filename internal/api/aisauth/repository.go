package aisauth

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
)

type AisAuthRepository interface {
	CreateUserRoleBinding(ctx context.Context, user *models.User, role *models.Role) error
	CreateUserStudentBinding(ctx context.Context, user *models.User, student *models.Student) error

	GetUserRoleID(ctx context.Context, user *models.User) (int, error)
	GetUserRole(ctx context.Context, userID int) (*models.Role, error)

	GetUserStudentID(ctx context.Context, userID int) (*int, error)
	GetStudentUserID(ctx context.Context, student *models.Student) (int, error)

	GetRoles(ctx context.Context) ([]*models.Role, error)

	PromoteUser(ctx context.Context, userID int, roleID int) error
}
