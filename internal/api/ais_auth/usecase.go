package ais_auth

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
)

type UseCase interface {
	CreateStudentWithCreds(ctx context.Context, student *models.Student, username, password string, role int) error
}
