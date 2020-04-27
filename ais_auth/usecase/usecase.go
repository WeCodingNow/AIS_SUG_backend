package usecase

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/auth"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
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
