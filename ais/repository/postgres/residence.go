package postgres

import (
	"context"
	"database/sql"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type Residence struct {
	ID        int
	Address   string
	City      string
	Community bool
	// StudentsIDs  []int
}

func toPostgresResidence(c *models.Residence) *Residence {
	return &Residence{
		ID:        c.ID,
		Address:   c.Address,
		City:      c.City,
		Community: c.Community,
		// Students....
	}
}

func toModelResidence(r DBAisRepository, ctx context.Context, c *Residence) *models.Residence {
	// student, err := r.GetStudent(ctx, c.StudentID)

	// if err != nil {
	// 	panic(err)
	// }

	return &models.Residence{
		ID:        c.ID,
		Address:   c.Address,
		City:      c.City,
		Community: c.Community,
	}
}

// const createCathedraQuery = `INSERT INTO Кафедра(название, короткое_название) VALUES ( $1, $2 )`

// func (r AisRepository) CreateCathedra(ctx context.Context, name, shortName string) error {
// 	_, err := r.db.ExecContext(ctx, createCathedraQuery,
// 		name, shortName,
// 	)

// 	return err
// }

const getResidenceQuery = `SELECT id, адрес, город, общежитие FROM МестоЖительства WHERE id = $1`

func (r DBAisRepository) GetResidence(ctx context.Context, residenceID int) (*models.Residence, error) {
	row := r.db.QueryRowContext(ctx, getResidenceQuery, residenceID)

	residence := new(Residence)
	err := row.Scan(&residence.ID, &residence.Address, &residence.City, &residence.Community)
	// TODO: добавить студентов

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ais.ErrResidenceNotFound
		}
		return nil, err
	}

	return toModelResidence(r, ctx, residence), nil
}

const getAllResidencesQuery = `SELECT id, адрес, город, общежитие FROM МестоЖительства`

func (r DBAisRepository) GetAllResidences(ctx context.Context) ([]*models.Residence, error) {
	rows, err := r.db.QueryContext(ctx, getAllResidencesQuery)
	residences := make([]*models.Residence, 0)

	if err != nil {
		return residences, err
	}

	for rows.Next() {
		residence := new(Residence)
		if err := rows.Scan(&residence.ID, &residence.Address, &residence.City, &residence.Community); err != nil {
			return []*models.Residence{}, err
		}
		residences = append(residences, toModelResidence(r, ctx, residence))
	}

	return residences, nil
}
