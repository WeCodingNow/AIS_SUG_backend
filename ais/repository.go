package ais

import (
	"context"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type AisRepository interface {
	// CreateCathedra(ctx context.Context, name, shortName string) error
	// UpdateCathedra(ctx context.Context, group *models.Cathedra) error
	// DeleteCathedra(ctx context.Context, cathedra *models.Cathedra) error
	GetCathedra(ctx context.Context, cathedraID int) (*models.Cathedra, error)
	GetAllCathedras(ctx context.Context) ([]*models.Cathedra, error)

	CreateSemester(ctx context.Context, number int, beginning time.Time, end *time.Time) error
	// UpdateSemester(ctx context.Context, group *models.Group) error
	// // DeleteSemester(ctx context.Context, semester *models.Semester) error
	GetSemester(ctx context.Context, semesterID int) (*models.Semester, error)
	GetAllSemesters(ctx context.Context) ([]*models.Semester, error)

	// CreateGroup(ctx context.Context, cathedraID int, number int) error
	// UpdateGroup(ctx context.Context, group *models.Group) error
	// // DeleteGroup(ctx context.Context, group *models.Group) error
	// GetGroup(ctx context.Context, studentID int) (*models.Group, error)

	// CreateStudent(ctx context.Context, name, secondName string, thirdName *string, groupID int) error
	// UpdateStudent(ctx context.Context, student *models.Student) error
	// // DeleteStudent(ctx context.Context, student *models.Student) error
	// GetStudent(ctx context.Context, studentID int) (*models.Student, error)
}
