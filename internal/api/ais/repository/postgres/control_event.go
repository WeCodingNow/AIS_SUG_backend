package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/WeCodingNow/AIS_SUG_backend/pkg/pgorm"
)

// CREATE TABLE КонтрольноеМероприятие(
//     id SERIAL,
//     id_типа int NOT NULL references ТипКонтрольногоМероприятия(id) ON DELETE CASCADE,
//     id_дисциплины int NOT NULL references Дисциплина(id) ON DELETE CASCADE,
//     id_семестра int NOT NULL references Семестр(id) ON DELETE CASCADE,
//     дата_проведения date NOT NULL,
//     CONSTRAINT контрольное_мероприятие_pk PRIMARY KEY (id)
// );

type repoControlEvent struct {
	ID   int
	Date sql.NullTime

	ControlEventType *repoControlEventType
	Discipline       *repoDiscipline
	Semester         *repoSemester
	Marks            map[int]*repoMark

	model *models.ControlEvent
}

func NewRepoControlEvent() *repoControlEvent {
	return &repoControlEvent{
		Marks: make(map[int]*repoMark),
	}
}

func (s *repoControlEvent) Fill(scannable pgorm.Scannable) error {
	return scannable.Scan(&s.ID, &s.Date)
}

func (s repoControlEvent) GetID() int {
	return s.ID
}

const controlEventTable = "КонтрольноеМероприятие"
const controlEventFields = "id,дата_проведения"
const controlEventControlEventTypeFK = "id_типа"
const controlEventDisciplineFK = "id_дисциплины"
const controlEventSemesterFK = "id_семестра"

func (c repoControlEvent) GetDescription() pgorm.ModelDescription {
	return pgorm.ModelDescription{
		Table:  controlEventTable,
		Fields: controlEventFields,
		Dependencies: []pgorm.ModelDependency{
			{
				DependencyType:  pgorm.ManyToOne,
				ForeignKeyField: controlEventSemesterFK,
				ModelMaker:      func() pgorm.RepoModel { return NewRepoSemester() },
			},
			{
				DependencyType:  pgorm.ManyToOne,
				ForeignKeyField: controlEventControlEventTypeFK,
				ModelMaker:      func() pgorm.RepoModel { return NewRepoControlEventType() },
			},
			{
				DependencyType:  pgorm.ManyToOne,
				ForeignKeyField: controlEventDisciplineFK,
				ModelMaker:      func() pgorm.RepoModel { return NewRepoDiscipline() },
			},
			{
				DependencyType:     pgorm.OneToMany,
				DepForeignKeyField: markControlEventFK,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoMark() },
			},
		},
	}
}

func (c *repoControlEvent) toModel() *models.ControlEvent {
	if c.model == nil {
		c.model = &models.ControlEvent{
			ID:   c.ID,
			Date: c.Date.Time,
		}

		c.model.ControlEventType = c.ControlEventType.toModel()
		c.model.Semester = c.Semester.toModel()
		c.model.Discipline = c.Discipline.toModel()

		marks := make([]*models.Mark, 0, len(c.Marks))

		for _, repoM := range c.Marks {
			marks = append(marks, repoM.toModel())
		}

		c.model.Marks = marks
	}

	return c.model
}

func (s *repoControlEvent) AcceptDep(dep interface{}) error {
	switch dep := dep.(type) {
	case *repoSemester:
		s.Semester = dep
	case *repoControlEventType:
		s.ControlEventType = dep
	case *repoDiscipline:
		s.Discipline = dep
	case *repoMark:
		s.Marks[dep.ID] = dep
	default:
		return fmt.Errorf("no dependency for %v", dep)
	}
	return nil
}

func (r *DBAisRepository) GetControlEvent(ctx context.Context, id int) (*models.ControlEvent, error) {
	controlEvent := NewRepoControlEvent()
	filler, err := pgorm.MakeFiller(ctx, r.db, controlEventFields, controlEventTable, &id)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, ais.ErrControlEventNotFound
	}

	err = filler.Fill(controlEvent)

	return controlEvent.toModel(), nil
}

func (r *DBAisRepository) GetAllControlEvents(ctx context.Context) ([]*models.ControlEvent, error) {
	controlEvents := make([]*models.ControlEvent, 0)
	filler, err := pgorm.MakeFiller(ctx, r.db, controlEventFields, controlEventTable, nil)

	if err != nil {
		return nil, err
	}

	for filler.Next() {
		newRepoControlEvent := NewRepoControlEvent()
		err = filler.Fill(newRepoControlEvent)
		if err != nil {
			return nil, err
		}
		controlEvents = append(controlEvents, newRepoControlEvent.toModel())
	}

	return controlEvents, nil
}
