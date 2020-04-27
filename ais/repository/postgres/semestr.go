package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

// CREATE TABLE Семестр(
//     id SERIAL,
//     номер int NOT NULL,
//     начало date NOT NULL,
//     конец date,
//     CONSTRAINT семестр_pk PRIMARY KEY (id)
// );

type Semester struct {
	ID        int
	Number    int
	Beginning sql.NullTime
	End       sql.NullTime
}

func toPostgresSemester(c *models.Semester) *Semester {
	endTime := sql.NullTime{time.Time{}, false}

	if c.End != nil {
		endTime = sql.NullTime{*c.End, true}
	}

	return &Semester{
		ID:        c.ID,
		Number:    c.Number,
		Beginning: sql.NullTime{c.Beginning, true},
		End:       endTime,
	}
}

func toModelSemester(c *Semester) *models.Semester {
	retModel := &models.Semester{
		ID:        c.ID,
		Number:    c.Number,
		Beginning: c.Beginning.Time,
	}

	c.End.Scan(&retModel.End)

	return retModel
}

const createSemesterQuery = `INSERT INTO Семестр(номер, начало, конец) VALUES ( $1, $2, $3 )`

func (r AisRepository) CreateSemester(ctx context.Context, number int, beginning time.Time, end *time.Time) error {
	thirdArg := sql.NullTime{Time: time.Time{}, Valid: false}

	if end != nil {
		thirdArg.Time = *end
	}

	_, err := r.db.ExecContext(ctx, createSemesterQuery,
		number, beginning, thirdArg,
	)

	return err
}

const getSemesterQuery = `SELECT * FROM Семестр WHERE id = $1`

func (r AisRepository) GetSemester(ctx context.Context, semesterID int) (*models.Semester, error) {
	row := r.db.QueryRowContext(ctx, getSemesterQuery, semesterID)

	semester := new(Semester)
	err := row.Scan(&semester.ID, &semester.Number, &semester.Beginning, &semester.End)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ais.ErrSemesterNotFound
		}
		return nil, err
	}

	return toModelSemester(semester), nil
}

const getAllSemestersQuery = `SELECT * FROM Семестр`

func (r AisRepository) GetAllSemesters(ctx context.Context) ([]*models.Semester, error) {
	rows, err := r.db.QueryContext(ctx, getAllSemestersQuery)
	semesters := make([]*models.Semester, 0)

	if err != nil {
		return semesters, err
	}

	for rows.Next() {
		semester := new(Semester)
		if err := rows.Scan(&semester.ID, &semester.Number, &semester.Beginning, &semester.End); err != nil {
			return []*models.Semester{}, err
		}
		semesters = append(semesters, toModelSemester(semester))
	}

	return semesters, nil
}
