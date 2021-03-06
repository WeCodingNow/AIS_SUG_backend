package usecase

import (
	"context"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
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

func (ais AisUseCase) GetGroup(ctx context.Context, groupID int) (*models.Group, error) {
	return ais.aisRepo.GetGroup(ctx, groupID)
}

func (ais AisUseCase) GetAllGroups(ctx context.Context) ([]*models.Group, error) {
	return ais.aisRepo.GetAllGroups(ctx)
}

func (ais AisUseCase) GetStudent(ctx context.Context, studentID int) (*models.Student, error) {
	return ais.aisRepo.GetStudent(ctx, studentID)
}

func (ais AisUseCase) CreateStudent(ctx context.Context, name, secondName string, thirdName *string, groupID, residenceID int) (*models.Student, error) {
	return ais.aisRepo.CreateStudent(ctx, name, secondName, thirdName, groupID, residenceID)
}

func (ais AisUseCase) GetAllStudents(ctx context.Context) ([]*models.Student, error) {
	return ais.aisRepo.GetAllStudents(ctx)
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

func (ais AisUseCase) CreateResidence(ctx context.Context, address, city string, community bool) (*models.Residence, error) {
	return ais.aisRepo.CreateResidence(ctx, address, city, community)
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

func (ais AisUseCase) GetControlEvent(ctx context.Context, controlEventID int) (*models.ControlEvent, error) {
	return ais.aisRepo.GetControlEvent(ctx, controlEventID)
}

func (ais AisUseCase) GetAllControlEvents(ctx context.Context) ([]*models.ControlEvent, error) {
	return ais.aisRepo.GetAllControlEvents(ctx)
}

func (ais AisUseCase) GetDiscipline(ctx context.Context, controlEventTypeID int) (*models.Discipline, error) {
	return ais.aisRepo.GetDiscipline(ctx, controlEventTypeID)
}

func (ais AisUseCase) GetAllDisciplines(ctx context.Context) ([]*models.Discipline, error) {
	return ais.aisRepo.GetAllDisciplines(ctx)
}

func (ais AisUseCase) CreateMark(ctx context.Context, date time.Time, value, controlEventID, studentID int) (*models.Mark, error) {
	return ais.aisRepo.CreateMark(ctx, date, value, controlEventID, studentID)
}

func (ais AisUseCase) GetMark(ctx context.Context, markID int) (*models.Mark, error) {
	return ais.aisRepo.GetMark(ctx, markID)
}

func (ais AisUseCase) GetAllMarks(ctx context.Context) ([]*models.Mark, error) {
	return ais.aisRepo.GetAllMarks(ctx)
}

func (ais AisUseCase) CreateBacklog(ctx context.Context, description string, disciplineID, studentID int) (*models.Backlog, error) {
	return ais.aisRepo.CreateBacklog(ctx, description, disciplineID, studentID)
}

func (ais AisUseCase) GetBacklog(ctx context.Context, backlogID int) (*models.Backlog, error) {
	return ais.aisRepo.GetBacklog(ctx, backlogID)
}

func (ais AisUseCase) GetAllBacklogs(ctx context.Context) ([]*models.Backlog, error) {
	return ais.aisRepo.GetAllBacklogs(ctx)
}
