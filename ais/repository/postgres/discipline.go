package postgres

import (
	"context"
	"database/sql"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type Discipline struct {
	ID    int
	Name  string
	Hours int
}

func toPostgresDiscipline(c *models.Discipline) *Discipline {
	return &Discipline{
		c.ID,
		c.Name,
		c.Hours,
	}
}

func toModelDiscipline(c *Discipline) *models.Discipline {
	return &models.Discipline{
		c.ID,
		c.Name,
		c.Hours,
	}
}

// const createCathedraQuery = `INSERT INTO Кафедра(название, короткое_название) VALUES ( $1, $2 )`

// func (r AisRepository) CreateCathedra(ctx context.Context, name, shortName string) error {
// 	_, err := r.db.ExecContext(ctx, createCathedraQuery,
// 		name, shortName,
// 	)

// 	return err
// }

const getDisciplineQuery = `SELECT id, название, часы FROM Дисциплина WHERE id = $1`

func (r AisRepository) GetDiscipline(ctx context.Context, disciplineID int) (*models.Discipline, error) {
	row := r.db.QueryRowContext(ctx, getDisciplineQuery, disciplineID)

	discipline := new(Discipline)
	err := row.Scan(&discipline.ID, &discipline.Name, &discipline.Hours)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ais.ErrDisciplineNotFound
		}
		return nil, err
	}

	return toModelDiscipline(discipline), nil
}

const getAllDisciplinesQuery = `SELECT id, название, часы FROM Дисциплина`

func (r AisRepository) GetAllDisciplines(ctx context.Context) ([]*models.Discipline, error) {
	rows, err := r.db.QueryContext(ctx, getAllDisciplinesQuery)
	disciplines := make([]*models.Discipline, 0)

	if err != nil {
		return disciplines, err
	}

	for rows.Next() {
		discipline := new(Discipline)
		if err := rows.Scan(&discipline.ID, &discipline.Name, &discipline.Hours); err != nil {
			return []*models.Discipline{}, err
		}
		disciplines = append(disciplines, toModelDiscipline(discipline))
	}

	return disciplines, nil
}
