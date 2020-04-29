package postgres

import (
	"context"
	"time"

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
	return nil, nil
}

const getMarkQuery = `SELECT id, id_контрольного_мероприятия, id_студента, дата_получения, значение FROM Оценка WHERE id = $1`

func (m *Mark) hydrate(sc postgres.Scannable) error {
	return sc.Scan(&m.ID, &m.ControlEventID, &m.StudentID, &m.Date, &m.Value)
}

func (r DBAisRepository) GetMark(ctx context.Context, markID int) (*models.Mark, error) {
	// row := r.db.QueryRowContext(ctx, getMarkQuery, markID)

	return nil, nil
}

const getAllMarksQuery = `SELECT id, id_контрольного_мероприятия, id_студента, дата_получения, значение FROM Оценка`

func (r DBAisRepository) GetAllMarks(ctx context.Context) ([]*models.Mark, error) {
	errValue := []*models.Mark{}

	// rows, err := r.db.QueryContext(ctx, getAllMarksQuery)
	// marks := make([]*models.Mark, 0)

	return errValue, nil
}
