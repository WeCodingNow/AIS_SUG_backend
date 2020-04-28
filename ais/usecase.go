package ais

import (
	"context"
	// "time"

	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type UseCase interface {
	GetCathedra(ctx context.Context, cathedraID int) (*models.Cathedra, error)
	GetAllCathedras(ctx context.Context) ([]*models.Cathedra, error)

	// CreateSemester(ctx context.Context, number int, beginning time.Time, end *time.Time) error
	GetSemester(ctx context.Context, semesterID int) (*models.Semester, error)
	GetAllSemesters(ctx context.Context) ([]*models.Semester, error)

	// GetAllGroups(ctx context.Context) ([]*models.Group, error)
	// GetGroup(ctx context.Context, groupID int) (*models.Group, error)

	// GetAllStudents(ctx context.Context) ([]*models.Student, error)
	// GetStudent(ctx context.Context, studentID int) (*models.Student, error)

	GetContactType(ctx context.Context, contactTypeID int) (*models.ContactType, error)
	GetAllContactTypes(ctx context.Context) ([]*models.ContactType, error)

	GetContact(ctx context.Context, contactID int) (*models.Contact, error)
	GetAllContacts(ctx context.Context) ([]*models.Contact, error)

	// GetResidence(ctx context.Context, residenceID int) (*models.Residence, error)
	// GetAllResidences(ctx context.Context) ([]*models.Residence, error)

	// GetDiscipline(ctx context.Context, disciplineID int) (*models.Discipline, error)
	// GetAllDisciplines(ctx context.Context) ([]*models.Discipline, error)

	GetControlEventType(ctx context.Context, controlEventTypeID int) (*models.ControlEventType, error)
	GetAllControlEventTypes(ctx context.Context) ([]*models.ControlEventType, error)

	// GetControlEvent(ctx context.Context, controlEventID int) (*models.ControlEvent, error)
	// GetAllControlEvents(ctx context.Context) ([]*models.ControlEvent, error)

	// GetMark(ctx context.Context, markID int) (*models.Mark, error)
	// GetAllMarks(ctx context.Context) ([]*models.Mark, error)
}
