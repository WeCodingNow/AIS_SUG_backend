package usecase

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
)

type AisAuthUseCase struct {
	aisUC  ais.UseCase
	authUC auth.UseCase
}

func NewAisAuthUseCase(
	aisUC ais.UseCase,
	authUC auth.UseCase,
) *AisAuthUseCase {
	return &AisAuthUseCase{
		aisUC,
		authUC,
	}
}

func (ais AisAuthUseCase) CreateStudentWithCreds(ctx context.Context, student *models.Student, username, password string, role int) error {
	return nil
}
