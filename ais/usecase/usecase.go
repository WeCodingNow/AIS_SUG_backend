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

func (ais AisUseCase) GetCathedra(ctx context.Context, cathedraID int) (*models.Cathedra, error) {
	return ais.aisRepo.GetCathedra(ctx, cathedraID)
}

func (ais AisUseCase) GetAllCathedras(ctx context.Context) ([]*models.Cathedra, error) {
	return ais.aisRepo.GetAllCathedras(ctx)
}

func (ais AisUseCase) GetSemester(ctx context.Context, semesterID int) (*models.Semester, error) {
	return ais.aisRepo.GetSemester(ctx, semesterID)
}

func (ais AisUseCase) GetAllSemesters(ctx context.Context) ([]*models.Semester, error) {
	return ais.aisRepo.GetAllSemesters(ctx)
}

func (ais AisUseCase) GetContactType(ctx context.Context, contactTypeID int) (*models.ContactType, error) {
	return ais.aisRepo.GetContactType(ctx, contactTypeID)
}

func (ais AisUseCase) GetAllContactTypes(ctx context.Context) ([]*models.ContactType, error) {
	return ais.aisRepo.GetAllContactTypes(ctx)
}

func (ais AisUseCase) GetContact(ctx context.Context, contactID int) (*models.Contact, error) {
	return ais.aisRepo.GetContact(ctx, contactID)
}

func (ais AisUseCase) GetAllContacts(ctx context.Context) ([]*models.Contact, error) {
	return ais.aisRepo.GetAllContacts(ctx)
}

func (ais AisUseCase) GetResidence(ctx context.Context, residenceID int) (*models.Residence, error) {
	return ais.aisRepo.GetResidence(ctx, residenceID)
}

func (ais AisUseCase) GetAllResidences(ctx context.Context) ([]*models.Residence, error) {
	return ais.aisRepo.GetAllResidences(ctx)
}

func (ais AisUseCase) GetControlEventType(ctx context.Context, controlEventTypeID int) (*models.ControlEventType, error) {
	return ais.aisRepo.GetControlEventType(ctx, controlEventTypeID)
}

func (ais AisUseCase) GetAllControlEventTypes(ctx context.Context) ([]*models.ControlEventType, error) {
	return ais.aisRepo.GetAllControlEventTypes(ctx)
}

func (ais AisUseCase) GetDiscipline(ctx context.Context, controlEventTypeID int) (*models.Discipline, error) {
	return ais.aisRepo.GetDiscipline(ctx, controlEventTypeID)
}

func (ais AisUseCase) GetAllDisciplines(ctx context.Context) ([]*models.Discipline, error) {
	return ais.aisRepo.GetAllDisciplines(ctx)
}
