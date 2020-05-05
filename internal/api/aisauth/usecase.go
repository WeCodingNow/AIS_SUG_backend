package aisauth

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
)

type UseCase interface {
	CreateStudentWithCreds(ctx context.Context, user *models.User, role *models.Role, student *models.Student) error
	GetUserRoleID(ctx context.Context, user *models.User) (int, error)
	GetUserStudentID(ctx context.Context, user *models.User) (int, error)
}
