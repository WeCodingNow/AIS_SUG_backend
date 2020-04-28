package postgres

import (
	"context"
	"database/sql"
	"time"

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
	ID             int
	ControlEventID int
	StudentID      int
	Date           time.Time
	Value          int
}

func toPostgresMark(c *models.Mark) *Mark {
	return nil
}

func toModelMark(r DBAisRepository, ctx context.Context, m *Mark) (*models.Mark, error) {
	controlEvent, err := r.GetControlEvent(ctx, m.ControlEventID)
	if err != nil {
		return nil, err
	}

	student, err := r.GetStudent(ctx, m.StudentID)
	if err != nil {
		return nil, err
	}

	return &models.Mark{
		ID:           m.ID,
		ControlEvent: controlEvent,
		Student:      student,
		Date:         m.Date,
		Value:        m.Value,
	}, nil
}

const getMarkQuery = `SELECT id, id_контрольного_мероприятия, id_студента, дата_получения, значение FROM Оценка WHERE id = $1`

func (m *Mark) hydrate(sc postgres.Scannable) error {
	return sc.Scan(&m.ID, &m.ControlEventID, &m.StudentID, &m.Date, &m.Value)
}

func (m *Mark) Fill(sc postgres.Scannable) error {
	err := m.hydrate(sc)

	if err != nil {
		if err == sql.ErrNoRows {
			return ais.ErrMarkNotFound
		}
		return err
	}

	return nil
}

func (r DBAisRepository) GetMark(ctx context.Context, markID int) (*models.Mark, error) {
	row := r.db.QueryRowContext(ctx, getMarkQuery, markID)

	mark := new(Mark)
	err := mark.Fill(row)

	if err != nil {
		return nil, err
	}

	return toModelMark(r, ctx, mark)
}

const getAllMarksQuery = `SELECT id, id_контрольного_мероприятия, id_студента, дата_получения, значение FROM Оценка`

func (r DBAisRepository) GetAllMarks(ctx context.Context) ([]*models.Mark, error) {
	rows, err := r.db.QueryContext(ctx, getAllMarksQuery)
	marks := make([]*models.Mark, 0)

	if err != nil {
		return marks, err
	}

	for rows.Next() {
		mark := new(Mark)
		err := mark.Fill(rows)

		if err != nil {
			return []*models.Mark{}, nil
		}

		markModel, err := toModelMark(r, ctx, mark)
		if err != nil {
			return []*models.Mark{}, err
		}

		marks = append(marks, markModel)
	}

	return marks, nil
}
