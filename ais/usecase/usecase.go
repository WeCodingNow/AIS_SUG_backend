package usecase

import (
	"context"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type AisUseCase struct {
	aisRepo ais.AisRepository
}

func NewAisUseCase(
	aisRepo ais.AisRepository,
) *AisUseCase {
	return &AisUseCase{
		aisRepo,
	}
}

func (ais AisUseCase) CreateSemester(ctx context.Context, number int, beginning time.Time, end *time.Time) error {
	return ais.aisRepo.CreateSemester(ctx, number, beginning, end)
}

func (ais AisUseCase) GetAllCathedras(ctx context.Context) ([]*models.Cathedra, error) {
	return ais.aisRepo.GetAllCathedras(ctx)
}

func (ais AisUseCase) GetCathedra(ctx context.Context, cathedraID int) (*models.Cathedra, error) {
	return ais.aisRepo.GetCathedra(ctx, cathedraID)
}

func (ais AisUseCase) GetSemester(ctx context.Context, semesterID int) (*models.Semester, error) {
	return ais.aisRepo.GetSemester(ctx, semesterID)
}

func (ais AisUseCase) GetAllSemesters(ctx context.Context) ([]*models.Semester, error) {
	return ais.aisRepo.GetAllSemesters(ctx)
}
