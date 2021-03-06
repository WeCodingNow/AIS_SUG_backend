package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/utils"
	"github.com/WeCodingNow/AIS_SUG_backend/pkg/pgorm"
)

// CREATE TABLE Оценка(
//     id SERIAL,
//     id_контрольного_мероприятия int NOT NULL references КонтрольноеМероприятие(id) ON DELETE CASCADE,
//     id_студента int NOT NULL references Студент(id) ON DELETE CASCADE,
//     дата_получения date NOT NULL,
//     значение int NOT NULL,
//     CONSTRAINT оценка_pk PRIMARY KEY (id)
// );

type repoMark struct {
	ID    int
	Date  sql.NullTime
	Value int

	ControlEvent *repoControlEvent
	Student      *repoStudent

	model *models.Mark
}

func NewRepoMark() *repoMark {
	return &repoMark{}
}

func (s *repoMark) Fill(scannable pgorm.Scannable) error {
	return scannable.Scan(&s.ID, &s.Date, &s.Value)
}

func (s repoMark) GetID() int {
	return s.ID
}

const markTable = "Оценка"
const markFields = "id,дата_получения,значение"
const markControlEventFK = "id_контрольного_мероприятия"
const markStudentFK = "id_студента"

func (c repoMark) GetDescription() pgorm.ModelDescription {
	return pgorm.ModelDescription{
		Table:  markTable,
		Fields: markFields,
		Dependencies: []pgorm.ModelDependency{
			{
				DependencyType:  pgorm.ManyToOne,
				ForeignKeyField: markControlEventFK,
				ModelMaker:      func() pgorm.RepoModel { return NewRepoControlEvent() },
			},
			{
				DependencyType:  pgorm.ManyToOne,
				ForeignKeyField: markStudentFK,
				ModelMaker:      func() pgorm.RepoModel { return NewRepoStudent() },
			},
		},
	}
}

func (c *repoMark) toModel() *models.Mark {
	if c.model == nil {
		c.model = &models.Mark{
			ID:    c.ID,
			Date:  c.Date.Time,
			Value: c.Value,
		}

		c.model.ControlEvent = c.ControlEvent.toModel()
		c.model.Student = c.Student.toModel()
	}

	return c.model
}

func (s *repoMark) AcceptDep(dep interface{}) error {
	switch dep := dep.(type) {
	case *repoControlEvent:
		s.ControlEvent = dep
	case *repoStudent:
		s.Student = dep
	default:
		return fmt.Errorf("no dependency for %v", dep)
	}
	return nil
}

func (r *DBAisRepository) GetMark(ctx context.Context, id int) (*models.Mark, error) {
	mark := NewRepoMark()
	filler, err := pgorm.MakeFiller(ctx, r.db, markFields, markTable, &id)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, ais.ErrMarkNotFound
	}

	err = filler.Fill(mark)

	return mark.toModel(), nil
}

func (r *DBAisRepository) GetAllMarks(ctx context.Context) ([]*models.Mark, error) {
	marks := make([]*models.Mark, 0)
	filler, err := pgorm.MakeFiller(ctx, r.db, markFields, markTable, nil)

	if err != nil {
		return nil, err
	}

	for filler.Next() {
		newRepoMark := NewRepoMark()
		err = filler.Fill(newRepoMark)
		if err != nil {
			return nil, err
		}
		marks = append(marks, newRepoMark.toModel())
	}

	return marks, nil
}

func (r *DBAisRepository) CreateMark(ctx context.Context, date time.Time, value, controlEventID, studentID int) (*models.Mark, error) {
	row := r.db.QueryRowContext(ctx,
		fmt.Sprintf("%s RETURNING %s", utils.MakeInsertQuery(markTable, "id_контрольного_мероприятия", "id_студента", "дата_получения", "значение"), markFields),
		controlEventID, studentID, date, value,
	)

	repoMark := NewRepoMark()

	err := repoMark.Fill(row)

	if err != nil {
		return nil, err
	}

	return r.GetMark(ctx, repoMark.ID)
}
