package ais

import (
	"context"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type UseCase interface {
	GetAllCathedras(ctx context.Context) ([]*models.Cathedra, error)
	GetCathedra(ctx context.Context, cathedraID int) (*models.Cathedra, error)

	CreateSemester(ctx context.Context, number int, beginning time.Time, end *time.Time) error
	GetAllSemesters(ctx context.Context) ([]*models.Semester, error)
	GetSemester(ctx context.Context, semesterID int) (*models.Semester, error)
}
