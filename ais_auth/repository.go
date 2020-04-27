package ais_auth

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type AisAuthRepository interface {
	// user-student binding
	CreateStudentUserBinding(ctx context.Context, user *models.User, student *models.User) error
}
