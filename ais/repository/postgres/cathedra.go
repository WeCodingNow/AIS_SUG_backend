package postgres

import (
	"context"
	"database/sql"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

// CREATE TABLE Кафедра(
//     id SERIAL,
//     название varchar(100) NOT NULL UNIQUE,
//     короткое_название varchar(10) NOT NULL UNIQUE,
//     CONSTRAINT кафедра_pk PRIMARY KEY (id)
// );

type Cathedra struct {
	ID        int
	Name      string
	ShortName string
}

func toPostgresCathedra(c *models.Cathedra) *Cathedra {
	return &Cathedra{
		c.ID,
		c.Name,
		c.ShortName,
	}
}

func toModelCathedra(c *Cathedra) *models.Cathedra {
	return &models.Cathedra{
		c.ID,
		c.Name,
		c.ShortName,
	}
}

// const createCathedraQuery = `INSERT INTO Кафедра(название, короткое_название) VALUES ( $1, $2 )`

// func (r AisRepository) CreateCathedra(ctx context.Context, name, shortName string) error {
// 	_, err := r.db.ExecContext(ctx, createCathedraQuery,
// 		name, shortName,
// 	)

// 	return err
// }

const getCathedraQuery = `SELECT * FROM Кафедра WHERE id = $1`

func (r DBAisRepository) GetCathedra(ctx context.Context, cathedraID int) (*models.Cathedra, error) {
	row := r.db.QueryRowContext(ctx, getCathedraQuery, cathedraID)

	cathedra := new(Cathedra)
	err := row.Scan(&cathedra.ID, &cathedra.Name, &cathedra.ShortName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ais.ErrCathedraNotFound
		}
		return nil, err
	}

	return toModelCathedra(cathedra), nil
}

const getAllCathedrasQuery = `SELECT * FROM Кафедра`

func (r DBAisRepository) GetAllCathedras(ctx context.Context) ([]*models.Cathedra, error) {
	rows, err := r.db.QueryContext(ctx, getAllCathedrasQuery)
	cathedras := make([]*models.Cathedra, 0)

	if err != nil {
		return cathedras, err
	}

	for rows.Next() {
		cathedra := new(Cathedra)
		if err := rows.Scan(&cathedra.ID, &cathedra.Name, &cathedra.ShortName); err != nil {
			return []*models.Cathedra{}, err
		}
		cathedras = append(cathedras, toModelCathedra(cathedra))
	}

	return cathedras, nil
}
