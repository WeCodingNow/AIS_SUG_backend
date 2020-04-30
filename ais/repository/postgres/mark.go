package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/WeCodingNow/AIS_SUG_backend/utils/delivery/postgres"
)

// CREATE TABLE Оценка(
//     id SERIAL,
//     id_контрольного_мероприятия int NOT NULL references КонтрольноеМероприятие(id) ON DELETE CASCADE,
//     id_студента int NOT NULL references Студент(id) ON DELETE CASCADE,
//     дата_получения date NOT NULL,
//     значение int NOT NULL,
//     CONSTRAINT оценка_pk PRIMARY KEY (id)
// );
type Mark struct {
	ID    int
	Date  sql.NullTime
	Value int

	ControlEvent *ControlEvent
	Student      *Student
}

const markTable = "Оценка"
const markIDField = "id"
const markFields = "id,дата_получения,значение"
const markControlEventFK = "id_контрольного_мероприятия"
const markStudentFK = "id_студента"

func (m *Mark) toModel(controlEventRef *models.ControlEvent, studentRef *models.Student) *models.Mark {
	mark := &models.Mark{
		ID:    m.ID,
		Date:  m.Date.Time,
		Value: m.Value,

		ControlEvent: controlEventRef,
		Student:      studentRef,
	}

	if mark.Student == nil {
		mark.Student = m.Student.toModel(nil, nil, nil, mark)
	}

	if mark.ControlEvent == nil {
		mark.ControlEvent = m.ControlEvent.toModel(nil, nil, mark)
	}

	return mark
}

func NewPostgresMark(scannable postgres.Scannable) (*Mark, error) {
	mark := &Mark{}

	err := scannable.Scan(&mark.ID, &mark.Date, &mark.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ais.ErrMarkNotFound
		}
		return nil, err
	}

	return mark, nil
}

func (m *Mark) Associate(ctx context.Context, r DBAisRepository, controlEventRef *ControlEvent, studentRef *Student) error {
	if controlEventRef == nil {
		controlEventRow := r.db.QueryRowContext(
			ctx,
			postgres.MakeJoinQuery(controlEventTable, controlEventFields, controlEventIDField, markTable, markControlEventFK, markIDField),
			m.ID,
		)

		controlEvent, err := NewPostgresControlEvent(controlEventRow)

		if err != nil {
			return err
		}

		controlEvent.Associate(ctx, r, nil, nil, m)
		m.ControlEvent = controlEvent
	} else {
		m.ControlEvent = controlEventRef
	}

	if studentRef == nil {
		studentRow := r.db.QueryRowContext(
			ctx,
			postgres.MakeJoinQuery(studentTable, studentFields, studentIDField, markTable, markControlEventFK, markIDField),
			m.ID,
		)

		student, err := NewPostgresStudent(studentRow)

		if err != nil {
			return err
		}

		student.Associate(ctx, r, nil, nil, nil, m)
		m.Student = student
	} else {
		m.Student = studentRef
	}

	return nil
}

func makeMarkModel(ctx context.Context, r DBAisRepository, scannable postgres.Scannable) (*models.Mark, error) {
	mark, err := NewPostgresMark(scannable)

	if err != nil {
		return nil, err
	}

	err = mark.Associate(ctx, r, nil, nil)

	if err != nil {
		return nil, err
	}

	return mark.toModel(nil, nil), nil
}

func (r DBAisRepository) GetMark(ctx context.Context, markID int) (*models.Mark, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", markFields, markTable), markID)
	return makeMarkModel(ctx, r, row)
}

func (r DBAisRepository) GetAllMarks(ctx context.Context) ([]*models.Mark, error) {
	errValue := []*models.Mark{}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s", markFields, markTable))

	if err != nil {
		return errValue, err
	}

	marks := []*models.Mark{}
	for rows.Next() {
		mark, err := makeMarkModel(ctx, r, rows)

		if err != nil {
			return errValue, nil
		}

		marks = append(marks, mark)
	}

	return marks, nil
}
